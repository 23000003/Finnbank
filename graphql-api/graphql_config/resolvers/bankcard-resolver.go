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
			Name: "get_requests",
			Fields: graphql.Fields{
				"request_by_id": &graphql.Field{
					Type:        bcen.GetBankCardResponseEntity(),
					Description: "Get all bank card requests by id",
					Args: graphql.FieldConfigArgument{
						"request_id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Args["request_id"].(int)
						if !ok {
							return nil, fmt.Errorf("request_id argument is required and must be an int")
						}
						data, err := BCService.GetBankCardRequestsById(p.Context, id)
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
			Name: "create_request",
			Fields: graphql.Fields{
				"insert_card_requester": &graphql.Field{
					Type:        bcen.GetBankCardRequestEntity(),
					Description: "Create a new bank card request",
					Args: graphql.FieldConfigArgument{
						"first_name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"last_name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"card_type": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						fname, okfname := p.Args["first_name"].(string)
						lname, oklname := p.Args["last_name"].(string)
						ctype, okctype := p.Args["card_type"].(int)

						if !okfname {
							return nil, fmt.Errorf("first_name argument is required and must be an int")
						}

						if !oklname {
							return nil, fmt.Errorf("last_name argument is required and must be an int")
						}

						if !okctype {
							return nil, fmt.Errorf("card_type argument is required and must be an int")
						}

						data, err := BCService.CreateCardRequest(p.Context, fname, lname, ctype)

						if err != nil {
							s.log.Error("Error in Creating Card Request: %v", err)
							return nil, fmt.Errorf("Error in Creating Card Request: %v", err)
						}
						s.log.Info("Request Created Successfully: %v", fname)

						return data, nil
					},
				},
			},
		},
	)
}
