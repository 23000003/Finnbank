package services

import (
	"context"
	"finnbank/common/utils"
	"finnbank/graphql-api/types"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const DefaultTransactionStatus = "Pending"

type TransactionService struct {
	db *pgxpool.Pool
	l  *utils.Logger
}

func NewTransactionService(db *pgxpool.Pool, logger *utils.Logger) *TransactionService {
	return &TransactionService{
		db: db,
		l:  logger,
	}
}

func GenerateUUID() string {
	return uuid.New().String()
}

// GetTransactionByUserId retrieves all transactions for a specific user.
func (s *TransactionService) GetTransactionByUserId(ctx context.Context, userId string) ([]types.Transaction, error) {
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

	var transactions []types.Transaction
	for rows.Next() {
		var txn types.Transaction
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

// CreateTransaction creates a new transaction.
func (s *TransactionService) CreateTransaction(ctx context.Context, req types.Transaction) (types.Transaction, error) {
	s.l.Info("Creating transaction for user ID: %s", req.SenderID)

	// Input validation
	if req.SenderID == "" || req.ReceiverID == "" {
		return types.Transaction{}, fmt.Errorf("sender_id and receiver_id cannot be empty")
	}
	if req.Amount < 0 {
		return types.Transaction{}, fmt.Errorf("amount cannot be negative")
	}
	if req.TransactionFee < 0 {
		return types.Transaction{}, fmt.Errorf("transaction_fee cannot be negative")
	}

	query := `
        INSERT INTO transactions (transaction_id, ref_no, sender_id, receiver_id, transaction_type, amount, transaction_status, date_transaction, transaction_fee, notes)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING transaction_id, ref_no, sender_id, receiver_id, transaction_type, amount, transaction_status, date_transaction, transaction_fee, notes
    `

	newTransaction := req
	newTransaction.TransactionID = GenerateUUID() // Generate a unique ID
	newTransaction.DateTransaction = time.Now()
	newTransaction.TransactionStatus = DefaultTransactionStatus

	row := s.db.QueryRow(ctx, query,
		newTransaction.TransactionID,
		newTransaction.RefNo,
		newTransaction.SenderID,
		newTransaction.ReceiverID,
		newTransaction.TransactionType,
		newTransaction.Amount,
		newTransaction.TransactionStatus,
		newTransaction.DateTransaction,
		newTransaction.TransactionFee,
		newTransaction.Notes,
	)

	err := row.Scan(
		&newTransaction.TransactionID,
		&newTransaction.RefNo,
		&newTransaction.SenderID,
		&newTransaction.ReceiverID,
		&newTransaction.TransactionType,
		&newTransaction.Amount,
		&newTransaction.TransactionStatus,
		&newTransaction.DateTransaction,
		&newTransaction.TransactionFee,
		&newTransaction.Notes,
	)
	if err != nil {
		s.l.Error("Error creating transaction with RefNo %s for user ID %s: %v", req.RefNo, req.SenderID, err)
		return types.Transaction{}, fmt.Errorf("failed to create transaction for user ID %s: %w", req.SenderID, err)
	}

	s.l.Info("Transaction created successfully: %+v", newTransaction)
	return newTransaction, nil
}
