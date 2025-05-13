package helpers

// prevent import cycle
import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SuccessTransaction(ctx context.Context, oaDB *pgxpool.Pool, transacDB *pgxpool.Pool, accDB *pgxpool.Pool, transacId int, receiverId int) error {
	oaConn, err := oaDB.Acquire(ctx)
	accConn, err1 := accDB.Acquire(ctx)
	transacConn, err2 := transacDB.Acquire(ctx)

	if err != nil || err1 != nil || err2 != nil {
		return fmt.Errorf("failed to acquire connection: %w %w %w", err, err1, err2)
	}
	defer oaConn.Release()

	validateReceiverStatus := `
		SELECT openedaccount_status, account_id
		FROM openedaccount
		WHERE openedaccount_id = $1 AND openedaccount_status = $2
	`

	var openAccStatus string
	var accountId string
	err = oaConn.QueryRow(ctx, validateReceiverStatus, receiverId, "Active").Scan(&openAccStatus, &accountId)
	if err != nil {
		return fmt.Errorf("validate query failed: %w", err)
	}
	if openAccStatus != "Active" {
		return fmt.Errorf("account is not active")
	}

	validateReceiverAcc := `
		SELECT account_status
		FROM account
		WHERE id = $1 AND account_status = $2
	`

	var AccountStatus string
	err = accConn.QueryRow(ctx, validateReceiverAcc, accountId, "Active").Scan(&AccountStatus)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	if AccountStatus != "Active" {
		return fmt.Errorf("account is not active")
	}

	query := `UPDATE transactions SET transaction_status = $1 
		 WHERE transaction_id = $2 
		 RETURNING transaction_id`

	var transac_id int
	err = transacConn.QueryRow(ctx, query, "Completed", transacId).Scan(&transac_id)

	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}

func FailedTransaction(ctx context.Context, db *pgxpool.Pool, transacId int) error {
	conn, err := db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	query := `UPDATE transactions SET transaction_status = $1 
		 WHERE transaction_id = $2 
		 RETURNING transaction_id`

	var transac_id int
	err = conn.QueryRow(ctx, query, "Failed", transacId).Scan(&transac_id)

	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}

func DeductOpenedAccountBalance(ctx context.Context, db *pgxpool.Pool, senderId int, receiverId int, amount float64) error {
	conn, err := db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	query := `UPDATE openedaccount SET balance = balance - $1 
		 WHERE openedaccount_id = $2 
		 RETURNING openedaccount_id`

	var open_acc_id int
	err = conn.QueryRow(ctx, query, amount, senderId).Scan(&open_acc_id)

	if err != nil {
		return fmt.Errorf("deduct update failed: %w", err)
	}

	query = `UPDATE openedaccount SET balance = balance + $1
		 WHERE openedaccount_id = $2
		 RETURNING openedaccount_id`

	err = conn.QueryRow(ctx, query, amount, receiverId).Scan(&open_acc_id)
	if err != nil {
		return fmt.Errorf("add update failed: %w", err)
	}

	return nil
}

func SendTransactionUpdate(ctx context.Context, conn *websocket.Conn, transacId int, status string) error {
	if conn == nil {
		return fmt.Errorf("websocket connection is nil")
	}

	message := map[string]interface{}{
		"transaction_id": transacId,
		"transaction_status": status,
	}

	err := conn.WriteJSON(message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}