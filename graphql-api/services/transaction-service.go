package services

import (
	"context"
	"finnbank/common/utils"
	q "finnbank/graphql-api/queue"
	t "finnbank/graphql-api/types"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	DefaultTransactionStatus = "Pending"
	RefNoLength              = 12
)

type TransactionService struct {
	db *pgxpool.Pool
	l  *utils.Logger
	queue *q.Queue
}

func NewTransactionService(db *pgxpool.Pool, logger *utils.Logger, q *q.Queue) *TransactionService {
	return &TransactionService{db: db, l: logger, queue: q}
}

// generateRefNo returns a random numeric string of length RefNoLength.

func (s *TransactionService) GetTransactionByUserId(ctx context.Context, creditId int, debitId int, savingsId int, limit int) ([]t.Transaction, error) {

	query := `
		SELECT 
			transaction_id, ref_no, sender_id, receiver_id, transaction_type, amount, 
			transaction_status, date_transaction, transaction_fee, notes
		FROM transactions
		WHERE sender_id IN ($1, $2, $3) OR receiver_id IN ($1, $2, $3)
		ORDER BY date_transaction DESC
		LIMIT $4;
	`

	rows, err := s.db.Query(ctx, query, creditId, debitId, savingsId, limit)
	if err != nil {
		s.l.Error("Error fetching transactions %v", err)
		return nil, fmt.Errorf("failed to fetch transactions %w", err)
	}
	defer rows.Close()

	var transactions []t.Transaction
	for rows.Next() {
		var txn t.Transaction
		err := rows.Scan(
			&txn.TransactionID,
			&txn.RefNo,
			&txn.SenderID,
			&txn.ReceiverID,
			&txn.TransactionType,
			&txn.Amount,
			&txn.TransactionStatus,
			&txn.DateTransaction,
			&txn.TransactionFee,
			&txn.Notes,
		)
		if err != nil {
			s.l.Error("Error scanning transaction row %v", err)
			return nil, fmt.Errorf("failed to scan transaction row %w", err)
		}
		transactions = append(transactions, txn)
	}

	if rows.Err() != nil {
		s.l.Error("Error iterating through transaction rows %v", rows.Err())
		return nil, fmt.Errorf("failed to iterate through transaction rows %w", rows.Err())
	}

	return transactions, nil
}

func (s *TransactionService) GetRecentlySent(ctx context.Context, creditId int, debitId int, savingsId int) ([]t.Transaction, error) {

	query := `
		SELECT 
			transaction_id, ref_no, sender_id, receiver_id, transaction_type, amount, 
			transaction_status, date_transaction, transaction_fee, notes
		FROM transactions
		WHERE sender_id IN ($1, $2, $3) AND transaction_status = 'Completed'
		ORDER BY date_transaction DESC
		LIMIT 2;
	`

	rows, err := s.db.Query(ctx, query, creditId, debitId, savingsId)
	if err != nil {
		s.l.Error("Error fetching transactions %v", err)
		return nil, fmt.Errorf("failed to fetch transactions %w", err)
	}
	defer rows.Close()

	var transactions []t.Transaction
	for rows.Next() {
		var txn t.Transaction
		err := rows.Scan(
			&txn.TransactionID,
			&txn.RefNo,
			&txn.SenderID,
			&txn.ReceiverID,
			&txn.TransactionType,
			&txn.Amount,
			&txn.TransactionStatus,
			&txn.DateTransaction,
			&txn.TransactionFee,
			&txn.Notes,
		)
		if err != nil {
			s.l.Error("Error scanning transaction row %v", err)
			return nil, fmt.Errorf("failed to scan transaction row %w", err)
		}
		transactions = append(transactions, txn)
	}

	if rows.Err() != nil {
		s.l.Error("Error iterating through transaction rows %v", rows.Err())
		return nil, fmt.Errorf("failed to iterate through transaction rows %w", rows.Err())
	}

	return transactions, nil
}


// CreateTransaction creates a new transaction, auto‑generating a numeric ref_no.
func (s *TransactionService) CreateTransaction(ctx context.Context, req t.Transaction, transacConn *websocket.Conn) (t.Transaction, error) {

	if req.Amount < 0 {
		return t.Transaction{}, fmt.Errorf("amount cannot be negative")
	}
	if req.TransactionFee < 0 {
		return t.Transaction{}, fmt.Errorf("transaction_fee cannot be negative")
	}

	// 2) Generate numeric-only ref_no
	refNo, err := generateRandomNumber(RefNoLength)
	if err != nil {
		s.l.Error("Error generating ref_no: %v", err)
		return t.Transaction{}, fmt.Errorf("failed to generate ref_no: %w", err)
	}

	// 3) Insert, letting the DB fill transaction_id and timestamp
	query := `
        INSERT INTO transactions
          (ref_no, sender_id, receiver_id,
           transaction_type, amount, transaction_status,
           transaction_fee, notes)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
        RETURNING
          transaction_id, ref_no, sender_id, receiver_id,
          transaction_type, amount, transaction_status,
          date_transaction, transaction_fee, notes
    `
	row := s.db.QueryRow(ctx, query,
		refNo,
		req.SenderID,
		req.ReceiverID,
		req.TransactionType,
		req.Amount,
		DefaultTransactionStatus,
		req.TransactionFee,
		req.Notes,
	)

	var created t.Transaction
	if err := row.Scan(
		&created.TransactionID,
		&created.RefNo,
		&created.SenderID,
		&created.ReceiverID,
		&created.TransactionType,
		&created.Amount,
		&created.TransactionStatus,
		&created.DateTransaction,
		&created.TransactionFee,
		&created.Notes,
	); err != nil {
		s.l.Error("Error creating transaction: %v", err)
		return t.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	s.l.Info("Transaction created: %+v", created)

	s.queue.Enqueue(created.TransactionID, created.SenderID, created.ReceiverID, created.Amount)

	// sendTransac := t.Transaction{
	// 	TransactionID:     created.TransactionID,
	// 	RefNo:             created.RefNo,
	// 	SenderID:          created.SenderID,
	// 	ReceiverID:        created.ReceiverID,
	// 	TransactionType:   created.TransactionType,
	// 	Amount:            created.Amount,
	// 	TransactionStatus: created.TransactionStatus,
	// 	DateTransaction:   created.DateTransaction,
	// 	TransactionFee:    created.TransactionFee,
	// 	Notes:             created.Notes,
	// }

	// if err := transacConn.WriteJSON(sendTransac); err != nil {
	// 	s.l.Error("Error sending transaction: %v", err)
	// 	return t.Transaction{}, err
	// }

	return created, nil
}

// todo: generate random 12 digit ref_no for transaction
// add to the createTransaction function to generate a random 12 digit number and assign it to the ref_no field of the transaction struct before inserting it into the database.
// This will ensure that each transaction has a unique reference number.

//What changed?
// Removed all calls to GenerateUUID()—we let Postgres SERIAL produce transaction_id.
// Generate ref_no exactly once via generateRefNo().
// Omit transaction_id and date_transaction from the column list; the table’s defaults will fill them.
// Return the full, DB‑populated record (including your new numeric ref_no).
// With this in place, every call to CreateTransaction yields:
// A guaranteed unique, 12‑digit numeric ref_no
// A DB‑assigned transaction_id
// A timestamp set by the database
// And your Go code never mutates the same struct twice or needs to orchestrate UUIDs itself.

// ============================Get TransactionByTimeStampByUserId============================
// acceps userid, start and end time as arguments
// this method fetches transactions for a specific user between the specified start and end time
// run sql and return list of matching transaction obejct

// func (s *TransactionService) GetTransactionByUserId(ctx context.Context, creditId int, debitId int, savingsId int, limit int) ([]t.Transaction, error) {

// 	query := `
// 		SELECT 
// 			transaction_id, ref_no, sender_id, receiver_id, transaction_type, amount, 
// 			transaction_status, date_transaction, transaction_fee, notes
// 		FROM transactions
// 		WHERE sender_id IN ($1, $2, $3) OR receiver_id IN ($1, $2, $3)
// 		ORDER BY date_transaction DESC
// 		LIMIT $4;
// 	`


func (s *TransactionService) GetTransactionByTimestampByUserId(
	ctx context.Context,
	creditId int, debitId int, savingsId int,
	start, end time.Time,
) ([]t.Transaction, error) {

	query := `
		SELECT
		  transaction_id, ref_no, sender_id, receiver_id,
		  transaction_type, amount, transaction_status,
		  date_transaction, transaction_fee, notes
		FROM public.transactions
		WHERE (sender_id IN ($1, $2, $3) OR receiver_id IN ($1, $2, $3))
		  AND date_transaction BETWEEN $4 AND $5
		ORDER BY date_transaction DESC;
	`

	rows, err := s.db.Query(ctx, query, creditId, debitId, savingsId, start, end)
	if err != nil {
		s.l.Error("Error fetching by timestamp for user %d %d %d: %v", creditId, debitId, savingsId, err)
		return nil, fmt.Errorf("failed to fetch transactions in range: %w", err)
	}
	defer rows.Close()

	var txns []t.Transaction
	for rows.Next() {
		var t t.Transaction
		if err := rows.Scan(
			&t.TransactionID,
			&t.RefNo,
			&t.SenderID,
			&t.ReceiverID,
			&t.TransactionType,
			&t.Amount,
			&t.TransactionStatus,
			&t.DateTransaction,
			&t.TransactionFee,
			&t.Notes,
		); err != nil {
			s.l.Error("Scan error for user %d %d %d: %v", creditId, debitId, savingsId, err)
			return nil, fmt.Errorf("scan error: %w", err)
		}
		txns = append(txns, t)
	}
	if err := rows.Err(); err != nil {
		s.l.Error("Row iteration error: %v", err)
		return nil, fmt.Errorf("row error: %w", err)
	}

	return txns, nil
}

func (s *TransactionService) GetIsAccountAtLimitByAccountId(
	ctx context.Context, accountType string,
	creditId int, debitId int, savingsId int,
) ([]bool, error) {

	creditLimit := 100000
	debitLimit := 100000
	savingsLimit := 50000

	if accountType == "Business" {
		creditLimit = 300000
		debitLimit = 250000
		savingsLimit = 100000
	}

	query := `
		SELECT
			COALESCE(SUM(CASE WHEN sender_id = $1 THEN amount ELSE 0 END), 0) > $4 AS credit_limit_exceeded,
			COALESCE(SUM(CASE WHEN sender_id = $2 THEN amount ELSE 0 END), 0) > $5 AS debit_limit_exceeded,
			COALESCE(SUM(CASE WHEN sender_id = $3 THEN amount ELSE 0 END), 0) > $6 AS savings_limit_exceeded
		FROM transactions
		WHERE date_transaction BETWEEN NOW() - INTERVAL '1 day' AND NOW()
	`

	row := s.db.QueryRow(ctx, query, creditId, debitId, savingsId, creditLimit, debitLimit, savingsLimit)

	var credit_limit_exceeded, debit_limit_exceeded, savings_limit_exceeded bool
	if err := row.Scan(&credit_limit_exceeded, &debit_limit_exceeded, &savings_limit_exceeded); err != nil {
		s.l.Error("Error scanning limit check: %v", err)
		return nil, fmt.Errorf("failed to scan limit check: %w", err)
	}

	return []bool{credit_limit_exceeded, debit_limit_exceeded, savings_limit_exceeded}, nil
}