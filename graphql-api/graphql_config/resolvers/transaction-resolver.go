// func (s *StructGraphQLResolvers) GetTransactionQueryType(#) *graphql.Object {

// }

// func (s *StructGraphQLResolvers) GetTransactionMutationType(#) *graphql.Object {

// }

// Connects the Graphql schema to the businesslogic (services) and the database (db).
package resolvers

import (
	"finnbank/graphql-api/services"
	t "finnbank/graphql-api/types"
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
func (r *StructGraphQLResolvers) GetTransactionQueryType(txSvc *services.TransactionService) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query", // or "TransactionQuery"
		Fields: graphql.Fields{
			"getTransactionsByUserId": &graphql.Field{
				Type:        graphql.NewList(transactionType),
				Description: "Fetch all transactions for a specific user by their user ID.",
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "The ID of the user whose transactions are to be fetched.",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userId := p.Args["userId"].(string)
					return txSvc.GetTransactionByUserId(p.Context, userId)
				},
			},
		},
	})
}

// for creating a new transaction
// mutation name: createTransaction
// mutation accepts a single argument: transaction (object)
// this argument specifies the transaction details to be created
// when called resolve function is executed
// ensure that all requred feilds are provided and valid
// creates a transaction object from the input
// calls the CreateTransaction method of the TrsanctionService to create the transaction
// GetTransactionMutationType defines the mutation for creating a transaction.
func (r *StructGraphQLResolvers) GetTransactionMutationType(txSvc *services.TransactionService) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation", // or "TransactionMutation"
		Fields: graphql.Fields{
			"createTransaction": &graphql.Field{
				Type:        transactionType,
				Description: "Create a new transaction (ref_no auto‑generated).",
				Args: graphql.FieldConfigArgument{
					"transaction": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(transactionInputType),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					in := p.Args["transaction"].(map[string]interface{})
					txn := t.Transaction{
						SenderID:        in["sender_id"].(string),
						ReceiverID:      in["receiver_id"].(string),
						TransactionType: in["transaction_type"].(string),
						Amount:          in["amount"].(float64),
						TransactionFee:  in["transaction_fee"].(float64),
						Notes:           in["notes"].(string),
					}
					return txSvc.CreateTransaction(p.Context, txn)
				},
			},
		},
	})
}

// for cmmit: change because to define one only. No duplicate names. easier maintenance. cleaner.
//GetTransactionQueryType: Defines a query to fetch all transactions for a specific user by their userId.
//GetTransactionMutationType: Defines a mutation to create a new transaction with detailed input fields.
