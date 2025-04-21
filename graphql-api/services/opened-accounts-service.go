package services

import (
	"context"
	"finnbank/common/utils"
	"fmt"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AllOpenedAccounts struct {
	OpenedAccountID     int       `json:"openedaccount_id"`
	BankCardID          *int      `json:"bankcard_id"`
	Balance             float64   `json:"balance"`
	AccountType         string    `json:"account_type"`
	DateCreated         time.Time `json:"date_created"`
	OpenedAccountStatus string    `json:"openedaccount_status"`
}

type OpenedAccountService struct {
	db     *pgxpool.Pool
	l *utils.Logger
}

func NewOpenedAccountService(db *pgxpool.Pool, logger *utils.Logger) *OpenedAccountService {
	return &OpenedAccountService{
		db:     db,
		l: logger,
	}
}

func (s *OpenedAccountService) GetAllOpenedAccountsByUserId(ctx context.Context, id int) ([]AllOpenedAccounts, error) {
	
	// Get a connection from the pool
	conn, err := s.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, 
		`SELECT 
			openedaccount_id, bankcard_id, balance, 
			account_type, openedaccount_status, date_created 
		FROM 
			openedaccount WHERE account_id = $1`, 
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	// iterates over the rows and scans the data into the struct
	var results []AllOpenedAccounts
	for rows.Next() {
		var acc AllOpenedAccounts
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
		results = append(results, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	s.l.Info("All opened accounts: %v", results)

	return results, nil
}
