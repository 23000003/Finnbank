package services

import (
	"context"
	"finnbank/common/utils"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OpenedAccounts struct {
	OpenedAccountID     int       `json:"openedaccount_id"`
	BankCardID          *int      `json:"bankcard_id"`
	Balance             float64   `json:"balance"`
	AccountType         string    `json:"account_type"`
	DateCreated         time.Time `json:"date_created"`
	OpenedAccountStatus string    `json:"openedaccount_status"`
}

type OpenedAccountService struct {
	db *pgxpool.Pool
	l  *utils.Logger
}

func NewOpenedAccountService(db *pgxpool.Pool, logger *utils.Logger) *OpenedAccountService {
	return &OpenedAccountService{
		db: db,
		l:  logger,
	}
}

func (s *OpenedAccountService) GetAllOpenedAccountsByUserId(ctx context.Context, id int) ([]*OpenedAccounts, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx,
		`SELECT 
			openedaccount_id, bankcard_id, balance, 
			account_type, openedaccount_status, date_created 
		FROM openedaccount 
		WHERE account_id = $1`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []*OpenedAccounts
	for rows.Next() {
		var acc OpenedAccounts
		if err := rows.Scan(
			&acc.OpenedAccountID,
			&acc.BankCardID,
			&acc.Balance,
			&acc.AccountType,
			&acc.OpenedAccountStatus,
			&acc.DateCreated,
		); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		results = append(results, &acc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	s.l.Info("All opened accounts: %v", results)
	return results, nil
}

func (s *OpenedAccountService) GetOpenedAccountById(ctx context.Context, id int) (*OpenedAccounts, error) {
	var acc OpenedAccounts

	query := `
		SELECT 
			openedaccount_id, bankcard_id, balance, 
			account_type, openedaccount_status, date_created 
		FROM openedaccount 
		WHERE openedaccount_id = $1
	`

	err := s.db.QueryRow(ctx, query, id).Scan(
		&acc.OpenedAccountID,
		&acc.BankCardID,
		&acc.Balance,
		&acc.AccountType,
		&acc.OpenedAccountStatus,
		&acc.DateCreated,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch opened account: %w", err)
	}

	return &acc, nil
}

// Create a new opened account
func (s *OpenedAccountService) CreateOpenedAccount(ctx context.Context, accountId int, accountType string, balance float64) (*OpenedAccounts, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	var bankcardId *int = nil
	if accountType != "savings" {
		bankcardId = new(int)
		*bankcardId = 1 // TODO: Replace this with actual BankCardService integration
	}

	var openedAccountId int
	err = conn.QueryRow(ctx,
		`INSERT INTO openedaccount (account_id, bankcard_id, balance, account_type) 
		 VALUES ($1, $2, $3, $4) RETURNING openedaccount_id`,
		accountId, bankcardId, balance, accountType,
	).Scan(&openedAccountId)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return &OpenedAccounts{
		OpenedAccountID:     openedAccountId,
		AccountType:         accountType,
		Balance:             balance,
		DateCreated:         time.Now(),
		BankCardID:          bankcardId,
		OpenedAccountStatus: "Active", // default status (optional)
	}, nil
}

func (s *OpenedAccountService) UpdateOpenedAccountStatus(ctx context.Context, openedAccountId int, status string) (string, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	var openedaccount_id int
	err = conn.QueryRow(ctx,
		`UPDATE openedaccount SET openedaccount_status = $1 
		 WHERE openedaccount_id = $2 
		 RETURNING openedaccount_id`,
		status, openedAccountId).Scan(&openedaccount_id)

	if err != nil {
		return "", fmt.Errorf("update failed: %w", err)
	}

	return "Update Successful", nil
}
