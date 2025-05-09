package types

import (
	q "finnbank/graphql-api/queue"

	"github.com/gorilla/websocket"
	"github.com/graphql-go/handler"
)

type IGraphQLHandlers interface {
	ProductServicesHandler(connAddress string) *handler.Handler
	AccountServicesHandler(connAddress string) *handler.Handler
	BankCardServicesHandler(connAddress string) *handler.Handler
	StatementServicesHandler(connAddress string) *handler.Handler
	OpenedAccountServicesHandler(connAddress string) *handler.Handler
	NotificationServicesHandler(connAddress string, notifConn *websocket.Conn) *handler.Handler
	TransactionServicesHandler(connAddress string, queue *q.Queue, transacConn *websocket.Conn) *handler.Handler
}
