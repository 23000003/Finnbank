
package main

import (
	"fmt"
	"sync"
	"time"
)
// **********************************************************************
// This is to simulate from Pending to Completed / Failed in Transaction
// If the transaction is completed, open account balance is deducted
// **********************************************************************

type Queue struct {
	transacQuery []string
	openAccQuery []string
	mu    sync.Mutex
}

func (q *Queue) Enqueue(transacQuery string, openAccQuery string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.transacQuery = append(q.transacQuery, transacQuery)
	q.openAccQuery = append(q.openAccQuery, openAccQuery)
	fmt.Println("Enqueued:", transacQuery)
}

func (q *Queue) Dequeue() (string, string, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.transacQuery) == 0 || len(q.openAccQuery) == 0 {
		return "", "", false
	}

	transacQuery := q.transacQuery[0]
	openAccQuery := q.openAccQuery[0]

	q.transacQuery = q.transacQuery[1:]
	q.openAccQuery = q.openAccQuery[1:]

	return transacQuery, openAccQuery, true
}

// StartAutoDequeue starts a ticker that dequeues every 1 minute
func (q *Queue) StartAutoDequeue() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			transac, openAcc, ok := q.Dequeue()
			if ok {
				//  Will call the db query service here
				fmt.Println("Dequeued:")
				fmt.Println(" - Transaction Query:", transac)
				fmt.Println(" - Open Account Query:", openAcc)
			} else {
				fmt.Println("Queue is empty.")
			}
		}
	}()
}

// func newQueue() *Queue {
// 	return &Queue{
// 		transacQuery: []string{},
// 		openAccQuery: []string{},
// 	}
// }