package main

import (
	"finnbank/services/api-gateway/routers"
	"finnbank/services/api-gateway/services"
	"finnbank/services/common/utils"
	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func CorsMiddleware(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"*"},
        AllowCredentials: true,
    }))
}

func main() {
    router := gin.New()
    logger, err := utils.NewLogger()
    if err != nil {
        panic(err)
    }
    logger.Info("Starting the application...")

    
    serviceAPI := router.Group("/api")
    {
        serviceAPI.GET("/health", HealthCheck)
        serviceAPI.GET("/graphql-health", GraphQLAPIHealthCheck)
        
        apiServices := services.NewApiGatewayServices(logger)
        gatewayRouters := routers.NewGatewayRouter(logger, serviceAPI, apiServices)
        gatewayRouters.ConfigureGatewayRouter()
    }
    
    logger.Info("Server running on http://localhost:8080")

    if err := router.Run("localhost:8080"); err != nil {
        logger.Error("Failed to start server: %v", err)
    }
}