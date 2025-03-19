package resolvers

import (
	"finnbank/common/grpc/account"
	"fmt"

	"github.com/graphql-go/graphql"
)

// FUTURE: GET THIS DONE
func (s *StructGraphQLResolvers) GetAccountQueryType(accServer account.AccountServiceClient) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"account_by_id": &graphql.Field{
					Type:        accountType,
					Description: "Get account by account number",
					Args: graphql.FieldConfigArgument{
						"account_number": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						acc_num, ok := p.Args["account_number"].(string)
						if ok {
							res, err := accServer.GetAccountById(p.Context, &account.AccountRequest{
								AccountNumber: acc_num,
							})

							if err != nil {
								s.log.Error("gRPC server error: %v", err)
								return nil, fmt.Errorf("gRPC error occured: %v", err)
							}
							return res.Account, err
						}
						return nil, fmt.Errorf("account_number argument is required and must be a string")
					},
				},
				"account_by_email": &graphql.Field{
					Type:        accountType,
					Description: "Get account by email",
					Args: graphql.FieldConfigArgument{
						"email": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						email, ok := p.Args["email"].(string)
						if ok {
							res, err := accServer.GetAccountByEmail(p.Context, &account.EmailRequest{
								Email: email,
							})
							if err != nil {
								s.log.Error("gRPC server error: %v", err)
								return nil, fmt.Errorf("gRPC error occured: %v", err)
							}
							return res.Account, err
						}
						return nil, fmt.Errorf("email argument is required and must be a string")
					},
				},
				// FUTRURE: ADD MORE RESOLVERS
			},
		})
}

// FUTURE: ADD THIS SHIT
func (s *StructGraphQLResolvers) GetAccountMutationType(accServer account.AccountServiceClient) *graphql.Object {
	return nil
}
