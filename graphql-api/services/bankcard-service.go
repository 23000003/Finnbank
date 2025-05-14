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
			bankcard_id, card_number, expiry_date, 
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
			&bc.BankCardId,
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
func (b *BankcardService) CreateCardRequest(ctx context.Context, user_id string) ([]int, error) {
	conn, err := b.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	debit, err := generateRandomNumber(16)
	credit, err1 := generateRandomNumber(16)
	if err != nil || err1 != nil {
		return nil, fmt.Errorf("failed to generate card number: %w %w", err, err1)
	}
	debitCVV, err := generateRandomNumber(3)
	creditCVV, err1 := generateRandomNumber(3)
	if err != nil || err1 != nil {
		return nil, fmt.Errorf("failed to generate cvv: %w %w", err, err1)
	}
	expiryDate := time.Now().AddDate(4, 0, 0)

	defaultPin := 1234
	var bankcardIDs []int
	rows, err := conn.Query(ctx,
		`INSERT INTO bankcard (account_id, card_number, card_type, pin_number, cvv, expiry_date) 
		 VALUES 
			($1, $2, $3, $4, $5, $6), 
			($1, $7, $8, $4, $9, $6)
		 RETURNING bankcard_id`,
		user_id, credit, "Credit", defaultPin, creditCVV, expiryDate, debit, "Debit", debitCVV,
	)
	if err != nil {
		return nil, fmt.Errorf("bankcard insert failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan bankcard_id: %w", err)
		}
		bankcardIDs = append(bankcardIDs, id)
	}
	
	return bankcardIDs, nil
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


func (b *BankcardService) UpdateBankcardPinNumberById(ctx context.Context, bankcard_id int, pin_number string) (string, error) {
	conn, err := b.db.Acquire(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	var bankcardId int
	err = conn.QueryRow(ctx,
		`UPDATE bankcard SET pin_number = $1 
		 WHERE bankcard_id = $2 
		 RETURNING bankcard_id`,
		bankcard_id, pin_number).Scan(&bankcardId)

	if err != nil {
		return "", fmt.Errorf("update failed: %w", err)
	}

	return "Update Successful", nil
}

