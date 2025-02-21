package main

/**
	TEST SERVICE
**/

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"finnbank/services/product/utils"
)

func GetRoot(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Welcome to the Home Page!");
}

func RunHttpServer (logger *utils.Logger) {

	router := gin.New();

	logger.Info("Starting the application...")

	api := router.Group("/api/product")
	{
		// Use the group for your routes
		api.GET("/testprod", GetRoot)
	}

	logger.Info("Server running on http://localhost:9080")

	if err := router.Run("localhost:9080"); err != nil {
		logger.Error("Failed to start server: %v", err)
	}
}