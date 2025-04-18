package main

import (
	"finnbank/common/utils"
	"finnbank/graphql-api/graphql_config"
	"finnbank/graphql-api/graphql_config/handlers"
	"finnbank/graphql-api/graphql_config/resolvers"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func CorsMiddleware() *cors.Cors {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	})

	return c
}

func InitializeGraphQL(logger *utils.Logger) {
	graphql_resolvers := resolvers.NewGraphQLResolvers(logger)
	graphql_handlers := handlers.NewGraphQLServicesHandler(logger, graphql_resolvers)
	graphql := graphql_config.NewGraphQL(logger, graphql_handlers)
	graphql.ConfigureGraphQLHandlers()
}

func ApiHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Graphql API is OK."))
}

func main() {

	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("Starting the application...")

	err = godotenv.Load()
	if err != nil {
		logger.Warn("Can't find Environment Variables")
	}

	InitializeGraphQL(logger)

	logger.Info("Server running on http://localhost:8083")

	corsHandler := CorsMiddleware().Handler(http.DefaultServeMux)

	http.HandleFunc("/graphql/health", ApiHealthCheck)
	err = http.ListenAndServe(":8083", corsHandler)

	if err != nil {
		logger.Error("Failed to start server: %v", err)
	}
}
