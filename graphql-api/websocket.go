package main

import (
	"context"
	"finnbank/common/utils"
	"finnbank/graphql-api/types"
	"net/http"
	"github.com/gorilla/websocket"
)

func initializeWebsockets(log *utils.Logger, ctx context.Context) *types.StructWebSocketConnections {
	
	headers := http.Header{}
	headers.Add("Origin", "http://localhost:8083")
	transacUrl := "ws://localhost:8080/api/ws/listen-to-transaction"
	notifUrl := "ws://localhost:8080/api/ws/listen-to-notification"

	transacConn, _, err := websocket.DefaultDialer.Dial(transacUrl, headers)
	notifConn, _, err1 := websocket.DefaultDialer.Dial(notifUrl, headers)

	if err != nil || err1 != nil {
		log.Fatal("Error connecting to websocket: %v %v", err, err1)
		return nil
	}

	// Start listening on both connections
	go listenAndHandle(ctx, transacConn, log, "transaction")
	go listenAndHandle(ctx, notifConn, log, "notification")

	return &types.StructWebSocketConnections{
		TransactionConn:  transacConn,
		NotificationConn: notifConn,
	}
}

func listenAndHandle(ctx context.Context, conn *websocket.Conn, log *utils.Logger, tag string) {
	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping WebSocket listener for %s", tag)
			return
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Info("WebSocket [%s] read error: %v", tag, err)
				return
			}
			log.Info("WebSocket [%s] received: %s", tag, string(msg))
		}
	}
}
