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
	"math/rand"
	"time"

	"github.com/google/uuid"
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
	BankCardNumber string `json: "bankcard_number"`
	PinNumber      string `json: "bankcard_pin"`
	ExpiryDate     string `json: "expiry_date"`
	AccountId      string `json: "account_id"`
	CardType       string `json: "card_type"`
}

type BankCardNumberResponse struct {
	BankCardNumber string `json: "bankcard_number"`
	PinNumber      string `json: "bankcard_pin"`
	ExpiryDate     string `json: "expiry_date"`
	AccountId      string `json: "account_id"`
	CardType       string `json: "card_type"`
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

func init() {
	rand.Seed(time.Now().UnixNano())

}

func GenerateBankCardNumber(first_name string, last_name string, card_type string) string {
	// Step 1: Determine prefix based on card type
	var prefix string
	switch card_type {
	case "debit":
		prefix = "51"
	case "credit":
		prefix = "52"
	case "prepaid":
		prefix = "53"
	default:
		prefix = "50" // Fallback for unknown types
	}

	// Step 2: Combine name and type for hashing
	combined := first_name + last_name + card_type

	// Step 3: Hash and extract number
	hasher := sha256.Sum256([]byte(combined))
	num := binary.BigEndian.Uint64(hasher[:8]) % 100000000000 // 11 digits

	// Step 4: Return 13-digit card number with prefix
	return fmt.Sprintf("%s%011d", prefix, num)
}

func GenerateBankCardPinNumber() string {
	// Generate a random 4-digit PIN number
	pin := make([]byte, 4)
	for i := range pin {
		pin[i] = byte('0' + rand.Intn(10)) // Random digit between '0' and '9'
	}
	return string(pin)
}

// Bank Card Requests
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

func (b *BankcardService) CreateCardRequest(ctx context.Context, first_name string, last_name string, card_type string) (*BankCardRequest, error) {
	var req BankCardRequest

	_, err := b.db.Exec(ctx, `
	INSERT INTO card_request_list (first_name, last_name, card_type) 
	VALUES ($1, $2, $3)`, first_name, last_name, card_type)

	if err != nil {
		return nil, fmt.Errorf("failed to insert data into the table: %w", err)
	}

	return &req, nil
}

// Bank Card Creation
func (b *BankcardService) CreateBankCardForUser(ctx context.Context, fname string, lname string, cardtype string, account_holder_id uuid.UUID) (*BankCardNumberGenerated, error) {
	var res BankCardNumberGenerated

	var card_number = GenerateBankCardNumber(fname, lname, cardtype)

	var card_pin = GenerateBankCardPinNumber()

	expiryDate := time.Now().AddDate(5, 0, 0).Format("2006-01-02") // "YYYY-MM-DD"

	_, err := b.db.Exec(ctx, `
		INSERT INTO bankcard_list (bankcard_number, bankcard_pin, expiry_date, account_id, card_type)
		VALUES ($1, $2, $3, $4, $5)
	`, card_number, card_pin, expiryDate, account_holder_id, cardtype)

	if err != nil {
		return nil, fmt.Errorf("failed to insert data into the table: %w", err)
	}

	return &res, nil
}
