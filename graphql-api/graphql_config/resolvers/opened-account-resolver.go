package resolvers

import (
	sv "finnbank/graphql-api/services"
	"fmt"
	"github.com/graphql-go/graphql"
)

func (s *StructGraphQLResolvers) GetOpenedAccountQueryType(OAService *sv.OpenedAccountService, BCService *sv.BankcardService) *graphql.Object {
	
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				// http://localhost:8083/graphql/opened-account?query={get_all(account_id:1){<entities>}}
				"get_all": &graphql.Field{
					Type:        graphql.NewList(openedAccountType),
					Description: "Get all opened accounts by user id",
					Args: graphql.FieldConfigArgument{
						"account_id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Args["account_id"].(string)
						if ok {
							data, err := OAService.GetAllOpenedAccountsByUserId(p.Context, id)
							if err != nil {
								return nil, err
							}
							return data, nil
						}
						return nil, fmt.Errorf("invalid argument: %v", p.Args["account_id"])
					},
				},
				// http://localhost:8083/graphql/opened-account?query={get_by_id(openedaccount_id:1){<entities>}}
				"get_by_id": &graphql.Field{
					Type:        openedAccountType,
					Description: "Get product list",
					Args: graphql.FieldConfigArgument{
						"openedaccount_id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Args["openedaccount_id"].(int)
						if ok {
							data, err := OAService.GetOpenedAccountById(p.Context, id)
							return data, err
						}
						return nil, fmt.Errorf("invalid argument: %v", p.Args["openedaccount_id"])
					},
				},
			"find_by_account_num": &graphql.Field{
				Type:        graphql.Int,
				Description: "Get opened account by account number",
				Args: graphql.FieldConfigArgument{
					"account_number": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					accountNum, ok := p.Args["account_number"].(string)
					if ok {
						data, err := OAService.GetOpenedAccountIdByAccountNumber(p.Context, accountNum)
						return data, err
					}
					return nil, fmt.Errorf("invalid argument: %v", p.Args["account_number"])
				},
			},
			"find_both_account_num": &graphql.Field{
				Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
					Name: "ReceiptAccountNumbers",
					Fields: graphql.Fields{
						"openedaccount_id": &graphql.Field{Type: graphql.Int},
						"account_number":   &graphql.Field{Type: graphql.String},
					},
				})),
				Description: "Get both account number by sender and receiver id",
				Args: graphql.FieldConfigArgument{
					"sender_id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"receiver_id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					senderId, ok := p.Args["sender_id"].(int)
					receiverId, ok1 := p.Args["receiver_id"].(int)

					if ok && ok1 {
						data, err := OAService.GetBothAccountNumberForReceipt(p.Context, senderId, receiverId)
						return data, err
					}
					return nil, fmt.Errorf("invalid argument: %v %v", p.Args["sender_id"], p.Args["receiver_id"])
				},
			},
		},
	})
}

func (s *StructGraphQLResolvers) GetOpenedAccountMutationType(OAService *sv.OpenedAccountService, BCService *sv.BankcardService) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			/* Open a new account
			http://localhost:8083/graphql/opened-account?query=mutation+_{create_account(account_id:1,account_type:"string",balance:1.99){<entities>}}
			// */
			"create_account": &graphql.Field{
				Type:        graphql.NewList(openedAccountType),
				Description: "Open a new account",
				Args: graphql.FieldConfigArgument{
					"account_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					account_id, ok := params.Args["account_id"].(string)

					if ok {
						data, err := OAService.CreateOpenedAccount(params.Context, BCService, account_id);
						return data, err
					}
					
					return nil, fmt.Errorf("invalid argument: %v", params.Args["account_id"]);
				},
			},
			/* Update account Status
			http://localhost:8083/graphql/opened-account?query=mutation+_{update_account_status(openedaccount_id:1,openedaccount_status:"string"){<entities>}}
			// */
			"update_account_status": &graphql.Field{
				Type:        graphql.String,
				Description: "Update openedaccount status",
				Args: graphql.FieldConfigArgument{
					"openedaccount_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"openedaccount_status": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, ok := params.Args["openedaccount_id"].(int)
					status, ok1 := params.Args["openedaccount_status"].(string)

					if ok && ok1 {
						data, err := OAService.UpdateOpenedAccountStatus(params.Context, id, status)
						return data, err
					}
					
					return nil, fmt.Errorf("invalid argument: %v %v", params.Args["openedaccount_id"], params.Args["openedaccount_status"]);
				},
			},
		},
	})
}