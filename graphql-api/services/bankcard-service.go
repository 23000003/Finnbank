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
			card_type, cvv, date_created, pin_number
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
			&bc.PinNumber,
		); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		if bc.PinNumber != "1234" {
			bc.PinNumber = "****"
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

	defaultPin := "1234"
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

func (b *BankcardService) VerifyBankcardPinNumber(ctx context.Context, bankcard_id int, pin_number string) (bool, error) {
	conn, err := b.db.Acquire(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	var hashedPin string

	err = conn.QueryRow(ctx,
		`SELECT pin_number 
		 FROM bankcard 
		 WHERE bankcard_id = $1`,
		bankcard_id).Scan(&hashedPin)

	if err != nil {
		return false, fmt.Errorf("query failed: %w", err)
	}

	b.l.Info("Hashed pin: %s", hashedPin)
	b.l.Info("Pin number: %s", pin_number)

	// this assumes hashedPin is "1234"
	if hashedPin == pin_number {
		return true, nil
	} 

	return verifyPassword(hashedPin, pin_number), nil
}

func (b *BankcardService) UpdateBankcardPinNumberById(ctx context.Context, bankcard_id int, pin_number string) (bool, error) {
	conn, err := b.db.Acquire(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	hashedPin, err := hashPassword(pin_number)
	if err != nil {
		return false, fmt.Errorf("failed to hash pin: %w", err)
	}

	var bankcardId string
	err = conn.QueryRow(ctx,
		`UPDATE bankcard SET pin_number = $2 
		 WHERE bankcard_id = $1 
		 RETURNING bankcard_id`,
		bankcard_id, hashedPin).Scan(&bankcardId)

	if err != nil {
		return false, fmt.Errorf("update failed: %w", err)
	}

	return true, nil
}
