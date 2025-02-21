package graphql_config

/**
 Handles all the graphql requests and Configure Routes
**/

import (
	"finnbank/services/graphql/graphql_config/resolvers"
	"net/http"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    resolvers.GetQueryType(),
		Mutation: resolvers.GetMutationType(),
	},
)

func GraphQLHandlers() {

	httpHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql/product", httpHandler)
}
