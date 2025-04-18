package resolvers

import (
	"finnbank/common/grpc/auth"
	"finnbank/graphql-api/services"
	"finnbank/graphql-api/types"
	"fmt"
	sv "finnbank/graphql-api/services"

	// "finnbank/graphql-api/db"
	"github.com/graphql-go/graphql"
)

// THEN REPLACE RESOLVER LOGIC WITH THE HELPERS
func (s *StructGraphQLResolvers) GetAccountQueryType(ACCService *sv.AccountService) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"account_by_account_num": &graphql.Field{
					Type:        accountType,
					Description: "Get account by account number",
					Args: graphql.FieldConfigArgument{
						"account_number": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (any, error) {
						req, ok := p.Args["account_number"].(string)
						if !ok {
							return nil, fmt.Errorf("account_number argument is required and must be a string")
						}

						res, err := services.FetchUserByAccountNumber(&p.Context, req, DB)
						if err != nil {
							s.log.Error("Error fetching account: %v", err)
							return nil, fmt.Errorf("error fetching account: %v", err)
						}
						s.log.Info("Account fetched successfully: %v", res.Account.ID)
						return res.Account, nil
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
					Resolve: func(p graphql.ResolveParams) (any, error) {
						req, ok := p.Args["email"].(string)
						if !ok {
							return nil, fmt.Errorf("email argument is required and must be a string")
						}

						res, err := services.FetchUserByEmail(&p.Context, req, DB)
						if err != nil {
							s.log.Error("Error fetching account: %v", err)
							return nil, fmt.Errorf("error fetching account: %v", err)
						}
						s.log.Info("Account fetched successfully: %v", res.Account.ID)
						return res.Account, nil
					},
				},
				"account_by_phone": &graphql.Field{
					Type:        accountType,
					Description: "Get account by phone number",
					Args: graphql.FieldConfigArgument{
						"phone_number": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (any, error) {
						req, ok := p.Args["phone_number"].(string)
						if !ok {
							return nil, fmt.Errorf("phone_number argument is required and must be a string")
						}
						res, err := services.FetchUserByPhone(&p.Context, req, DB)
						if err != nil {
							s.log.Error("Error fetching account: %v", err)
							return nil, fmt.Errorf("error fetching account: %v", err)
						}
						s.log.Info("Account fetched successfully: %v", res.Account.ID)
						return res.Account, nil
					},
				},
				"account_by_auth_id": &graphql.Field{
					Type:        accountType,
					Description: "Get account by phone number",
					Args: graphql.FieldConfigArgument{
						"auth_id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (any, error) {
						req, ok := p.Args["auth_id"].(string)
						if !ok {
							return nil, fmt.Errorf("auth_id argument is required and must be a string")
						}
						res, err := services.FetchUserByAuthID(&p.Context, req, DB)
						if err != nil {
							s.log.Error("Error fetching account: %v", err)
							return nil, fmt.Errorf("error fetching account: %v", err)
						}
						s.log.Info("Account fetched successfully: %v", res.Account.ID)
						return res.Account, nil
					},
				},
			},
		})
}

func (s *StructGraphQLResolvers) GetAccountMutationType(ACCService *sv.AccountService, authServer auth.AuthServiceClient) *graphql.Object {
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
			"login": &graphql.Field{
				Type:        types.AuthResponseType,
				Description: "Login to an account",
				Args: graphql.FieldConfigArgument{
					"account": &graphql.ArgumentConfig{
						Type: types.LoginInputType,
					},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					loginInput, ok := p.Args["account"].(map[string]any)
					if !ok {
						return nil, fmt.Errorf("invalid account input")
					}
					email, _ := loginInput["email"].(string)
					password, _ := loginInput["password"].(string)
					req := &types.LoginRequest{
						Email:    email,
						Password: password,
					}
					res, err := services.Login(&p.Context, req, authServer)
					if err != nil {
						s.log.Error("Error logging in: %v", err)
						return nil, fmt.Errorf("error logging in: %v", err)
					}
					s.log.Info("Login successful: %v", res.AuthID)
					return res, nil
				},
			},
		},
	})

}
