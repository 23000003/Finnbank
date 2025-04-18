package resolvers

import (
	"finnbank/graphql-api/services"

	"github.com/graphql-go/graphql"
)

func (s *StructGraphQLResolvers) GetBankCardQueryType(db interface{}) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "BankCardQuery",
		Fields: graphql.Fields{
			"getBankCardByUserId": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "BankCard",
					Fields: graphql.Fields{
						"BankCard_ID": &graphql.Field{Type: graphql.Int},
						"Card_Number": &graphql.Field{Type: graphql.String},
						"Expiry":      &graphql.Field{Type: graphql.DateTime},
						"Card_Type":   &graphql.Field{Type: graphql.String},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"account_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					accountID := p.Args["account_id"].(int)
					return services.GetBankCardOfUserById(db, accountID)
				},
			},
		},
	})
}

func (s *StructGraphQLResolvers) GetBankCardMutationType(db interface{}) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "BankCardMutation",
		Fields: graphql.Fields{
			"createBankCard": &graphql.Field{
				Type: graphql.String, // Return success message or ID
				Args: graphql.FieldConfigArgument{
					"account_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"card_type":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"expiry":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return services.CreateBankCardForUser(db, p.Args)
				},
			},
			"updateExpiry": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"account_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"new_expiry": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return services.UpdateBankcardExpiryDateByUserId(db, p.Args)
				},
			},
		},
	})
}
