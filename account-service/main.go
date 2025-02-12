package main

import (
	"account-service/routers"
	"account-service/service"
	"account-service/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	client := service.SupabaseInit()
	accountService := &service.AccountService{Client: client}

	router := gin.New()
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("Starting the application...")
	routers.InitializeSwagger(router)

	serviceAPI := router.Group("/api/account-service") // base path
	routers.AccountRouter(serviceAPI, accountService)

	if err := router.Run("localhost:8081"); err != nil {
		logger.Fatal("Failed to start server: %v", err)
	}
}
