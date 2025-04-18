package db

import (
	"context"
	"finnbank/common/utils"
	"os"

	"github.com/jackc/pgx/v5"
)

func NewNotificationDB(ctx context.Context) (*pgx.Conn, error) {
	logger, err1 := utils.NewLogger()
	if err1 != nil {
		panic(err1)
	}

	dbURL := os.Getenv("CONNECTION_STRING")
	if dbURL == "" {
		logger.Fatal("NOTIF_DATABASE_URL is missing")
	}
	logger.Info(dbURL)

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		logger.Fatal("Failed to connect to PostgreSQL: %v", err)
		return nil, err
	}
	logger.Info("Successfully connected to PostgreSQL")
	return conn, nil
}
