package handlers

import (
	"finnbank/common/grpc/products"
	"finnbank/common/utils"
	"finnbank/internal-services/graphql-api/graphql_config/resolvers"
	"finnbank/internal-services/graphql-api/grpc_client"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type StructGraphQLHandler struct {
	l *utils.Logger
	r *resolvers.StructGraphQLResolvers
}

func NewGraphQLServicesHandler(l *utils.Logger, r *resolvers.StructGraphQLResolvers) *StructGraphQLHandler {
	return &StructGraphQLHandler{
		l: l,
		r: r,
	}
}

func (g *StructGraphQLHandler) ProductServicesHandler(connAddress string) *handler.Handler {

	grpcConnection := grpc_client.NewGRPCClient(connAddress)
	prodsServer := products.NewProductServiceClient(grpcConnection)

	productSchema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    g.r.GetProductQueryType(prodsServer),
			Mutation: g.r.GetProductMutationType(prodsServer),
		},
	)

	if err != nil {
		g.l.Fatal("Failed to configure Product Handler Schema %v", err)
	}

	productHandler := handler.New(&handler.Config{
		Schema:   &productSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	return productHandler
}

// ================== Below have no Protobuf yet (will wait for u if ur assigned to it)==================

func (g *StructGraphQLHandler) AccountServicesHandler(connAddress string) *handler.Handler {

	// grpcConnection := grpc_client.NewGRPCClient(connAddress)
	// proto server

	return nil
}

func (g *StructGraphQLHandler) BankCardServicesHandler(connAddress string) *handler.Handler {

	// grpcConnection := grpc_client.NewGRPCClient(connAddress)
	// proto server

	return nil
}

func (g *StructGraphQLHandler) NotificationServicesHandler(connAddress string) *handler.Handler {

	// grpcConnection := grpc_client.NewGRPCClient(connAddress)
	// proto server

	return nil
}

func (g *StructGraphQLHandler) StatementServicesHandler(connAddress string) *handler.Handler {

	// grpcConnection := grpc_client.NewGRPCClient(connAddress)
	// proto server

	return nil
}

func (g *StructGraphQLHandler) TransactionServicesHandler(connAddress string) *handler.Handler {

	// grpcConnection := grpc_client.NewGRPCClient(connAddress)
	// proto server

	return nil
}
