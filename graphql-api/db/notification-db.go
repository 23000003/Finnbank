package db

import (
	"context"
	"finnbank/common/utils"
	"os"
	"time"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"fmt"
)

func newNotificationDB(ctx context.Context) (*pgxpool.Pool, error) {
	logger, err := utils.NewLogger()
	if err != nil {
		return nil, err
	}

	err = godotenv.Load()
	if err != nil {
		logger.Warn("Can't find Environment Variables")
	}

	dbURL := os.Getenv("NOTIFICATION_DB_URL")
	if dbURL == "" {
		logger.Fatal("NOTIFICATION_DB_URL is missing")
		return nil, fmt.Errorf("NOTIFICATION_DB_URL is missing")
	}

	// Configure with forced simple protocol
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	// Critical settings to prevent prepared statement conflicts
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, "DISCARD ALL")
		return err
	}

	// Pool settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	logger.Info("Notification Database Pool Initialized")
	return pool, nil
}