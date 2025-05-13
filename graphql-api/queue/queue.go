package queue

import (
	"context"
	"finnbank/common/utils"
	h "finnbank/graphql-api/helpers"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

// **********************************************************************
// This is to simulate from Pending to Completed / Failed in Transaction
// If the transaction is completed, open account balance is deducted
// **********************************************************************

type Queue struct {
	transacId	 []int // To update status
	senderId	 []int // To deduct
	receiverId []int
	openAccAmount []float64 // Num to deduct
	mu    sync.Mutex
	l 		*utils.Logger
	ctx  	context.Context
}


func NewQueue(l *utils.Logger, c context.Context) *Queue {
	return &Queue{
		transacId: []int{},
		senderId: []int{},
		receiverId: []int{},
		openAccAmount: []float64{},
		mu: sync.Mutex{},
		l: l,
		ctx: c,
	}
}

func (q *Queue) Enqueue(transacId int, senderId int, receiverId int, openAccAmount float64) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.transacId = append(q.transacId, transacId)
	q.senderId = append(q.senderId, senderId)
	q.receiverId = append(q.receiverId, receiverId)
	q.openAccAmount = append(q.openAccAmount, openAccAmount)
	
	q.l.Info("Enqueued:")
	q.l.Info(" - Transaction ID: %d", transacId)
	q.l.Info(" - Sender ID: %d", senderId)
	q.l.Info(" - Receiver ID: %d", receiverId)
	q.l.Info(" - Open Account Balance: %f", openAccAmount)
}

func (q *Queue) Dequeue() (int, int, int, float64, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.openAccAmount) == 0 || len(q.receiverId) == 0 || len(q.senderId) == 0 || len(q.transacId) == 0 {
		return 0, 0, 0, 0, false
	}

	transacId := q.transacId[0]
	senderId := q.senderId[0]
	receiverId := q.receiverId[0]
	openAccAmount := q.openAccAmount[0]

	q.openAccAmount = q.openAccAmount[1:]	
	q.transacId = q.transacId[1:]
	q.receiverId = q.receiverId[1:]
	q.senderId = q.senderId[1:]

	return transacId, senderId, receiverId, openAccAmount, true
}

// StartAutoDequeue starts a ticker that dequeues every 1 minute
func (q *Queue) StartAutoDequeue(OADBPool *pgxpool.Pool, TXDBPool *pgxpool.Pool, ACCDBPool *pgxpool.Pool, transacConn *websocket.Conn) {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-q.ctx.Done():
				q.l.Info("AutoDequeue stopped due to context cancellation.")
				return
			case <-ticker.C:
				transacId, senderId, receiverId, openAccAmount, ok := q.Dequeue()
				if ok {
					q.l.Info("Dequeued:")
					err := h.SuccessTransaction(q.ctx, OADBPool, TXDBPool, ACCDBPool, transacId, receiverId)
					if err != nil {
						_ = h.FailedTransaction(q.ctx, TXDBPool, transacId)
						q.l.Error("Failed to update success transaction: %v", err)
						err := h.SendTransactionUpdate(q.ctx, transacConn, transacId, "FAILED")
						if err != nil {
							q.l.Error("Failed to send transaction update: %v", err)
						}
					} else {
						err1 := h.DeductOpenedAccountBalance(q.ctx, OADBPool, senderId, receiverId, openAccAmount)
						err2 := h.SendTransactionUpdate(q.ctx, transacConn, transacId, "COMPLETED")
						if err1 != nil {
							q.l.Error("Failed to deduct opened account balance: %v", err1)
						}
						if err2 != nil {
							q.l.Error("Failed to send transaction update: %v", err2)
						}
					}
				} else {
					q.l.Info("Queue is empty.")
				}
			}
		}
	}()
}
