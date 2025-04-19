package services

// Use this for resolvers business logic

// GetBankCardOfUserById, (Query)
// CreateBankCardForUser,  (Mutation)
// UpdateBankcardExpiryDateByUserId  (Mutation)

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"finnbank/common/utils"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AllBankCardRequests struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CardType  string `json:"card_type"`
}

type Requester struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CardType  string `json:"card_type"`
	BirthDate string `json:"birth_date"`
}

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

func GenrateBankCardNumber(req Requester) string {
	// Step 1: Combine key strings
	combined := req.FirstName + req.LastName + req.CardType + req.BirthDate

	// Step 2: Hash the combined string using SHA-1 (or any other hashing algo)
	hasher := sha256.Sum256([]byte(combined))

	num := binary.BigEndian.Uint64(hasher[:8]) // Get the first 8 bytes of the hash

	return fmt.Sprintf("%012d", num%100000000000000000+9000000000000000) // Ensure it starts with 4 (Visa) and is 16 digits long
}

func (b *BankcardService) CreateCardRequest(ctx context.Context, req Requester) error {
	conn, err := b.db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `
        INSERT INTO card_request (first_name, last_name, card_type)
        VALUES ($1, $2, $3)
    `, req.FirstName, req.LastName, req.CardType)

	if err != nil {
		return fmt.Errorf("failed to insert card request: %w", err)
	}

	return nil
}

func (b *BankcardService) GetAllBankCardRequestsById(ctx context.Context, id int) ([]AllBankCardRequests, error) {

	conn, err := b.db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(ctx,
		`select first_name, last_name, card_type from card_request where request_id = $1`,
		id,
	)
	defer rows.Close()

	var results []AllBankCardRequests
	for rows.Next() {
		var bcn AllBankCardRequests
		if err := rows.Scan(
			&bcn.FirstName,
			&bcn.LastName,
			&bcn.CardType,
		); err != nil {
			return nil, err
		}
		results = append(results, bcn)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (b *BankcardService) CreateBankCardForUser(ctx context.Context, req Requester) (string, error) {
	conn, err := b.db.Acquire(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	cardNumber := GenrateBankCardNumber(req)

	_, err = conn.Exec(ctx, `
		INSERT INTO bank_card (card_number, expiry_date)
		VALUES ($1, $2)
	`, cardNumber, req.BirthDate)

	if err != nil {
		return "", fmt.Errorf("failed to insert bank card: %w", err)
	}

	return cardNumber, nil
}
