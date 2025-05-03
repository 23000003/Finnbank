package services

import (
	"context"
	"finnbank/common/utils"
	t "finnbank/graphql-api/types"
	"fmt"
	"time"

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
			account_type, openedaccount_status, date_created, account_number
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
			&acc.AccountNumber,
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

func (s *OpenedAccountService) GetOpenedAccountIdByAccountNumber(ctx context.Context, account_num string) (int, error) {

	query := `
		SELECT 
			openedaccount_id
		FROM openedaccount 
		WHERE account_number = $1
	`

	var openedAccountId int
	err := s.db.QueryRow(ctx, query, account_num).Scan(&openedAccountId)
	if err != nil {
		return -1, fmt.Errorf("failed to fetch opened account ID: %w", err)
	}
	
	return openedAccountId, nil
}

func (s *OpenedAccountService) GetBothAccountNumberForReceipt(ctx context.Context, sent_id int, receive_id int) ([]*t.OpenedAccountNumber, error) {

	query := `
		SELECT 
			openedaccount_id, account_number
		FROM openedaccount 
		WHERE openedaccount_id = $1 OR openedaccount_id = $2
	`

	s.l.Info("Id: %d, %d", sent_id, receive_id)

	rows, err := s.db.Query(ctx, query, sent_id, receive_id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch opened account IDs: %w", err)
	}
	defer rows.Close()

	var results []*t.OpenedAccountNumber
	for rows.Next() {
		var acc t.OpenedAccountNumber
		if err := rows.Scan(&acc.OpenedAccountID, &acc.AccountNumber); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		results = append(results, &acc)
	}

	s.l.Info("Both account numbers: %v", results)
	
	return results, nil
}

func (s *OpenedAccountService) GetOpenedAccountById(ctx context.Context, id int) (*t.OpenedAccounts, error) {
	var acc t.OpenedAccounts

	query := `
		SELECT 
			openedaccount_id, bankcard_id, balance, account_id,
			account_type, openedaccount_status, date_created, account_number
		FROM openedaccount 
		WHERE openedaccount_id = $1
	`

	err := s.db.QueryRow(ctx, query, id).Scan(
		&acc.OpenedAccountID,
		&acc.BankCardID,
		&acc.Balance,
		&acc.AccountID,
		&acc.AccountType,
		&acc.OpenedAccountStatus,
		&acc.DateCreated,
		&acc.AccountNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch opened account: %w", err)
	}

	return &acc, nil
}

// Create a new opened account
func (s *OpenedAccountService) CreateOpenedAccount(ctx context.Context, BCService *BankcardService, user_id string) ([]*t.OpenedAccounts, error) {
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	var bankcardId []int = nil
	id, err := BCService.CreateCardRequest(ctx, user_id)
	s.l.Info("Bankcard ID: %d", id)
	if err != nil {
		return nil, fmt.Errorf("failed to create card request: %w", err)
	}
	bankcardId = id

	debit, err := generateRandomNumber(16)
	credit, err1 := generateRandomNumber(16)
	savings, err2 := generateRandomNumber(16)

	if err != nil || err1 != nil || err2 != nil {
		return nil, fmt.Errorf("failed to generate random numbers: %w", err)
	}

	var accounts []*t.OpenedAccounts
	rows, err := conn.Query(ctx,
		`INSERT INTO openedaccount (account_id, bankcard_id, balance, account_type, openedaccount_status, account_number) 
		 VALUES 
			($1, $2, $3, $4, $8, $10), 
			($1, $5, $3, $6, $8, $11),
			($1, NULL, $3, $7, $9, $12)
		RETURNING openedaccount_id, account_type, bankcard_id, balance, openedaccount_status`,
		user_id, bankcardId[0], 10000, "Credit",
		bankcardId[1], "Checking", "Savings", "Closed", "Active", credit, debit, savings,
	)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var acc t.OpenedAccounts
		var bankCardIDNullable *int
		if err := rows.Scan(&acc.OpenedAccountID, &acc.AccountType, &bankCardIDNullable, &acc.Balance, &acc.OpenedAccountStatus); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		acc.DateCreated = time.Now()
		acc.BankCardID = bankCardIDNullable
		accounts = append(accounts, &acc)
	}

	return accounts, nil
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
