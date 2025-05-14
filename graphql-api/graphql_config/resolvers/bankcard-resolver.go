package resolvers

import (
	sv "finnbank/graphql-api/services"
	"fmt"
	"github.com/graphql-go/graphql"
)

func (s *StructGraphQLResolvers) GetBankCardQueryType(BCService *sv.BankcardService) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"get_all_bankcard": &graphql.Field{
					Type: graphql.NewList(bankCardType),
					Description: "Get all bank card of user",
					Args: graphql.FieldConfigArgument{
						"user_id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						user_id, ok := p.Args["user_id"].(string)
						if !ok {
							return nil, fmt.Errorf("user_id argument is required and must be an int")
						}
						data, err := BCService.GetAllBankCardOfUserById(p.Context, user_id)
						if err != nil {
							return nil, fmt.Errorf("failed to get bank card: %w", err)
						}
						return data, nil
					},
				},
			},
		},
	)
}

func (s *StructGraphQLResolvers) GetBankCardMutationType(BCService *sv.BankcardService) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"update_bankcard_expiry": &graphql.Field{
					Type:        bankCardType,
					Description: "Update bankcard expiry date",
					Args: graphql.FieldConfigArgument{
						"bankcard_id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						bankcard_id, ok := p.Args["bankcard_id"].(int)
						if !ok {
							return nil, fmt.Errorf("bankcard_id argument is required and must be an int")
						}
						data, err := BCService.UpdateBankcardExpiryDateByUserId(p.Context, bankcard_id)
						if err != nil {
							return nil, fmt.Errorf("failed to update bankcard expiry date: %w", err)
						}
						return data, nil
					},
				},
				"update_pin_number": &graphql.Field{
					Type:        bankCardType,
					Description: "Update bankcard pin number",
					Args: graphql.FieldConfigArgument{
						"bankcard_id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"pin_number": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						bankcard_id, ok := p.Args["bankcard_id"].(int)
						pin_number, ok2 := p.Args["pin_number"].(string)
						if !ok || !ok2 {
							return nil, fmt.Errorf("bankcard_id argument is required and must be an int")
						}
						data, err := BCService.UpdateBankcardPinNumberById(p.Context, bankcard_id, pin_number)
						if err != nil {
							return nil, fmt.Errorf("failed to update bankcard pin number: %w", err)
						}
						return data, nil
					},
				},
			},
		},
	)
}
