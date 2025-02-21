package main

import (
	"finnbank/services/account/db"
	"finnbank/services/account/handlers"
	"finnbank/services/account/service"
	"finnbank/services/common/utils"
	"github.com/gin-gonic/gin"
)

/*Transfer http configuration to http.go*/
func main() {
	router := gin.New()
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}

	client, auth := db.SupabaseInit()
	accountService := &service.AccountService{Client: client, Auth: auth, Logger: logger}

	logger.Info("Starting the application...")

	serviceAPI := router.Group("/api/account") // base path
	handlers.AccountRouter(serviceAPI, accountService)

	if err := router.Run("localhost:8082"); err != nil {
		logger.Fatal("Failed to start server: %v", err)
	}
}
