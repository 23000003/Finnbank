package types

import (
	"github.com/graphql-go/handler"
)

type IGraphQLHandlers interface {
	ProductServicesHandler(connAddress string) *handler.Handler
	AccountServicesHandler() *handler.Handler
	BankCardServicesHandler(connAddress string) *handler.Handler
	NotificationServicesHandler(connAddress string) *handler.Handler
	StatementServicesHandler(connAddress string) *handler.Handler
	TransactionServicesHandler(connAddress string) *handler.Handler
}
