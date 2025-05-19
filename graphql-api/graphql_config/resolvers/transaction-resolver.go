// func (s *StructGraphQLResolvers) GetTransactionQueryType(#) *graphql.Object {

// }

// func (s *StructGraphQLResolvers) GetTransactionMutationType(#) *graphql.Object {

// }

// Connects the Graphql schema to the businesslogic (services) and the database (db).
package resolvers

import (
	e "finnbank/graphql-api/graphql_config/entities"
	"finnbank/graphql-api/services"
	t "finnbank/graphql-api/types"
	"time"

	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
)

type TransactionResolver struct {
	Service *services.TransactionService
}

//query name: getTransactionsByUserId
// query acceps a single argument: userId (string)
// this argument specifies the user whose transactions should be fetched.
// when called resolve function is executed
// extracts userid from the arguments and calls the service method to fetch transactions
// returns the transactions or an error if any occurs

// GetTransactionQueryType defines the query for fetching transactions.
// just use the ById below
// func (r *StructGraphQLResolvers) GetTransactionQueryType(txSvc *services.TransactionService) *graphql.Object {
// 	return graphql.NewObject(graphql.ObjectConfig{
// 		Name: "Query", // or "TransactionQuery"
// 		Fields: graphql.Fields{
// 			"getTransactionsByOpenAccountId": &graphql.Field{
// 				Type:        graphql.NewList(transactionType),
// 				Description: "Fetch all transactions for a specific user by their opened accounts ID.",
// 				Args: graphql.FieldConfigArgument{
// 					"accountId": &graphql.ArgumentConfig{
// 						Type:        graphql.NewNonNull(graphql.Int),
// 						Description: "The ID of the user whose transactions are to be fetched.",
// 					},
// 				},
// 				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 					userId := p.Args["accountId"].(int)
// 					return txSvc.GetTransactionByUserId(p.Context, userId)
// 				},
// 			},
// 		},
// 	})
// }

// for creating a new transaction
// mutation name: createTransaction
// mutation accepts a single argument: transaction (object)
// this argument specifies the transaction details to be created
// when called resolve function is executed
// ensure that all requred feilds are provided and valid
// creates a transaction object from the input
// calls the CreateTransaction method of the TrsanctionService to create the transaction
// GetTransactionMutationType defines the mutation for creating a transaction.
// ===========
// func (r *StructGraphQLResolvers) GetTransactionMutationType(txSvc *services.TransactionService) *graphql.Object {
// 	return graphql.NewObject(graphql.ObjectConfig{
// 		Name: "Mutation", // or "TransactionMutation"
// 		Fields: graphql.Fields{
// 			"createTransaction": &graphql.Field{
// 				Type:        transactionType,
// 				Description: "Create a new transaction (ref_no auto‑generated).",
// 				Args: graphql.FieldConfigArgument{
// 					"transaction": &graphql.ArgumentConfig{
// 						Type: graphql.NewNonNull(transactionInputType),
// 					},
// 				},
// 				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 					in := p.Args["transaction"].(map[string]interface{})
// 					txn := t.Transaction{
// 						SenderID:        in["sender_id"].(int),
// 						ReceiverID:      in["receiver_id"].(int),
// 						TransactionType: in["transaction_type"].(string),
// 						Amount:          in["amount"].(float64),
// 						TransactionFee:  in["transaction_fee"].(float64),
// 						Notes:           in["notes"].(string),
// 					}
// 					return txSvc.CreateTransaction(p.Context, txn)
// 				},
// 			},
// 		},
// 	})
// }



// for cmmit: change because to define one only. No duplicate names. easier maintenance. cleaner.
//GetTransactionQueryType: Defines a query to fetch all transactions for a specific user by their userId.
//GetTransactionMutationType: Defines a mutation to create a new transaction with detailed input fields.

// GetTransactionByTimeStamByUserId

func (r *StructGraphQLResolvers) TransactionQueryFields(
	txSvc *services.TransactionService,
) graphql.Fields {
	return graphql.Fields{
		"getTransactionsByUserId": &graphql.Field{
			Type:        graphql.NewList(e.GetTransactionEntityType()),
			Description: "Fetch all transactions for a specific user by userId.",
			Args: graphql.FieldConfigArgument{
				"creditId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"debitId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"savingsId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"limit": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				creditId := p.Args["creditId"].(int)
				debitId := p.Args["debitId"].(int)
				savingsId := p.Args["savingsId"].(int)
				limit := p.Args["limit"].(int)
				return txSvc.GetTransactionByUserId(p.Context, creditId, debitId, savingsId, limit)
			},
		},
		"getIsAccountAtLimit": &graphql.Field{
			Type:        graphql.NewList(graphql.Boolean),
			Description: "Identifies if opened account is at limit.",
			Args: graphql.FieldConfigArgument{
				"creditId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"debitId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"savingsId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"account_type": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				creditId := p.Args["creditId"].(int)
				debitId := p.Args["debitId"].(int)
				savingsId := p.Args["savingsId"].(int)
				accountType := p.Args["account_type"].(string)
				return txSvc.GetIsAccountAtLimitByAccountId(p.Context, accountType, creditId, debitId, savingsId)
			},
		},
	}
}

// returns the fields for “getTransactionsByTimeStampByUserId”
func (r *StructGraphQLResolvers) TransactionTimeQueryFields(
	txSvc *services.TransactionService,
) graphql.Fields {
	return graphql.Fields{
		// http://localhost:8083/graphql/transaction?query=
		// {getTransactionsByTimeStampByUserId(creditId:1, debitId:2, savingsId:3, 
		// startTime:"2025-04-19T00:00:00Z", endTime:"2025-04-29T00:00:00Z"){ref_no}}
		"getTransactionsByTimeStampByUserId": &graphql.Field{
			Type:        graphql.NewList(e.GetTransactionEntityType()),
			Description: "Fetch transactions for a user between two timestamps.",
			Args: graphql.FieldConfigArgument{
				"creditId":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"debitId":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"savingsId":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"startTime": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.DateTime)},
				"endTime":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.DateTime)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				creditId := p.Args["creditId"].(int)
				debitId := p.Args["debitId"].(int)
				savingsId := p.Args["savingsId"].(int)
				start := p.Args["startTime"].(time.Time)
				end := p.Args["endTime"].(time.Time)
				return txSvc.GetTransactionByTimestampByUserId(p.Context, creditId, debitId, savingsId, start, end)
			},
		},
	}
}

// returns the fields for your createTransaction mutation
func (r *StructGraphQLResolvers) TransactionMutationFields(
	txSvc *services.TransactionService,
	transacConn *websocket.Conn,
) graphql.Fields {
	return graphql.Fields{
		"createTransaction": &graphql.Field{
			Type:        e.GetTransactionEntityType(),
			Description: "Create a new transaction (ref_no auto-generated).",
			Args: graphql.FieldConfigArgument{
				"transaction": &graphql.ArgumentConfig{Type: graphql.NewNonNull(e.GetTransactionInputType())},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				in := p.Args["transaction"].(map[string]interface{})
				txn := t.Transaction{
					SenderID:        in["sender_id"].(int),
					ReceiverID:      in["receiver_id"].(int),
					TransactionType: in["transaction_type"].(string),
					Amount:          in["amount"].(float64),
					TransactionFee:  in["transaction_fee"].(float64),
					Notes:           in["notes"].(string),
				}
				return txSvc.CreateTransaction(p.Context, txn, transacConn)
			},
		},
	}
}
