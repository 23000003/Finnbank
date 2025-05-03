package main

import (
	"context"
	"finnbank/api-gateway/routers"
	"finnbank/api-gateway/services"
	"finnbank/common/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	t "finnbank/api-gateway/types"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// func CorsMiddleware(r *gin.Engine) {
// 	r.Use(cors.New(cors.Config{
// 		AllowOrigins:     []string{"http://localhost:5173"},
// 		AllowMethods:     []string{"*"},
// 		AllowCredentials: true,
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
// 	}))
// }

func CorsMiddleware(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
}

func main() {
	router := gin.New()
	godotenv.Load()
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

		startAndShutdownServer(logger, router, apiServices)
	}
}

func startAndShutdownServer(logger *utils.Logger, router *gin.Engine, apiServices *t.ApiGatewayServices) {
	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	// Graceful shutdown channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server: %v", err)
		}
	}()

	logger.Info("Server running on http://localhost:8080")

	<-quit
	logger.Info("Shutting down server...")

	// Creates shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown WebSocket service first
	apiServices.RealTimeService.Shutdown()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited properly")
}