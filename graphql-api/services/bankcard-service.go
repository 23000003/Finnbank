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

type BankCardResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CardType  string `json:"card_type"`
	Status    string `json: "status"`
}

type BankCardRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CardType  string `json:"card_type"`
}

type BankCardNumberGenerated struct {
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

func GenrateBankCardNumber(first_name string, last_name string, card_type int) string {
	// Step 1: Combine key strings
	combined := first_name + last_name + fmt.Sprintf("%d", card_type)

	// Step 2: Hash the combined string using SHA-1 (or any other hashing algo)
	hasher := sha256.Sum256([]byte(combined))

	num := binary.BigEndian.Uint64(hasher[:8]) // Get the first 8 bytes of the hash

	return fmt.Sprintf("%012d", num%100000000000000000+9000000000000000) // Ensure it starts with 4 (Visa) and is 16 digits long
}

func (b *BankcardService) GetBankCardRequestsById(ctx context.Context, id int) (*BankCardResponse, error) {
	var res BankCardResponse
	_, err := b.db, b.db.QueryRow(ctx, `
		select crl.first_name, crl.last_name, ct.name as card_type, crl.status from card_request_list crl
	join card_types ct on crl.card_type = ct.card_type_id
	where crl.request_id = $1;
	`, id).Scan(
		&res.FirstName,
		&res.LastName,
		&res.CardType,
		&res.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank card request: %w", err)
	}
	return &res, nil
}

func (b *BankcardService) CreateCardRequest(ctx context.Context, first_name string, last_name string, card_type int) (*BankCardRequest, error) {
	var req BankCardRequest

	_, err := b.db.Exec(ctx, `
	INSERT INTO card_request_list (first_name, last_name, card_type) 
	VALUES ($1, $2, $3)`, first_name, last_name, card_type)

	if err != nil {
		return nil, fmt.Errorf("failed to insert date into the table: %w", err)
	}

	return &req, nil
}

func (b *BankcardService) CreateBankCardForUser(ctx context.Context) (*BankCardNumberGenerated, error) {
	var res BankCardNumberGenerated

	return &res, nil
}
