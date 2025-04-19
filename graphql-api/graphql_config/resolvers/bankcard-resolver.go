package resolvers

import (
	sv "finnbank/graphql-api/services"

	"github.com/graphql-go/graphql"
)

var bankCardType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "BankCard",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.Int,
				Description: "The ID of the bank card",
			},
			"card_number": &graphql.Field{
				Type:        graphql.String,
				Description: "The number of the bank card",
			},
			"expiry_date": &graphql.Field{
				Type:        graphql.String,
				Description: "The expiry date of the bank card",
			},
		},
	},
)

func (s *StructGraphQLResolvers) GetBankCardQueryType(BCService *sv.BankcardService) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "get_all",
			Fields: graphql.Fields{
				"request_by_id": &graphql.Field{
					Type:        graphql.NewList(bankCardType),
					Description: "Get all bank card requests by id",
					Args: graphql.FieldConfigArgument{
						"request_id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Args["request_id"].(int)
						if ok {
							data, err := BCService.GetAllBankCardRequestsById(p.Context, id)
							return data, err
						}
						return nil, nil
					},
				},
			},
		},
	)
}

func (s *StructGraphQLResolvers) GetBankCardMutationType(BCService *sv.BankcardService) *graphql.Object {
	return nil
}
