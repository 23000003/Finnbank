package entities

import (
	ty "finnbank/graphql-api/types"
	"time"

	"github.com/graphql-go/graphql"
)

func GetProductEntityType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Product",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"info": &graphql.Field{
					Type: graphql.String,
				},
				"price": &graphql.Field{
					Type: graphql.Float,
				},
			},
		},
	)
}

func GetAccountEntityType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Account",
			Fields: graphql.Fields{
				"account_id": &graphql.Field{
					Type: graphql.String, // UUID in DB
				},
				"email": &graphql.Field{
					Type: graphql.String,
				},
				"full_name": &graphql.Field{
					Type: graphql.String,
				},
				"phone_number": &graphql.Field{
					Type: graphql.String,
				},
				"password": &graphql.Field{
					Type: graphql.String, // encrypted
				},
				"date_created": &graphql.Field{
					Type: graphql.DateTime,
				},
				"date_updated": &graphql.Field{
					Type: graphql.DateTime,
				},
				"account_number": &graphql.Field{
					Type: graphql.String,
				},
				"has_card": &graphql.Field{
					Type: graphql.Boolean,
				},
				"address": &graphql.Field{
					Type: graphql.String,
				},
				"balance": &graphql.Field{
					Type: graphql.Float,
				},
				"account_type": &graphql.Field{
					Type: graphql.String,
				},
				"nationality": &graphql.Field{
					Type: graphql.String,
				},
				"auth_id": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)
}

func GetBankCardEntityType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "BankCard",
			Fields: graphql.Fields{
				"card_number": &graphql.Field{
					Type: graphql.String, // UUID in DB
				},
				"expiry": &graphql.Field{
					Type: graphql.DateTime,
				},
				"account_id": &graphql.Field{
					Type: graphql.Int,
				},
				"cvv": &graphql.Field{
					Type: graphql.Int,
				},
				"pin_number": &graphql.Field{
					Type: graphql.String, // encrypted
				},
				"date_created": &graphql.Field{
					Type: graphql.DateTime,
				},
				"card_type": &graphql.Field{
					Type: graphql.EnumValueType,
				},
			},
		},
	)
}

//	func GetTransactionEntityType() *graphql.Object {
//		return graphql.NewObject(
//			graphql.ObjectConfig{
//				Name: "Transaction",
//				Fields: graphql.Fields{
//					"ref_number": &graphql.Field{
//						Type: graphql.String, // UUID in DB
//					},
//					"sender": &graphql.Field{
//						Type: graphql.Int,
//					},
//					"receiver": &graphql.Field{
//						Type: graphql.Int,
//					},
//					"transaction_type": &graphql.Field{
//						Type: graphql.EnumValueType,
//					},
//					"amount": &graphql.Field{
//						Type: graphql.Int,
//					},
//					"date_created": &graphql.Field{
//						Type: graphql.DateTime,
//					},
//				},
//			},
//		)
//	}
//
// ① Declare it once at package‐init time:
var TransactionEntityType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Transaction",
	Fields: graphql.Fields{
		"transaction_id":     &graphql.Field{Type: graphql.String},
		"ref_no":             &graphql.Field{Type: graphql.String},
		"sender_id":          &graphql.Field{Type: graphql.String},
		"receiver_id":        &graphql.Field{Type: graphql.String},
		"transaction_type":   &graphql.Field{Type: ty.TransactionTypeEnum},
		"amount":             &graphql.Field{Type: graphql.Float},
		"transaction_status": &graphql.Field{Type: ty.TransactionStatusEnum},
		"date_transaction":   &graphql.Field{Type: graphql.DateTime},
		"transaction_fee":    &graphql.Field{Type: graphql.Float},
		"notes":              &graphql.Field{Type: graphql.String},
	},
})

// ② And similarly create your InputObject only once:
var TransactionInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "TransactionInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"ref_no":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"sender_id":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"receiver_id":      &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"transaction_type": &graphql.InputObjectFieldConfig{Type: ty.TransactionTypeEnum},
		"amount":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
		"transaction_fee":  &graphql.InputObjectFieldConfig{Type: graphql.Float},
		"notes":            &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

// ③ If you still want getters:
func GetTransactionEntityType() *graphql.Object {
	return TransactionEntityType
}

func GetTransactionInputType() *graphql.InputObject {
	return TransactionInputType
}

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

func GetNotificationEntityType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Notification",
			Fields: graphql.Fields{
				"notif_id": &graphql.Field{
					Type: graphql.Int,
				},
				"notif_type": &graphql.Field{
					Type: graphql.EnumValueType,
				},
				"account_id": &graphql.Field{
					Type: graphql.Int,
				},
				"redirect_url": &graphql.Field{
					Type: graphql.String,
				},
				"date_notified": &graphql.Field{
					Type: graphql.DateTime,
				},
			},
		},
	)
}

func GetOpenedAccountEntityType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "OpenedAccount",
			Fields: graphql.Fields{
				"openedaccount_id": &graphql.Field{
					Type: graphql.Int,
				},
				"account_id": &graphql.Field{
					Type: graphql.Int,
				},
				"bankcard_id": &graphql.Field{
					Type: graphql.Int,
				},
				"account_type": &graphql.Field{
					Type: graphql.String,
				},
				"balance": &graphql.Field{
					Type: graphql.Float,
				},
				"openedaccount_status": &graphql.Field{
					Type: graphql.String,
				},
				"date_created": &graphql.Field{
					Type: graphql.DateTime,
				},
			},
		},
	)
}
