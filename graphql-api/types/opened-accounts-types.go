package types

import "time"

type OpenedAccounts struct {
	OpenedAccountID     int       `json:"openedaccount_id"`
	BankCardID          *int      `json:"bankcard_id"`
	Balance             float64   `json:"balance"`
	AccountType         string    `json:"account_type"`
	DateCreated         time.Time `json:"date_created"`
	OpenedAccountStatus string    `json:"openedaccount_status"`
	AccountNumber          string    `json:"account_number"`
}
