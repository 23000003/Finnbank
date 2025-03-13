package main

// setup http (for response test only)

import (
	"finnbank/common/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRoot(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Welcome to the Home Page!")
}

func RunHttpServer(logger *utils.Logger) {

	router := gin.New()

	api := router.Group("/api/statement")
	{
		// Use the group for your routes
		api.GET("/test", GetRoot)
	}

	logger.Info("Http server running on http://localhost:8084")

	if err := router.Run("localhost:8084"); err != nil {
		logger.Error("Failed to start http server: %v", err)
	}
}
