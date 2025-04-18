package resolvers

import (
	sv "finnbank/graphql-api/services"
	"fmt"
	"github.com/graphql-go/graphql"
)

func (s *StructGraphQLResolvers) GetOpenedAccountQueryType(OAService *sv.OpenedAccountService) *graphql.Object {
	
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "get_all",
			Fields: graphql.Fields{
				// http://localhost:8083/graphql/opened-account?query={opened_account(account_id:1){<entities>}}
				"opened_account": &graphql.Field{
					Type:        graphql.NewList(openedAccountType),
					Description: "Get all opened accounts by user id",
					Args: graphql.FieldConfigArgument{
						"account_id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Args["account_id"].(int)
						if ok {
							data, err := OAService.GetAllOpenedAccountsByUserId(p.Context, id)
							return data, err
						}
						return nil, fmt.Errorf("invalid argument: %v", p.Args["account_id"])
					},
				},
				// http://localhost:8083/graphql/product?query={list{id,name,info,price}}
				"list": &graphql.Field{
					Type:        graphql.NewList(productType),
					Description: "Get product list",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {

						
						return nil, nil
					},
				},
			},
		},
	)
}

func (s *StructGraphQLResolvers) GetOpenedAccountMutationType(OAService *sv.OpenedAccountService) *graphql.Object {
	return nil
}