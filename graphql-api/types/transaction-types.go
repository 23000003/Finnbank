package types

import (
	"time"

	"github.com/graphql-go/graphql"
)

type Transaction struct {
	TransactionID     string    `json:"transaction_id"`
	RefNo             string    `json:"ref_no"`
	SenderID          string    `json:"sender_id"`
	ReceiverID        string    `json:"receiver_id"`
	TransactionType   string    `json:"transaction_type"`
	Amount            float64   `json:"amount"`
	TransactionStatus string    `json:"transaction_status"`
	DateTransaction   time.Time `json:"date_transaction"`
	TransactionFee    float64   `json:"transaction_fee"`
	Notes             string    `json:"notes"`
}

// TransactionTypeEnum defines the allowed values for transaction types.
var TransactionTypeEnum = graphql.NewEnum(graphql.EnumConfig{
	Name:        "TransactionType",
	Description: "The type of the transaction (e.g., Transfer, Payment, Deposit, etc.).",
	Values: graphql.EnumValueConfigMap{
		"TRANSFER": &graphql.EnumValueConfig{
			Value:       "Transfer",
			Description: "A transfer transaction.",
		},
		"PAYMENT": &graphql.EnumValueConfig{
			Value:       "Payment",
			Description: "A payment transaction.",
		},
		"DEPOSIT": &graphql.EnumValueConfig{
			Value:       "Deposit",
			Description: "A deposit transaction.",
		},
		"WITHDRAW": &graphql.EnumValueConfig{
			Value:       "Withdraw",
			Description: "A withdrawal transaction.",
		},
		"LOAN": &graphql.EnumValueConfig{
			Value:       "Loan",
			Description: "A loan transaction.",
		},
	},
})

// TransactionStatusEnum defines the allowed values for transaction statuses.
var TransactionStatusEnum = graphql.NewEnum(graphql.EnumConfig{
	Name:        "TransactionStatus",
	Description: "The status of the transaction (e.g., Pending, Completed, Failed).",
	Values: graphql.EnumValueConfigMap{
		"PENDING": &graphql.EnumValueConfig{
			Value:       "Pending",
			Description: "The transaction is pending.",
		},
		"COMPLETED": &graphql.EnumValueConfig{
			Value:       "Completed",
			Description: "The transaction has been completed.",
		},
		"FAILED": &graphql.EnumValueConfig{
			Value:       "Failed",
			Description: "The transaction has failed.",
		},
	},
})
