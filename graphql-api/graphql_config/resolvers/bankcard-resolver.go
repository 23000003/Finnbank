package resolvers

import (
	bcen "finnbank/graphql-api/graphql_config/entities"
	sv "finnbank/graphql-api/services"
	"fmt"

	"github.com/graphql-go/graphql"
)

func (s *StructGraphQLResolvers) GetBankCardQueryType(BCService *sv.BankcardService) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"bankcard_by_bankcard_number": &graphql.Field{
					Type:        bcen.GetBankCardEntity(),
					Description: "Get all bank cards by bankcard number",
					Args: graphql.FieldConfigArgument{
						"bankcard_number": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						bankcard_number, ok := p.Args["bankcard_number"].(string)
						if !ok {
							return nil, fmt.Errorf("bankcard_number argument is required and must be an int")
						}
						data, err := BCService.GetBankCardByNumber(p.Context, bankcard_number)
						return data, err
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
				"create_bankcard_function": &graphql.Field{
					Type:        bcen.GetBankCardEntity(),
					Description: "Create a new bank card",
					Args: graphql.FieldConfigArgument{
						"first_name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"last_name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"cardtype": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"account_holder_id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						fname, okfname := p.Args["first_name"].(string)
						lname, oklname := p.Args["last_name"].(string)
						ctype, okctype := p.Args["cardtype"].(string)
						cardholder_id, okcardholder_id := p.Args["account_holder_id"].(string)

						if !okfname {
							return nil, fmt.Errorf("first_name argument is required and must be an string")
						}

						if !oklname {
							return nil, fmt.Errorf("last_name argument is required and must be an string")
						}

						if !okctype {
							return nil, fmt.Errorf("card_type argument is required and must be an string")
						}

						if !okcardholder_id {
							return nil, fmt.Errorf("cardholder_id argument is required and must be an int")
						}

						data, err := BCService.CreateBankCardForUser(p.Context, ctype, cardholder_id, fname, lname)

						if err != nil {
							s.log.Error("Error in Creating Card Request: %v", err)
							return nil, fmt.Errorf("error in Creating Card Request: %v", err)
						}
						s.log.Info("Request Created Successfully: %v", fname)

						return data, nil
					},
				},
			},
		},
	)
}
