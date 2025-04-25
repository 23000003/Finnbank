package entities

import (
	ty "finnbank/graphql-api/types"

	"github.com/graphql-go/graphql"
)

func GetProductEntityType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Product",
			Fields: graphql.Fields{
				"id": 		&graphql.Field{ Type: graphql.Int	},
				"name": 	&graphql.Field{ Type: graphql.String },
				"info": 	&graphql.Field{ Type: graphql.String },
				"price": 	&graphql.Field{ Type: graphql.Float },
			},
		},
	)
}

func GetAccountEntityType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Account",
			Fields: graphql.Fields{
				"account_id": 		&graphql.Field{ Type: graphql.String },
				"email": 					&graphql.Field{ Type: graphql.String },
				"first_name": 		&graphql.Field{ Type: graphql.String },
				"middle_name": 		&graphql.Field{ Type: graphql.String },
				"last_name": 			&graphql.Field{ Type: graphql.String },
				"phone_number": 	&graphql.Field{ Type: graphql.String },
				"password": 			&graphql.Field{ Type: graphql.String },
				"date_created": 	&graphql.Field{ Type: graphql.DateTime },
				"date_updated": 	&graphql.Field{ Type: graphql.DateTime },
				"account_number": &graphql.Field{ Type: graphql.String },
				"has_card": 			&graphql.Field{ Type: graphql.Boolean },
				"address": 				&graphql.Field{ Type: graphql.String },
				"account_type": 	&graphql.Field{ Type: graphql.String },
				"nationality": 		&graphql.Field{ Type: graphql.String },
				"auth_id": 				&graphql.Field{ Type: graphql.String },
				"birthdate": 			&graphql.Field{ Type: graphql.DateTime },
				"national_id": 		&graphql.Field{ Type: graphql.String },
				"account_status": &graphql.Field{ Type: graphql.String },
			},
		},
	)
}

var bankCardEntityType = graphql.NewObject(graphql.ObjectConfig{
	Name: "BankCardEntity",
	Fields: graphql.Fields{
		"bankcard_id": 	&graphql.Field{ Type: graphql.Int },
		"card_number": 	&graphql.Field{ Type: graphql.String },
		"expiry_date": 	&graphql.Field{ Type: graphql.DateTime },
		"account_id":		&graphql.Field{ Type: graphql.String },
		"cvv": 					&graphql.Field{ Type: graphql.String },
		"pin_number": 	&graphql.Field{ Type: graphql.String },
		"card_type": 		&graphql.Field{ Type: graphql.String },
		"date_created": &graphql.Field{ Type: graphql.DateTime },
	},
})


func GetBankCardEntity() *graphql.Object {
	return bankCardEntityType
}

var TransactionEntityType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Transaction",
	Fields: graphql.Fields{
		"transaction_id":     &graphql.Field{Type: graphql.Int},
		"ref_no":             &graphql.Field{Type: graphql.String},
		"sender_id":          &graphql.Field{Type: graphql.Int},
		"receiver_id":        &graphql.Field{Type: graphql.Int},
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
		// "ref_no":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"sender_id":        &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
		"receiver_id":      &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
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

// Experimenting
var notifTypeEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "NotificationType",
	Values: graphql.EnumValueConfigMap{
		"Transaction": &graphql.EnumValueConfig{Value: "Transaction"},
		"System":      &graphql.EnumValueConfig{Value: "System"},
	},
})

var notificationEntityType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Notification",
	Fields: graphql.Fields{
		"notif_id":        &graphql.Field{Type: graphql.String},
		"notif_type":      &graphql.Field{Type: notifTypeEnum},
		"notif_to_id":     &graphql.Field{Type: graphql.String},
		"notif_from_name": &graphql.Field{Type: graphql.String},
		"content":         &graphql.Field{Type: graphql.String},
		"is_read":         &graphql.Field{Type: graphql.Boolean},
		"redirect_url":    &graphql.Field{Type: graphql.String},
		"date_notified":   &graphql.Field{Type: graphql.DateTime},
		"date_read":       &graphql.Field{Type: graphql.DateTime},
	},
})

func GetNotificationEntityType() *graphql.Object {
	return notificationEntityType
}

func GetOpenedAccountEntityType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "OpenedAccount",
			Fields: graphql.Fields{
				"openedaccount_id":     &graphql.Field{Type: graphql.Int},
				"account_id":           &graphql.Field{Type: graphql.Int},
				"bankcard_id":          &graphql.Field{Type: graphql.Int},
				"account_type":         &graphql.Field{Type: graphql.String},
				"balance":              &graphql.Field{Type: graphql.Float},
				"openedaccount_status": &graphql.Field{Type: graphql.String},
				"date_created":         &graphql.Field{Type: graphql.DateTime},
				"account_number":     &graphql.Field{Type: graphql.String},
			},
		},
	)
}
