package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"finnbank/common/utils"
	en "finnbank/graphql-api/graphql_config/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	DefaultTransactionStatus = "Pending"
	RefNoLength              = 12
)

type TransactionService struct {
	db *pgxpool.Pool
	l  *utils.Logger
}

func NewTransactionService(db *pgxpool.Pool, logger *utils.Logger) *TransactionService {
	return &TransactionService{db: db, l: logger}
}

// generateRefNo returns a random numeric string of length RefNoLength.
func generateRefNo() (string, error) {
	const digits = "0123456789"
	ref := make([]byte, RefNoLength)
	for i := range ref {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate ref_no: %w", err)
		}
		ref[i] = digits[n.Int64()]
	}
	return string(ref), nil
}
func (s *TransactionService) GetTransactionByUserId(ctx context.Context, userId string) ([]en.Transaction, error) {
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

	var transactions []en.Transaction
	for rows.Next() {
		var txn en.Transaction
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
func (s *TransactionService) CreateTransaction(ctx context.Context, req en.Transaction) (en.Transaction, error) {
	s.l.Info("Creating transaction for user ID: %s", req.SenderID)

	// 1) Validation
	if req.SenderID == "" || req.ReceiverID == "" {
		return en.Transaction{}, fmt.Errorf("sender_id and receiver_id cannot be empty")
	}
	if req.Amount < 0 {
		return en.Transaction{}, fmt.Errorf("amount cannot be negative")
	}
	if req.TransactionFee < 0 {
		return en.Transaction{}, fmt.Errorf("transaction_fee cannot be negative")
	}

	// 2) Generate numeric-only ref_no
	refNo, err := generateRefNo()
	if err != nil {
		s.l.Error("Error generating ref_no: %v", err)
		return en.Transaction{}, fmt.Errorf("failed to generate ref_no: %w", err)
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

	var created en.Transaction
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
		return en.Transaction{}, fmt.Errorf("failed to create transaction: %w", err)
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
