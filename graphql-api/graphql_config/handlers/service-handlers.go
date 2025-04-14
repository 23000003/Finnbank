package handlers

import (
	"context"
	"finnbank/common/grpc/products"
	"finnbank/common/utils"
	"finnbank/graphql-api/graphql_config/resolvers"
	"finnbank/graphql-api/grpc_client"

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

// HIGHLIGHT
// Instead of adding gRPC here, initialize DB connector add pass it to the Query and Mutations

func (g *StructGraphQLHandler) AccountServicesHandler() *handler.Handler {
	ctx := context.Background()
	db, _ := NewDbHandler(ctx)
	accountSchema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    g.r.GetAccountQueryType(db),
			Mutation: g.r.GetAccountMutationType(db),
		},
	)
	if err != nil {
		g.l.Fatal("Failed to configure Account Handler Schema: %v", err)
	}

	accountHandler := handler.New(&handler.Config{
		Schema:   &accountSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	return accountHandler
}

// >>>>>>>>>>>>>============= ADD HERE ========== <<<<<<<<<<<<<<

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
