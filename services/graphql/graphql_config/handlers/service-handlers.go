package handlers

import (
	"finnbank/services/common/grpc/account"
	"finnbank/services/common/grpc/products"
	"finnbank/services/common/utils"
	"finnbank/services/graphql/graphql_config/resolvers"
	"finnbank/services/graphql/grpc_client"

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
// UPDATE: ADDED MINE, WILL HAVE THE REST LEFT BLANK FOR THEM - SEP

func (g *StructGraphQLHandler) AccountServicesHandler(connAddress string) *handler.Handler {

	grpcConnection := grpc_client.NewGRPCClient(connAddress)
	accServer := account.NewAccountServiceClient(grpcConnection)

	accountSchema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    g.r.GetAccountQueryType(accServer),
			Mutation: g.r.GetAccountMutationType(accServer),
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
