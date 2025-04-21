package entities

import (
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
				"pin_number": &graphql.Field{
					Type: graphql.String,
				},
				"date_created": &graphql.Field{
					Type: graphql.DateTime,
				},
				"card_type_final": &graphql.Field{
					Type: graphql.EnumValueType,
				},
			},
		},
	)
}

var bankCardReponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "BankCardResponse",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"first_name": &graphql.Field{
			Type: graphql.String,
		},
		"last_name": &graphql.Field{
			Type: graphql.String,
		},
		"card_type": &graphql.Field{
			Type: graphql.String,
		},
		"status": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var bankCardRequestType = graphql.NewObject(graphql.ObjectConfig{
	Name: "BankCardRequest",
	Fields: graphql.Fields{
		"first_name": &graphql.Field{
			Type: graphql.String,
		},
		"last_name": &graphql.Field{
			Type: graphql.String,
		},
		"card_type": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func GetBankCardResponseEntity() *graphql.Object {
	return bankCardReponseType
}

func GetBankCardRequestEntity() *graphql.Object {
	return bankCardRequestType
}

var transactionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Transaction",
	Fields: graphql.Fields{
		"transaction_id":     &graphql.Field{Type: graphql.String},
		"ref_no":             &graphql.Field{Type: graphql.String},
		"sender_id":          &graphql.Field{Type: graphql.String},
		"receiver_id":        &graphql.Field{Type: graphql.String},
		"transaction_type":   &graphql.Field{Type: graphql.String},
		"amount":             &graphql.Field{Type: graphql.Int},
		"transaction_status": &graphql.Field{Type: graphql.String},
		"date_transaction":   &graphql.Field{Type: graphql.DateTime},
		"transaction_fee":    &graphql.Field{Type: graphql.Int},
		"notes":              &graphql.Field{Type: graphql.String},
	},
})

// 2️⃣ Getter just returns the shared var
func GetTransactionEntityType() *graphql.Object {
	return transactionType
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
