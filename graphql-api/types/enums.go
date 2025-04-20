package types

import "github.com/graphql-go/graphql"

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
