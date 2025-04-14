package db

import (
	"context"
	"finnbank/common/utils"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func InitDb(ctx context.Context) (*pgx.Conn, error) {
	logger, err1 := utils.NewLogger()
	if err1 != nil {
		panic(err1)
	}

	err := godotenv.Load()
	if err != nil {
		logger.Warn("Can't find Environment Variables")
	}
	// LOCAL_DB_URL <-- LOCAL Database
	// ACC_DATABASE_URL <-- PROD Database
	dbURL := os.Getenv("LOCAL_DB_URL")
	if dbURL == "" {
		logger.Fatal("ACC_DATABASE_URL is missing")
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
