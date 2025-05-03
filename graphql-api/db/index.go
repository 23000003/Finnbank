package db

import (
	"context"
	"finnbank/common/utils"
	"finnbank/graphql-api/types"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitializeServiceDatabases(logger *utils.Logger) *types.StructServiceDatabasePools {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	BCDBPool, err := newBankcardDB(ctx)
	ACCDBPool, err1 := newAccountDB(ctx)
	OADBPool, err2 := newOpenedAccountsDB(ctx)
	NOTIFDBPool, err3 := newNotificationDB(ctx)
	TRANSACDBPool, err4 := newTransactionDB(ctx)

	if err != nil || err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		logger.Error("Failed to initialize database pools: %v, %v, %v, %v, %v", err, err1, err2, err3, err4)
		return nil
	}

	return &types.StructServiceDatabasePools{
		OpenedAccountDBPool: OADBPool,
		AccountDBPool:       ACCDBPool,
		TransactionDBPool:   TRANSACDBPool,
		BankCardDBPool:      BCDBPool,
		NotificationDBPool:  NOTIFDBPool,
	}
}

func CleanupDatabase(dbPools *types.StructServiceDatabasePools, logger *utils.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Cleanup all pools in parallel (error channel)
	errChan := make(chan error, 5)

	go cleanupPool(ctx, dbPools.OpenedAccountDBPool, "OpenedAccount", errChan, logger)
	time.Sleep(3 * time.Second)
	go cleanupPool(ctx, dbPools.AccountDBPool, "Account", errChan, logger)
	time.Sleep(3 * time.Second)
	go cleanupPool(ctx, dbPools.TransactionDBPool, "Transaction", errChan, logger)
	time.Sleep(3 * time.Second)
	go cleanupPool(ctx, dbPools.BankCardDBPool, "BankCard", errChan, logger)
	time.Sleep(3 * time.Second)
	go cleanupPool(ctx, dbPools.NotificationDBPool, "Notification", errChan, logger)
	time.Sleep(3 * time.Second)

	// Wait for all cleanups to complete
	for range 5 {
		if err := <-errChan; err != nil { // WTF is <-
			logger.Error("Cleanup error: %v", err)
		}
	}

	logger.Info("All database pools closed successfully")
}

// Clean up individual pool
func cleanupPool(ctx context.Context, pool *pgxpool.Pool, poolName string, errChan chan<- error, logger *utils.Logger) {
	if pool == nil {
		errChan <- nil
		return
	}

	// Clear prepared statements
	if _, err := pool.Exec(ctx, "DISCARD ALL"); err != nil {
		errChan <- err // Bruh
		return
	}

	// Close the pool
	pool.Close()
	errChan <- nil
	logger.Info("%s pool closed", poolName)
}
