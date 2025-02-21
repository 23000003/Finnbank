package main

import (
    "github.com/gin-gonic/gin"
    "finnbank/services/common/utils"
    "finnbank/services/api-gateway/routers"
)

func main() {
    router := gin.New()
    logger, err := utils.NewLogger()
    if err != nil {
        panic(err)
    }
    logger.Info("Starting the application...")

    // router group with base path for each services (to identify the service)
    serviceAPI := router.Group("/api")
    {
        // Use the group for your routes
        routers.GatewayRouter(serviceAPI)
    }
    
    logger.Info("Server running on http://localhost:8080")

    if err := router.Run("localhost:8080"); err != nil {
        logger.Error("Failed to start server: %v", err)
    }
}