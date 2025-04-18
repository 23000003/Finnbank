package resolvers

import (
	"finnbank/common/grpc/auth"
	"finnbank/graphql-api/services"
	"finnbank/graphql-api/types"
	"fmt"

	// "finnbank/graphql-api/db"
	"github.com/graphql-go/graphql"
	"github.com/jackc/pgx/v5"
)

// THEN REPLACE RESOLVER LOGIC WITH THE HELPERS
func (s *StructGraphQLResolvers) GetAccountQueryType(DB *pgx.Conn) *graphql.Object {
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
						// acc_num, ok := p.Args["account_number"].(string)
						// if ok {
						// 	res, err := accServer.GetAccountById(p.Context, &account.AccountRequest{
						// 		AccountNumber: acc_num,
						// 	})

						// 	if err != nil {
						// 		s.log.Error("gRPC server error: %v", err)
						// 		return nil, fmt.Errorf("gRPC error occured: %v", err)
						// 	}
						// 	return res.Account, err
						// }
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
						// email, ok := p.Args["email"].(string)
						// if ok {
						// 	res, err := accServer.GetAccountByEmail(p.Context, &account.EmailRequest{
						// 		Email: email,
						// 	})
						// 	if err != nil {
						// 		s.log.Error("gRPC server error: %v", err)
						// 		return nil, fmt.Errorf("gRPC error occured: %v", err)
						// 	}
						// 	return res.Account, err
						// }
						return nil, fmt.Errorf("email argument is required and must be a string")
					},
				},
			},
		})
}

func (s *StructGraphQLResolvers) GetAccountMutationType(DB *pgx.Conn, authServer auth.AuthServiceClient) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"create_account": &graphql.Field{
				Type:        accountType,
				Description: "Create a new account",
				Args: graphql.FieldConfigArgument{
					"account": &graphql.ArgumentConfig{
						Type: types.AccountInputType,
					},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					accountInput, ok := p.Args["account"].(map[string]any)
					if !ok {
						return nil, fmt.Errorf("invalid account input")
					}
					email, _ := accountInput["email"].(string)
					password, _ := accountInput["password"].(string)
					firstName, _ := accountInput["first_name"].(string)
					lastName, _ := accountInput["last_name"].(string)
					phoneNumber, _ := accountInput["phone_number"].(string)
					address, _ := accountInput["address"].(string)
					accountType, _ := accountInput["account_type"].(string)
					nationality, _ := accountInput["nationality"].(string)
					req := &types.AddAccountRequest{
						Email:       email,
						Password:    password,
						FullName:    firstName + " " + lastName,
						PhoneNumber: phoneNumber,
						Address:     address,
						AccountType: accountType,
						Nationality: nationality,
					}
					account, err := services.CreateUser(&p.Context, req, DB, authServer)
					if err != nil {
						s.log.Error("Error creating account: %v", err)
						return nil, fmt.Errorf("error creating account: %v", err)
					}
					s.log.Info("Account created successfully: %v", account.ID)

					return account, nil
				},
			},
			// can define a new Field here
		},
	})

}
