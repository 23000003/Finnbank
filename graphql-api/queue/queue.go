package queue

import (
	"finnbank/common/utils"
	"fmt"
	"sync"
	"time"
	h "finnbank/graphql-api/helpers"
	"github.com/jackc/pgx/v5/pgxpool"
	"context"
)

// **********************************************************************
// This is to simulate from Pending to Completed / Failed in Transaction
// If the transaction is completed, open account balance is deducted
// **********************************************************************

type Queue struct {
	transacId	 []string
	openAccId	 []string
	openAccAmount []float64
	mu    sync.Mutex
	l 		*utils.Logger
	ctx  	context.Context
}


func NewQueue(l *utils.Logger, c context.Context) *Queue {
	return &Queue{
		transacId: []string{},
		openAccId: []string{},
		openAccAmount: []float64{},
		mu: sync.Mutex{},
		l: l,
		ctx: c,
	}
}

func (q *Queue) Enqueue(transacId string, openAccId string, openAccAmount float64) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.transacId = append(q.transacId, transacId)

	q.openAccId = append(q.openAccId, openAccId)
	q.openAccAmount = append(q.openAccAmount, openAccAmount)
	
	q.l.Info("Enqueued:")
	q.l.Info(" - Transaction ID: %s", transacId)
	q.l.Info(" - Open Account ID: %s", openAccId)
	q.l.Info(" - Open Account Balance: %f", openAccAmount)
}

func (q *Queue) Dequeue() (string, string, float64, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.openAccAmount) == 0 || len(q.openAccId) == 0 || len(q.transacId) == 0 {
		return "", "", 0, false
	}

	transacId := q.transacId[0]
	openAccId := q.openAccId[0]
	openAccAmount := q.openAccAmount[0]

	q.openAccAmount = q.openAccAmount[1:]	
	q.transacId = q.transacId[1:]
	q.openAccId = q.openAccId[1:]

	return transacId, openAccId, openAccAmount, true
}

// StartAutoDequeue starts a ticker that dequeues every 1 minute
func (q *Queue) StartAutoDequeue(OADBPool *pgxpool.Pool, TXDBPool *pgxpool.Pool) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			transacId, openAccId, openAccAmount, ok := q.Dequeue()
			if ok {
				err := h.SuccessTransaction(q.ctx, TXDBPool, transacId)
				if err != nil {
					q.l.Error("Failed to update transaction status: %v", err)
					h.FailedTransaction(q.ctx, TXDBPool, transacId)
					break
				}
				err1 := h.DeductOpenedAccountBalance(q.ctx, OADBPool, openAccId, openAccAmount)
				if err1 != nil {
					q.l.Error("Failed to deduct opened account balance: %v", err1)
				}
			} else {
				fmt.Println("Queue is empty.")
			}
		}
	}()
}
