// func (s *StructGraphQLResolvers) GetTransactionQueryType(#) *graphql.Object {

// }

// func (s *StructGraphQLResolvers) GetTransactionMutationType(#) *graphql.Object {

// }

// Connects the Graphql schema to the businesslogic (services) and the database (db).
package resolvers

import (
	"finnbank/graphql-api/services"
	"finnbank/graphql-api/types"
	"fmt"

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
func (r *StructGraphQLResolvers) GetTransactionQueryType(transactionService *services.TransactionService) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "TransactionQuery",
		Fields: graphql.Fields{
			"getTransactionsByUserId": &graphql.Field{
				Type:        graphql.NewList(types.GetTransactionEntityType()),
				Description: "Fetch all transactions for a specific user by their user ID.",
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "The ID of the user whose transactions are to be fetched.",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userId, ok := p.Args["userId"].(string)
					if !ok || userId == "" {
						return nil, fmt.Errorf("userId argument is required and must be a non-empty string")
					}
					return transactionService.GetTransactionByUserId(p.Context, userId)
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
func (r *StructGraphQLResolvers) GetTransactionMutationType(transactionService *services.TransactionService) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "TransactionMutation",
		Fields: graphql.Fields{
			"createTransaction": &graphql.Field{
				Type:        types.GetTransactionEntityType(),
				Description: "Create a new transaction.",
				Args: graphql.FieldConfigArgument{
					"transaction": &graphql.ArgumentConfig{Type: graphql.NewInputObject(graphql.InputObjectConfig{
						Name: "TransactionInput",
						Fields: graphql.InputObjectConfigFieldMap{
							"ref_no": &graphql.InputObjectFieldConfig{
								Type:        graphql.String,
								Description: "The reference number of the transaction.",
							},
							"sender_id": &graphql.InputObjectFieldConfig{
								Type:        graphql.String,
								Description: "The ID of the sender.",
							},
							"receiver_id": &graphql.InputObjectFieldConfig{
								Type:        graphql.String,
								Description: "The ID of the receiver.",
							},
							"transaction_type": &graphql.InputObjectFieldConfig{
								Type:        types.TransactionTypeEnum,
								Description: "The type of the transaction (e.g., Transfer, Payment, Deposit, etc.).",
							},
							"amount": &graphql.InputObjectFieldConfig{
								Type:        graphql.Float,
								Description: "The amount of the transaction.",
							},
							"transaction_fee": &graphql.InputObjectFieldConfig{
								Type:        graphql.Float,
								Description: "The transaction fee.",
							},
							"notes": &graphql.InputObjectFieldConfig{
								Type:        graphql.String,
								Description: "Additional notes for the transaction.",
							},
						},
					})},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					input, ok := p.Args["transaction"].(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("invalid input: transaction argument is required and must be an object")
					}
					if input["sender_id"] == "" || input["receiver_id"] == "" {
						return nil, fmt.Errorf("sender_id and receiver_id are required")
					}

					if input["amount"].(float64) <= 0 {
						return nil, fmt.Errorf("amount must be greater than zero")
					}

					transaction := types.Transaction{
						RefNo:           input["ref_no"].(string),
						SenderID:        input["sender_id"].(string),
						ReceiverID:      input["receiver_id"].(string),
						TransactionType: input["transaction_type"].(string),
						Amount:          input["amount"].(float64),
						TransactionFee:  input["transaction_fee"].(float64),
						Notes:           input["notes"].(string),
					}

					return transactionService.CreateTransaction(p.Context, transaction)
				},
			},
		},
	})
}

//GetTransactionQueryType: Defines a query to fetch all transactions for a specific user by their userId.
//GetTransactionMutationType: Defines a mutation to create a new transaction with detailed input fields.
