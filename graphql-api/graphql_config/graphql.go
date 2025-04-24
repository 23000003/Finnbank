package graphql_config

/**
 Handles all the graphql requests and Configure Routes
**/

import (
	"finnbank/common/utils"
	"finnbank/graphql-api/types"
	"net/http"
	"github.com/99designs/gqlgen/graphql/playground"
	q "finnbank/graphql-api/queue"
)

type StructGraphQL struct {
	h   types.IGraphQLHandlers
	s   types.StructGrpcServiceConnections
	log *utils.Logger
}

// Keeping Grpc future purposes
func NewGraphQL(log *utils.Logger, h types.IGraphQLHandlers) *StructGraphQL {
	return &StructGraphQL{
		log: log,
		h:   h,
		s: types.StructGrpcServiceConnections{
			ProductServer:       ":9000",
			BankCardServer:      ":9001",
			AuthServer:          ":9002",
			StatementServer:     ":9004",
			TransactionServer:   ":9005",
			NotificationServer:  ":9006",
			OpenedAccountServer: ":9007",
		},
	}
}

func (gql *StructGraphQL) ConfigureGraphQLHandlers(q *q.Queue) {

	gql.log.Info("Configuring GraphQL Handlers...")

	productHandler := gql.h.ProductServicesHandler(gql.s.ProductServer)
	accountHandler := gql.h.AccountServicesHandler(gql.s.AuthServer)
	bankCardHandler := gql.h.BankCardServicesHandler(gql.s.BankCardServer)
	statementHandler := gql.h.StatementServicesHandler(gql.s.StatementServer)
	transactionHandler := gql.h.TransactionServicesHandler(gql.s.TransactionServer, q)
	notificationHandler := gql.h.NotificationServicesHandler(gql.s.NotificationServer)
	openedAccountHandler := gql.h.OpenedAccountServicesHandler(gql.s.OpenedAccountServer)

	http.Handle("/graphql/account", accountHandler)
	http.Handle("/graphql/product", productHandler)
	http.Handle("/graphql/bankcard", bankCardHandler)
	http.Handle("/graphql/statement", statementHandler)
	http.Handle("/graphql/transaction", transactionHandler)
	http.Handle("/graphql/notification", notificationHandler)
	http.Handle("/graphql/opened-account", openedAccountHandler)

	// For graphql query/mutation Testing
	http.Handle("/playground/account", playground.Handler("Account GraphQL Playground", "/graphql/account"))
	http.Handle("/playground/product", playground.Handler("Product GraphQL Playground", "/graphql/product"))
	http.Handle("/playground/bankcard", playground.Handler("BankCard GraphQL Playground", "/graphql/bankcard"))
	http.Handle("/playground/statement", playground.Handler("Statement GraphQL Playground", "/graphql/statement"))
	http.Handle("/playground/transaction", playground.Handler("Transaction GraphQL Playground", "/graphql/transaction"))
	http.Handle("/playground/notification", playground.Handler("Notification GraphQL Playground", "/graphql/notification"))
	http.Handle("/playground/opened-account", playground.Handler("OpenedAccount GraphQL Playground", "/graphql/openedaccount"))

	gql.log.Info("Configured GraphQL Handlers")

}
