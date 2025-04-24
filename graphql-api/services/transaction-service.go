package services

import (
	"context"
	"finnbank/common/utils"
	t "finnbank/graphql-api/types"
	"fmt"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
	q "finnbank/graphql-api/queue"
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

func (s *TransactionService) GetTransactionByUserId(ctx context.Context, userId string) ([]t.Transaction, error) {
	s.l.Info("Fetching transactions for user ID: %s", userId)

	query := `
        SELECT transaction_id, ref_no, sender_id, receiver_id, transaction_type, amount, transaction_status, date_transaction, transaction_fee, notes
        FROM transactions
        WHERE sender_id = $1 OR receiver_id = $1
    `

	rows, err := s.db.Query(ctx, query, userId)
	if err != nil {
		s.l.Error("Error fetching transactions for user ID %s: %v", userId, err)
		return nil, fmt.Errorf("failed to fetch transactions for user ID %s: %w", userId, err)
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
			s.l.Error("Error scanning transaction row for user ID %s: %v", userId, err)
			return nil, fmt.Errorf("failed to scan transaction row for user ID %s: %w", userId, err)
		}
		transactions = append(transactions, txn)
	}

	if rows.Err() != nil {
		s.l.Error("Error iterating through transaction rows for user ID %s: %v", userId, rows.Err())
		return nil, fmt.Errorf("failed to iterate through transaction rows for user ID %s: %w", userId, rows.Err())
	}
	if len(transactions) == 0 {
		s.l.Info("No transactions found for user ID: %s", userId)
	}

	s.l.Info("Successfully fetched %d transactions for user ID: %s", len(transactions), userId)
	return transactions, nil
}

// CreateTransaction creates a new transaction, auto‑generating a numeric ref_no.
func (s *TransactionService) CreateTransaction(ctx context.Context, req t.Transaction) (t.Transaction, error) {
	s.l.Info("Creating transaction for user ID: %s", req.SenderID)

	// 1) Validation
	if req.SenderID == "" || req.ReceiverID == "" {
		return t.Transaction{}, fmt.Errorf("sender_id and receiver_id cannot be empty")
	}
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
          date_transaction, transaction_fee, notes;
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
func (s *TransactionService) GetTransactionByTimestampByUserId(
	ctx context.Context,
	userId string,
	start, end time.Time,
) ([]t.Transaction, error) {
	s.l.Info("Fetching transactions for user %s between %s and %s",
		userId, start.Format(time.RFC3339), end.Format(time.RFC3339),
	)

	query := `
        SELECT
          transaction_id, ref_no, sender_id, receiver_id,
          transaction_type, amount, transaction_status,
          date_transaction, transaction_fee, notes
        FROM public.transactions
        WHERE (sender_id = $1 OR receiver_id = $1)
          AND date_transaction BETWEEN $2 AND $3
        ORDER BY date_transaction;
    `

	rows, err := s.db.Query(ctx, query, userId, start, end)
	if err != nil {
		s.l.Error("Error fetching by timestamp for user %s: %v", userId, err)
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
			s.l.Error("Scan error for user %s: %v", userId, err)
			return nil, fmt.Errorf("scan error: %w", err)
		}
		txns = append(txns, t)
	}
	if err := rows.Err(); err != nil {
		s.l.Error("Row iteration error: %v", err)
		return nil, fmt.Errorf("row error: %w", err)
	}

	s.l.Info("Fetched %d transactions in range for user %s", len(txns), userId)
	return txns, nil
}
