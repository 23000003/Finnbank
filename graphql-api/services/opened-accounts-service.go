package services

import (
	"context"
	"finnbank/common/utils"
	"fmt"
	"time"
	t "finnbank/graphql-api/types"
	"github.com/jackc/pgx/v5/pgxpool"
)


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

func (s *OpenedAccountService) GetAllOpenedAccountsByUserId(ctx context.Context, id string) ([]*t.OpenedAccounts, error) {
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

	var results []*t.OpenedAccounts
	for rows.Next() {
		var acc t.OpenedAccounts
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

func (s *OpenedAccountService) GetOpenedAccountById(ctx context.Context, id int) (*t.OpenedAccounts, error) {
	var acc t.OpenedAccounts

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
func (s *OpenedAccountService) CreateOpenedAccount(ctx context.Context, BCService *BankcardService, data *t.CreateOpenedAccountRequest) (*t.OpenedAccounts, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	// Checks if account with same type already exists (only 1 per account_type)
	var existingAccountCount int
	err = conn.QueryRow(ctx,
		`SELECT COUNT(*) 
		 FROM openedaccount 
		 WHERE account_id = $1 AND account_type = $2`,
		data.AccountId, data.AccountType,
	).Scan(&existingAccountCount)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing account: %w", err)
	}
	if existingAccountCount > 0 {
		return nil, fmt.Errorf("account with type %s already exists for this user", data.AccountType)
	}

	var bankcardId *int = nil;
	if data.AccountType != "savings" {
		id, err := BCService.CreateCardRequest(ctx, data.AccountId, data.AccountType, data.PinNumber)
		s.l.Info("Bankcard ID: %d", id)
		if err != nil {
			return nil, fmt.Errorf("failed to create card request: %w", err)
		}
		bankcardId = &id
	}

	var openedAccountId int
	err = conn.QueryRow(ctx,
		`INSERT INTO openedaccount (account_id, bankcard_id, balance, account_type) 
		 VALUES ($1, $2, $3, $4) RETURNING openedaccount_id`,
		data.AccountId, bankcardId, data.Balance, data.AccountType,
	).Scan(&openedAccountId)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return &t.OpenedAccounts{
		OpenedAccountID:     openedAccountId,
		AccountType:         data.AccountType,
		Balance:             data.Balance,
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
