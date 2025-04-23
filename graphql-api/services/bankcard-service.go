package services

// Use this for resolvers business logic

// GetBankCardOfUserById, (Query)
// CreateBankCardForUser,  (Mutation)
// UpdateBankcardExpiryDateByUserId  (Mutation)

import (
	"context"
	"finnbank/common/utils"
	t "finnbank/graphql-api/types"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BankcardService struct {
	db *pgxpool.Pool
	l  *utils.Logger
}

func NewBankcardService(db *pgxpool.Pool, logger *utils.Logger) *BankcardService {
	return &BankcardService{
		db: db,
		l:  logger,
	}
}

// Bank Card Requests
func (b *BankcardService) GetAllBankCardOfUserById(ctx context.Context, user_id string) ([]*t.GetBankCardResponse, error) {
	conn, err := b.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx,
		`SELECT 
			card_number, expiry_date, 
			card_type, cvv, date_created
		FROM bankcard 
		WHERE account_id = $1`,
		user_id,
	)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []*t.GetBankCardResponse
	for rows.Next() {
		var bc t.GetBankCardResponse
		if err := rows.Scan(
			&bc.CardNumber,
			&bc.ExpiryDate,
			&bc.CardType,
			&bc.CVV,
			&bc.DateCreated,
		); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		results = append(results, &bc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	b.l.Info("All opened accounts: %v", results)
	return results, nil
}

// called only in opened-account
func (b *BankcardService) CreateCardRequest(ctx context.Context, user_id string, card_type string, pin_number string) (string, error) {
	conn, err := b.db.Acquire(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	cardNumber, err := generateRandomNumber(16)
	if err != nil {
		return "", fmt.Errorf("failed to generate card number: %w", err)
	}
	cvv, err := generateRandomNumber(3)
	if err != nil {
		return "", fmt.Errorf("failed to generate cvv: %w", err)
	}
	expiryDate := time.Now().AddDate(4, 0, 0)

	var bankcard_id int
	err = conn.QueryRow(ctx,
		`INSERT INTO bankcard (account_id, card_number, card_type, pin_number, cvv, expiry_date) 
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING bankcard_id`,
			user_id, cardNumber, card_type, pin_number, cvv, expiryDate,
	).Scan(&bankcard_id)
	if err != nil {
		return "", fmt.Errorf("bankcard insert failed: %w", err)
	}

	return "Card created successfully", nil;
}

func (b *BankcardService) UpdateBankcardExpiryDateByUserId(ctx context.Context, bankcard_id int) (string, error) {
	conn, err := b.db.Acquire(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	var newExpiry string

	var bankcardId int
	err = conn.QueryRow(ctx,
		`UPDATE bankcard SET expiry_date = $1 
		 WHERE bankcard_id = $2 
		 RETURNING bankcard_id`,
		newExpiry, bankcard_id).Scan(&bankcardId)

	if err != nil {
		return "", fmt.Errorf("update failed: %w", err)
	}

	return "Update Successful", nil
}

