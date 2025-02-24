package main

import (
	"finnbank/services/common/utils"
	"finnbank/services/graphql/graphql_config"
	"finnbank/services/graphql/graphql_config/handlers"
	"finnbank/services/graphql/graphql_config/resolvers"
	"net/http"
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

func InitializeGraphQLServer(logger *utils.Logger) {
	graphql_resolvers := resolvers.NewGraphQLResolvers(logger)
	graphql_handlers := handlers.NewGraphQLServicesHandler(logger, graphql_resolvers)
	graphql := graphql_config.NewGraphQL(graphql_handlers)
    graphql.ConfigureGraphQLHandlers(logger)
}

func main() {

	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("Starting the application...")
	
	InitializeGraphQLServer(logger)

	logger.Info("Server running on http://localhost:8083")

    corsHandler := CorsMiddleware().Handler(http.DefaultServeMux)
	err = http.ListenAndServe(":8083", corsHandler)

	if err != nil {
		logger.Error("Failed to start server: %v", err)
	}
}
