package main

import (
	"finnbank/services/graphql/graphql_config"
	"finnbank/services/common/utils"
	"net/http"
)

func main() {
	// router := gin.New()
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}

	graphql_config.GraphQLHandlers()

	// client, auth := db.SupabaseInit()

	logger.Info("Starting the application...")

	logger.Info("Server running on http://localhost:8083")
	err = http.ListenAndServe(":8083", nil)
	if err != nil {
		logger.Error("Failed to start server: %v", err)
	}
}
