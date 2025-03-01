package db

import (
	"context"
	"finnbank/services/graphql/utils"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

/**
	Supabase Connection string not client for postgre db
	since sir wants us to use use pressly for migrations
	dont use local postgredb since it will conflict our ports and stuffs
**/

// func SupabaseInit() (*supabase.Client, auth.Client) {
// 	// var local_url string = "LOCAL_DB_URL"
// 	// var local_key string = "LOCAL_DB_KEY"
// 	logger, err1 := utils.NewLogger()
// 	if err1 != nil {
// 		panic(err1)
// 	}
// 	var super_key string = "SERVICE_ROLE_KEY"
// 	err := godotenv.Load()
// 	if err != nil {
// 		logger.Fatal("Missing env files")
// 	}
// 	url := os.Getenv("DB_URL")
// 	key := os.Getenv(super_key)
// 	auth_url := os.Getenv("AUTH_DB_URL")
// 	if url == "" || key == "" || auth_url == "" {
// 		logger.Fatal("Supabase URL and Keys missing")
// 	}
// 	client, err := supabase.NewClient(url, key, &supabase.ClientOptions{})
// 	authClient := auth.New(auth_url, key)

// 	if err != nil {
// 		logger.Fatal("Failed to initialize Supabase client: %v", err)
// 	}

// 	return client, authClient
// }

func InitDb(ctx context.Context) (*pgx.Conn, error) {
	logger, err1 := utils.NewLogger()
	if err1 != nil {
		panic(err1)
	}

	err := godotenv.Load()
	if err != nil {
		logger.Warn("Can't find Environment Variables")
	}

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
