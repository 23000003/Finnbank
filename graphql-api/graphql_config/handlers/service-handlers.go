package handlers

import (
	"finnbank/common/grpc/auth"
	"finnbank/common/grpc/products"
	"finnbank/common/utils"
	"finnbank/graphql-api/graphql_config/resolvers"
	"finnbank/graphql-api/grpc_client"
	sv "finnbank/graphql-api/services"
	"finnbank/graphql-api/types"

	"fmt"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type StructGraphQLHandler struct {
	l  *utils.Logger
	r  *resolvers.StructGraphQLResolvers
	db *types.StructServiceDatabasePools
}

func NewGraphQLServicesHandler(l *utils.Logger, r *resolvers.StructGraphQLResolvers, db *types.StructServiceDatabasePools) *StructGraphQLHandler {
	return &StructGraphQLHandler{
		l:  l,
		r:  r,
		db: db,
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

func (g *StructGraphQLHandler) AccountServicesHandler(connAddress string) *handler.Handler {

	grpcConnection := grpc_client.NewGRPCClient(connAddress)
	authServer := auth.NewAuthServiceClient(grpcConnection)

	ACCService := sv.NewAccountService(g.db.AccountDBPool, g.l, authServer)

	accountSchema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    g.r.GetAccountQueryType(ACCService),
			Mutation: g.r.GetAccountMutationType(ACCService),
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

	// Follow AccountServiceHandler | OpenedAccountServicesHandler

	// grpcConnection := grpc_client.NewGRPCClient(connAddress)
	// proto server

	return nil
}

func (g *StructGraphQLHandler) NotificationServicesHandler(connAddress string) *handler.Handler {

	// Follow AccountServiceHandler | OpenedAccountServicesHandler

	// grpcConnection := grpc_client.NewGRPCClient(connAddress)
	// proto server

	notifService := sv.NewNotificationService(g.db.NotificationDBPool, g.l)

	notifSchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    g.r.GetNotificationQueryType(notifService),
		Mutation: g.r.GetNotificationMutationType(notifService),
	})
	if err != nil {
		g.l.Fatal("Failed to configure Notification Handler Schema: %v", err)
	}

	notificationHandler := handler.New(&handler.Config{
		Schema:   &notifSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	return notificationHandler

}

func (g *StructGraphQLHandler) StatementServicesHandler(connAddress string) *handler.Handler {

	// Follow AccountServiceHandler | OpenedAccountServicesHandler

	// grpcConnection := grpc_client.NewGRPCClient(connAddress)
	// proto server

	return nil
}

func (g *StructGraphQLHandler) transactionServicesHandler() (http.Handler, error) {
	if g.db.TransactionDBPool == nil {
		return nil, fmt.Errorf("transaction DB pool is not initialized")
	}

	g.l.Info("🛠 Initializing TransactionServicesHandler…")
	txSvc := sv.NewTransactionService(g.db.TransactionDBPool, g.l)

	userQuery := g.r.GetTransactionQueryType(txSvc)
	timeQuery := g.r.GetTransactionByTimeStampByUserId(txSvc)

	mergedFields := graphql.Fields{}
	for name, field := range userQuery.Fields() {
		mergedFields[name] = field
	}
	for name, field := range timeQuery.Fields() {
		mergedFields[name] = field
	}
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: mergedFields,
	})
	mutation := g.r.GetTransactionMutationType(txSvc)

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: mutation,
	})

	// schema, err := graphql.NewSchema(graphql.SchemaConfig{
	// 	Query:    g.r.GetTransactionQueryType(txSvc),
	// 	Mutation: g.r.GetTransactionMutationType(txSvc),
	// })
	if err != nil {
		g.l.Error("❌ Failed to configure Transaction schema: %v", err)
		return nil, fmt.Errorf("failed to configure transaction schema: %w", err)
	}
	g.l.Info("✅ Transaction GraphQL schema created")

	graphiql := os.Getenv("ENABLE_GRAPHIQL") == "true"
	gqlHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: graphiql,
	})

	g.l.Info("🚀 TransactionServicesHandler ready")
	return gqlHandler, nil
}
func (g *StructGraphQLHandler) TransactionServicesHandler(connAddress string) *handler.Handler {
	h, err := g.transactionServicesHandler()
	if err != nil {
		g.l.Fatal("Could not initialize TransactionServicesHandler: %v", err)
	}
	// at this point h is the *gqlhandler.Handler you created above
	return h.(*handler.Handler)
}
func (g *StructGraphQLHandler) OpenedAccountServicesHandler(connAddress string) *handler.Handler {

	// grpcConnection := grpc_client.NewGRPCClient(connAddress)
	// proto server

	OAService := sv.NewOpenedAccountService(g.db.OpenedAccountDBPool, g.l)

	openedAccountSchema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    g.r.GetOpenedAccountQueryType(OAService),
			Mutation: g.r.GetOpenedAccountMutationType(OAService),
		},
	)
	if err != nil {
		g.l.Fatal("Failed to configure Opened Account Handler Schema: %v", err)
	}

	openedAccountHandler := handler.New(&handler.Config{
		Schema:   &openedAccountSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	return openedAccountHandler

}
