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

	"github.com/jackc/pgx/v5/pgxpool"
)

type BankCardNumberGenerated struct {
	BankCardNumber string    `json:"card_number"`
	ExpiryDate     time.Time `json:"expiry_date"`
	AccountId      string    `json:"account_id"`
	CardType       string    `json:"card_type"`
}

type BankCardNumberResponse struct {
	BankCardNumber string    `json:"card_number"`
	PinNumber      string    `json:"pin_number"`
	ExpiryDate     time.Time `json:"expiry_date"`
	AccountId      string    `json:"account_id"`
	CardType       string    `json:"card_type"`
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

func GenerateBankCardCVV() string {
	b := make([]byte, 1)
	_, err := rand.Read(b)
	if err != nil {
		return "000"
	}
	return fmt.Sprintf("%03d", 100+int(b[0])%900) // 100–999
}

// Bank Card Creation
func (b *BankcardService) CreateBankCardForUser(ctx context.Context, cardType string, accountId string, first_name string, last_name string) (*BankCardNumberGenerated, error) {
	var res BankCardNumberGenerated

	cardNumber := GenerateBankCardNumber(first_name, last_name, cardType)
	cardPin := GenerateBankCardPinNumber()
	cvv := GenerateBankCardCVV()

	var expiryDate time.Time
	switch cardType {
	case "debit":
		expiryDate = time.Now().AddDate(8, 0, 0)
	case "credit":
		expiryDate = time.Now().AddDate(12, 0, 0)
	default:
		expiryDate = time.Now().AddDate(5, 0, 0)
	}

	_, err := b.db.Exec(ctx, `
		INSERT INTO BankCard (card_number, expiry_date, account_id, cvv, pin_number, card_type)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, cardNumber, expiryDate, accountId, cvv, cardPin, cardType)

	if err != nil {
		return nil, fmt.Errorf("failed to insert into BankCard: %w", err)
	}

	res.BankCardNumber = cardNumber
	res.ExpiryDate = expiryDate
	res.AccountId = accountId
	res.CardType = cardType

	return &res, nil
}

func (b *BankcardService) GetBankCardByNumber(ctx context.Context, cardNumber string) (*BankCardNumberResponse, error) {
	var res BankCardNumberResponse
	err := b.db.QueryRow(ctx, `
        SELECT 
            card_number, 
            pin_number, 
            expiry_date, 
            account_id, 
            card_type
        FROM BankCard
        WHERE card_number = $1
    `, cardNumber).Scan(
		&res.BankCardNumber,
		&res.PinNumber,
		&res.ExpiryDate,
		&res.AccountId,
		&res.CardType,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get bank card: %w", err)
	}

	return &res, nil
}
