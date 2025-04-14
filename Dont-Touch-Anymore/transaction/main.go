package main

import (
	"context"
	"finnbank/common/utils"
	"finnbank/internal-services/transaction/db"
)

func main() {
	ctx := context.Background()
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("Starting the application...")
	_, err = db.InitDb(ctx)
	if err != nil {
		logger.Fatal("Failed to connect to DB: %f", err)
		return
	}
}
