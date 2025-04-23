package resolvers

import (
	sv "finnbank/graphql-api/services"
	"finnbank/graphql-api/types"
	"fmt"
	"time"

	// "finnbank/graphql-api/db"
	"github.com/graphql-go/graphql"
)

// THEN REPLACE RESOLVER LOGIC WITH THE HELPERS
func (s *StructGraphQLResolvers) GetAccountQueryType(ACCService *sv.AccountService) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				// http://localhost:8083/graphql/account?query={account_by_account_num(account_number:1){<entities>}}
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

						res, err := ACCService.FetchUserByAccountNumber(&p.Context, req)
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

						res, err := ACCService.FetchUserByEmail(&p.Context, req)
						if err != nil {
							s.log.Error("Error fetching account: %v", err)
							return nil, fmt.Errorf("error fetching account: %v", err)
						}
						s.log.Info("Account fetched successfully: %v", res.Account.ID)
						return res.Account, nil
					},
				},
				"account_by_id": &graphql.Field{
					Type:        accountType,
					Description: "Get account by id",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (any, error) {
						req, ok := p.Args["id"].(string)
						if !ok {
							return nil, fmt.Errorf("id argument is required and must be a string")
						}
						res, err := ACCService.FetchUserById(&p.Context, req)
						if err != nil {
							s.log.Error("Error fetching account: %v", err)
							return nil, fmt.Errorf("error fetching account: %v", err)
						}
						s.log.Info("Account fetched successfully: %v", res.Account.ID)
						return res.Account, nil
					},
				},
				// http://localhost:8083/graphql/account?query={account_by_auth_id(auth_id:1){<entities>}}
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
						res, err := ACCService.FetchUserByAuthID(&p.Context, req)
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

func (s *StructGraphQLResolvers) GetAccountMutationType(ACCService *sv.AccountService) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			/* signup
			http://localhost:8083/graphql/account?query=mutation+_{create_account(
			// account: { email: "", password: "", first_name: "", last_name: "", phone_number: "", address: "", account_type: "", nationality: ""})
			// {access_token, email, auth_id}}
			// */
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
					middleName, _ := accountInput["middle_name"].(string)
					lastName, _ := accountInput["last_name"].(string)
					phoneNumber, _ := accountInput["phone_number"].(string)
					address, _ := accountInput["address"].(string)
					accountType, _ := accountInput["account_type"].(string)
					nationality, _ := accountInput["nationality"].(string)
					birthDate, _ := accountInput["birthdate"].(time.Time)
					nationalID, _ := accountInput["national_id"].(string)
					req := &types.AddAccountRequest{
						Email:       email,
						Password:    password,
						FirstName:   firstName,
						MiddleName:  middleName,
						LastName:    lastName,
						PhoneNumber: phoneNumber,
						Address:     address,
						AccountType: accountType,
						Nationality: nationality,
						BirthDate:   birthDate,
						NationalID:  nationalID,
					}
					res, err := ACCService.CreateUser(&p.Context, req)
					if err != nil {
						s.log.Error("Error creating account: %v", err)
						return nil, fmt.Errorf("error creating account: %v", err)
					}
					s.log.Info("Account created successfully: %v", res.Account.ID)

					return res.Account, nil
				},
			},
			/* login
			http://localhost:8083/graphql/account?query=mutation+_{login(account: { email: "", password: "" }){access_token, email, auth_id}}
			// */
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

					s.log.Info("Login input: %v", loginInput)

					if !ok {
						return nil, fmt.Errorf("invalid account input")
					}
					email, _ := loginInput["email"].(string)
					password, _ := loginInput["password"].(string)
					req := &types.LoginRequest{
						Email:    email,
						Password: password,
					}
					res, err := ACCService.Login(&p.Context, req)
					if err != nil {
						s.log.Error("Error logging in: %v", err)
						return nil, fmt.Errorf("error logging in: %v", err)
					}

					s.log.Info("Login successful: %v", res.AuthID)
					return res, nil
				},
			},
			/* update_account
			http://localhost:8083/graphql/account?query=mutation+_{update_account(
			// account: { email: "", password: "", first_name: "", last_name: "", phone_number: "", address: "", account_type: "",
			*/
			"update_password": &graphql.Field{
				Type:        accountType,
				Description: "Update user password",
				Args: graphql.FieldConfigArgument{
					"UpdatePasswordInput": &graphql.ArgumentConfig{
						Type: types.UpdatePasswordInputType,
					},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					updateInput, ok := p.Args["UpdatePasswordInput"].(map[string]any)
					if !ok {
						return nil, fmt.Errorf("invalid account input")
					}
					auth_id, _ := updateInput["auth_id"].(string)
					oldPassword, _ := updateInput["old_password"].(string)
					newPassword, _ := updateInput["new_password"].(string)
					req := &types.UpdatePasswordRequest{
						AuthID:      auth_id,
						OldPassword: oldPassword,
						NewPassword: newPassword,
					}
					res, err := ACCService.UpdatePassword(&p.Context, req)
					if err != nil {
						s.log.Error("Error updating password: %v", err)
						return nil, fmt.Errorf("error updating password: %v", err)
					}
					s.log.Info("Password updated successfully: %v", res.Account.ID)
					return res.Account, nil
				},
			},
		},
	})
}
