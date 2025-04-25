package types

import (
	"github.com/graphql-go/handler"
	q "finnbank/graphql-api/queue"
)

type IGraphQLHandlers interface {
	ProductServicesHandler(connAddress string) *handler.Handler
	AccountServicesHandler(connAddress string) *handler.Handler
	BankCardServicesHandler(connAddress string) *handler.Handler
	NotificationServicesHandler(connAddress string) *handler.Handler
	StatementServicesHandler(connAddress string) *handler.Handler
	TransactionServicesHandler(connAddress string, queue *q.Queue) *handler.Handler
	OpenedAccountServicesHandler(connAddress string) *handler.Handler
}
