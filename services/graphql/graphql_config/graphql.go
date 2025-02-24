package graphql_config

/**
 Handles all the graphql requests and Configure Routes
**/

import (
	"finnbank/services/common/utils"
	"finnbank/services/graphql/types"
	"net/http"
)


type StructGraphQL struct {
	h types.IGraphQLHandlers
	s types.StructGrpcServiceConnections
}

func NewGraphQL(h types.IGraphQLHandlers) *StructGraphQL {
	return &StructGraphQL{
		h: h,
		s: types.StructGrpcServiceConnections{
			ProductServer:      ":9000",
			BankCardServer:     ":9001",
			AccountServer:      ":9002",
			StatementServer:    ":9004",
			TransactionServer:  ":9005",
			NotificationServer: ":9006",
		},
	}
}

func (gql *StructGraphQL) ConfigureGraphQLHandlers(log *utils.Logger) {

	log.Info("Configuring GraphQL Handlers...")

	productHandler := gql.h.ProductServicesHandler(gql.s.ProductServer)
	accountHandler := gql.h.AccountServicesHandler(gql.s.AccountServer)
	bankCardHandler := gql.h.BankCardServicesHandler(gql.s.BankCardServer)
	statementHandler := gql.h.StatementServicesHandler(gql.s.StatementServer)
	transactionHandler := gql.h.TransactionServicesHandler(gql.s.TransactionServer)
	notificationHandler := gql.h.NotificationServicesHandler(gql.s.TransactionServer)

	http.Handle("/graphql/account", accountHandler)
	http.Handle("/graphql/product", productHandler)
	http.Handle("/graphql/bankcard", bankCardHandler)
	http.Handle("/graphql/statement", statementHandler)
	http.Handle("/graphql/transaction", transactionHandler)
	http.Handle("/graphql/notification", notificationHandler)

	log.Info("Configured GraphQL Handlers")
	
}
