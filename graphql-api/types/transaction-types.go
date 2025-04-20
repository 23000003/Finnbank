package types

import (
	"time"

	"github.com/graphql-go/graphql"
)

// Transaction represents a financial transaction.
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

// GetTransactionEntityType defines the GraphQL object type for a transaction entity.
func GetTransactionEntityType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Transaction",
		Fields: graphql.Fields{
			"transaction_id": &graphql.Field{
				Type:        graphql.String,
				Description: "The unique identifier for the transaction.",
			},
			"ref_no": &graphql.Field{
				Type:        graphql.String,
				Description: "The reference number for the transaction.",
			},
			"sender_id": &graphql.Field{
				Type:        graphql.String,
				Description: "The ID of the sender.",
			},
			"receiver_id": &graphql.Field{
				Type:        graphql.String,
				Description: "The ID of the receiver.",
			},
			"transaction_type": &graphql.Field{
				Type:        TransactionTypeEnum,
				Description: "The type of the transaction (e.g., Transfer, Payment, Deposit, etc.).",
			},
			"amount": &graphql.Field{
				Type:        graphql.Float,
				Description: "The amount of the transaction.",
			},
			"transaction_status": &graphql.Field{
				Type:        TransactionStatusEnum,
				Description: "The status of the transaction (e.g., Pending, Completed, Failed).",
			},
			"date_transaction": &graphql.Field{
				Type:        graphql.DateTime,
				Description: "The date and time when the transaction occurred.",
			},
			"transaction_fee": &graphql.Field{
				Type:        graphql.Float,
				Description: "The fee charged for the transaction.",
			},
			"notes": &graphql.Field{
				Type:        graphql.String,
				Description: "Additional notes or comments about the transaction.",
			},
		},
	})
}
