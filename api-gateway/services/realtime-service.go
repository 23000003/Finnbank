package services

import (
	"encoding/json"
	"finnbank/api-gateway/types"
	"finnbank/common/utils"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type RealTimeService struct {
	l                   *utils.Logger
	transactionClients  map[*websocket.Conn]bool
	notificationClients map[*websocket.Conn]bool
	mutex               sync.RWMutex             // Protect clients maps from concurrent access
}

func newRealTimeService(l *utils.Logger) *RealTimeService {
	return &RealTimeService{
		l:                   l,
		transactionClients:  make(map[*websocket.Conn]bool),
		notificationClients: make(map[*websocket.Conn]bool),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Transaction WebSocket Handler
func (r *RealTimeService) GetRealTimeTransaction(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		r.l.Error("Error upgrading connection: %v", err)
		return
	}

	r.mutex.Lock()
	r.transactionClients[conn] = true
	r.mutex.Unlock()

	// Configure connection handlers
	conn.SetCloseHandler(func(code int, text string) error {
		r.l.Info("Transaction connection closed by client: %d - %s", code, text)
		r.removeTransactionClient(conn)
		return nil
	})

	// Handle connection
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				r.l.Error("Error reading transaction message: %v", err)
			}
			break
		}

		var transac types.GetRealTimeTransaction
		if err := json.Unmarshal(message, &transac); err != nil {
			r.l.Error("Error unmarshaling transaction: %v", err)
			continue
		}

		r.l.Info("Received transaction - ID: %d, Ref: %s", transac.TransactionID, transac.RefNo)

		// Create and broadcast response
		ack := map[string]string{
			"transaction_id":      strconv.Itoa(transac.TransactionID),
			"ref_no":              transac.RefNo,
			"sender_id":           strconv.Itoa(transac.SenderID),
			"receiver_id":         strconv.Itoa(transac.ReceiverID),
			"transaction_type":    transac.TransactionType,
			"amount":             strconv.FormatFloat(transac.Amount, 'f', 2, 64),
			"transaction_status": transac.TransactionStatus,
			"date_transaction":   transac.DateTransaction.Format(time.RFC3339),
			"transaction_fee":    strconv.FormatFloat(transac.TransactionFee, 'f', 2, 64),
			"notes":             transac.Notes,
		}

		ackBytes, _ := json.Marshal(ack)
		r.broadcastToTransactionClients(ackBytes)
	}
}

// Notification WebSocket Handler
func (r *RealTimeService) GetRealTimeNotification(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		r.l.Error("Error upgrading notification connection: %v", err)
		return
	}

	r.mutex.Lock()
	r.notificationClients[conn] = true
	r.mutex.Unlock()

	// Configure connection handlers
	conn.SetCloseHandler(func(code int, text string) error {
		r.l.Info("Notification connection closed by client: %d - %s", code, text)
		r.removeNotificationClient(conn)
		return nil
	})

	// Handle connection
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				r.l.Error("Error reading notification message: %v", err)
			}
			break
		}

		var notif types.GetRealTimeNotification
		if err := json.Unmarshal(message, &notif); err != nil {
			r.l.Error("Error unmarshaling notification: %v", err)
			continue
		}

		r.l.Info("Received notification - ID: %d, Content: %s", notif.NotifID, notif.Content)

		// Create and broadcast response
		ack := map[string]string{
			"notif_id":        strconv.Itoa(notif.NotifID),
			"notif_type":      notif.NotifType,
			"notif_from_name": notif.NotifFromName,
			"content":         notif.Content,
			"is_read":         strconv.FormatBool(notif.IsRead),
			"date_notified":   notif.DateNotified.Format(time.RFC3339),
		}

		ackBytes, _ := json.Marshal(ack)
		r.broadcastToNotificationClients(ackBytes)
	}
}

// Broadcast to transaction clients
func (r *RealTimeService) broadcastToTransactionClients(message []byte) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for client := range r.transactionClients {
		// Timeout connection after 10 seconds since its not global
		client.SetWriteDeadline(time.Now().Add(10 * time.Second)) 
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			r.l.Error("Error writing to transaction client: %v", err)
			go r.removeTransactionClient(client)
		}
	}
}

// Broadcast to notification clients
func (r *RealTimeService) broadcastToNotificationClients(message []byte) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for client := range r.notificationClients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			r.l.Error("Error writing to notification client: %v", err)
			go r.removeNotificationClient(client)
		}
	}
}

// Remove transaction client
func (r *RealTimeService) removeTransactionClient(conn *websocket.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.transactionClients[conn]; ok {
		conn.Close()
		delete(r.transactionClients, conn)
		r.l.Info("Transaction client connection removed")
	}
}

// Remove notification client
func (r *RealTimeService) removeNotificationClient(conn *websocket.Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.notificationClients[conn]; ok {
		conn.Close()
		delete(r.notificationClients, conn)
		r.l.Info("Notification client connection removed")
	}
}

// Shutdown all WebSocket connections
func (r *RealTimeService) Shutdown() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Close transaction connections
	for client := range r.transactionClients {
		client.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server shutting down"),
			time.Now().Add(5*time.Second),
		)
		client.Close()
		delete(r.transactionClients, client)
	}

	// Close notification connections
	for client := range r.notificationClients {
		client.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server shutting down"),
			time.Now().Add(5*time.Second),
		)
		client.Close()
		delete(r.notificationClients, client)
	}

	r.l.Info("All WebSocket connections closed (transactions: %d, notifications: %d)", 
		len(r.transactionClients), len(r.notificationClients))
}