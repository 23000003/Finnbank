package main

import (
	"finnbank/common/utils"
	"github.com/rs/cors"
	"net/http"
	"finnbank/graphql-api/db"
	"context"
	"finnbank/graphql-api/graphql_config"
	"finnbank/graphql-api/graphql_config/handlers"
	"finnbank/graphql-api/graphql_config/resolvers"
	"os"
	"os/signal"
	"syscall"
	"time"
	"finnbank/graphql-api/types"
)

func CorsMiddleware() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	})
}

func main() {
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("Starting the application...")
	
	dbServicesPool := db.InitializeServiceDatabases(logger)
	defer db.CleanupDatabase(dbServicesPool, logger) 

	server := initializeGraphQL(logger, dbServicesPool)
	
	startAndShutdownServer(server, logger)
}

func initializeGraphQL(logger *utils.Logger, dbPool *types.StructServiceDatabasePools) *http.Server {
	resolvers := resolvers.NewGraphQLResolvers(logger)
	handlers := handlers.NewGraphQLServicesHandler(logger, resolvers, dbPool)
	graphql := graphql_config.NewGraphQL(logger, handlers)
	graphql.ConfigureGraphQLHandlers()

	http.HandleFunc("/graphql/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Graphql API is OK."))
	})

	return &http.Server{
		Addr:    ":8083",
		Handler: CorsMiddleware().Handler(http.DefaultServeMux),
	}
}

func startAndShutdownServer(server *http.Server, logger *utils.Logger) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Server running on http://localhost:8083")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: %v", err)
		}
	}()

	// Triggers Below When CTRL + C or Shutdown 
	<-done
	logger.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed: %v", err)
	}
	logger.Info("Server exited properly")
}

