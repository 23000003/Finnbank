package handlers

import (
	"finnbank/common/grpc/auth"
	"finnbank/common/grpc/products"
	"finnbank/common/grpc/statement"
	"finnbank/common/utils"
	"finnbank/graphql-api/graphql_config/resolvers"
	"finnbank/graphql-api/grpc_client"
	sv "finnbank/graphql-api/services"
	"finnbank/graphql-api/types"
	q "finnbank/graphql-api/queue"
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

	BCService := sv.NewBankcardService(g.db.BankCardDBPool, g.l)

	bankcardSchema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    g.r.GetBankCardQueryType(BCService),
			Mutation: g.r.GetBankCardMutationType(BCService),
		},
	)
	if err != nil {
		g.l.Fatal("Failed to configure Bank Card Handler Schema: %v", err)
	}

	bankcardHandler := handler.New(&handler.Config{
		Schema:   &bankcardSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	return bankcardHandler
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

	grpcConnection := grpc_client.NewGRPCClient(connAddress)
	statementServer := statement.NewStatementServiceClient(grpcConnection)

	STservice := sv.NewStatementService(g.db.AccountDBPool, g.l, statementServer)

	accountSchema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    g.r.GetStatementQueryType(STservice),
		},
	)
	if err != nil {
		g.l.Fatal("Failed to configure Statement Handler Schema: %v", err)
	}

	statementHandler := handler.New(&handler.Config{
		Schema:   &accountSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	return statementHandler
}

func mergeFields(fieldMaps ...graphql.Fields) graphql.Fields {
	out := graphql.Fields{}
	for _, m := range fieldMaps {
		for k, f := range m {
			out[k] = f
		}
	}
	return out
}

func (g *StructGraphQLHandler) transactionServicesHandler(queue *q.Queue) (http.Handler, error) {
	if g.db.TransactionDBPool == nil {
		return nil, fmt.Errorf("transaction DB pool is not initialized")
	}
	txSvc := sv.NewTransactionService(g.db.TransactionDBPool, g.l, queue)

	// pull the three field-maps
	q1 := g.r.TransactionQueryFields(txSvc)
	q2 := g.r.TransactionTimeQueryFields(txSvc)
	m := g.r.TransactionMutationFields(txSvc)

	// build root Query and Mutation
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: mergeFields(q1, q2),
	})
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: m,
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to configure schema: %w", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: os.Getenv("ENABLE_GRAPHIQL") == "true",
	})
	return h, nil
}

func (g *StructGraphQLHandler) TransactionServicesHandler(connAddress string, queue *q.Queue) *handler.Handler {
	h, err := g.transactionServicesHandler(queue)
	if err != nil {
		g.l.Fatal("Could not initialize TransactionServicesHandler: %v", err)
	}
	return h.(*handler.Handler)
}


func (g *StructGraphQLHandler) OpenedAccountServicesHandler(connAddress string) *handler.Handler {

	// grpcConnection := grpc_client.NewGRPCClient(connAddress)
	// proto server

	OAService := sv.NewOpenedAccountService(g.db.OpenedAccountDBPool, g.l)
	BCService := sv.NewBankcardService(g.db.BankCardDBPool, g.l)

	openedAccountSchema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    g.r.GetOpenedAccountQueryType(OAService, BCService),
			Mutation: g.r.GetOpenedAccountMutationType(OAService, BCService),
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
