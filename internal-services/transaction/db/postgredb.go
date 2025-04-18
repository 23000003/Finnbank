package db

/**
	Supabase Connection string not client for postgre db
	since sir wants us to use use pressly for migrations
	dont use local postgredb since it will conflict our ports and stuffs
**/
import (
	"context"
	"finnbank/common/utils"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func NewBankCardDB(ctx context.Context) (*pgx.Conn, error) {
	logger, err1 := utils.NewLogger()
	if err1 != nil {
		panic(err1)
	}

	err := godotenv.Load()
	if err != nil {
		logger.Warn("Can't find Environment Variables")
	}

	dbURL := os.Getenv("CONNECTION_STRING")
	if dbURL == "" {
		logger.Fatal("Connection string is missing")
	}

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		logger.Fatal("Failed to connect to PostgreSQL (BankCard): %v", err)
		return nil, err
	}

	logger.Info("Successfully connected to PostgreSQL (BankCard)")
	return conn, nil
}
