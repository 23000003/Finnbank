package helpers

// prevent import cycle
import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"fmt"
)

func SuccessTransaction(ctx context.Context, db *pgxpool.Pool, transacId int) error {
	conn, err := db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	query := `UPDATE transaction SET transaction_status = $1 
		 WHERE transaction_id = $2 
		 RETURNING transaction_id` 

	var transac_id int
	err = conn.QueryRow(ctx, query, "Completed" ,transacId).Scan(&transac_id)

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

	query := `UPDATE transaction SET transaction_status = $1 
		 WHERE transaction_id = $2 
		 RETURNING transaction_id` 

	var transac_id int
	err = conn.QueryRow(ctx, query, "Failed", transacId).Scan(&transac_id)

	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}

func DeductOpenedAccountBalance(ctx context.Context, db *pgxpool.Pool, openAccId int, amount float64) error {
	conn, err := db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	query := `UPDATE openedaccount SET balance = balance - $1 
		 WHERE openedaccount_id = $2 
		 RETURNING openedaccount_id` 

	var open_acc_id int
	err = conn.QueryRow(ctx, query, amount ,openAccId).Scan(&open_acc_id)

	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}