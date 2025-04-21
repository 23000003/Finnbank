package main

import (
	"finnbank/common/utils"
	"finnbank/api-gateway/routers"
	"finnbank/api-gateway/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
}

func main() {
	router := gin.New()
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("Starting the application...")

	CorsMiddleware(router)

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
