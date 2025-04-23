package types

import "time"

type BankCardResponse struct {
	BankCardId   string    `json:"bankcard_id"`
	CardNumber   string    `json:"card_number"`
	ExpiryDate   time.Time `json:"expiry_date"`
	CardType     string    `json:"card_type"`
	AccountId    string    `json:"account_id"`
	CVV          string    `json:"cvv"`
	PinNumber    string    `json:"pin_number"`
	DateCreated  time.Time `json:"date_created"`
}

type GetBankCardResponse struct {
	CardNumber   string    `json:"card_number"`
	ExpiryDate   time.Time `json:"expiry_date"`
	CardType     string    `json:"card_type"`
	CVV          string    `json:"cvv"`
	DateCreated  time.Time `json:"date_created"`
}


type BankCardRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CardType  string `json:"card_type"`
}

type BankCardNumberGenerated struct {
	BankCardNumber string `json:"bankcard_number"`
	PinNumber      string `json:"bankcard_pin"`
	ExpiryDate     string `json:"expiry_date"`
	AccountId      string `json:"account_id"`
	CardType       string `json:"card_type"`
}

type BankCardNumberResponse struct {
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	CardType       string    `json:"card_type"`
	BankCardNumber string    `json:"bankcard_number"`
	PinNumber      string    `json:"bankcard_pin"`
	ExpiryDate     time.Time `json:"expiry_date"`
	AccountId      string    `json:"account_id"`
}