package db

import (
	"finnbank/services/common/utils"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/auth-go"
	"github.com/supabase-community/supabase-go"
)

/**
	Supabase Connection string not client for postgre db
	since sir wants us to use use pressly for migrations
	dont use local postgredb since it will conflict our ports and stuffs
**/

func SupabaseInit() (*supabase.Client, auth.Client) {
	// var local_url string = "LOCAL_DB_URL"
	// var local_key string = "LOCAL_DB_KEY"
	logger, err1 := utils.NewLogger()
	if err1 != nil {
		panic(err1)
	}
	var super_key string = "SERVICE_ROLE_KEY"
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Missing env files")
	}
	url := os.Getenv("DB_URL")
	key := os.Getenv(super_key)
	auth_url := os.Getenv("AUTH_DB_URL")
	if url == "" || key == "" || auth_url == "" {
		logger.Fatal("Supabase URL and Keys missing")
	}
	client, err := supabase.NewClient(url, key, &supabase.ClientOptions{})
	authClient := auth.New(auth_url, key)

	if err != nil {
		logger.Fatal("Failed to initialize Supabase client: %v", err)
	}

	return client, authClient
}
