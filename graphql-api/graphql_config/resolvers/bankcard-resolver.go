package resolvers

import (
	bcen "finnbank/graphql-api/graphql_config/entities"
	sv "finnbank/graphql-api/services"
	"fmt"

	"github.com/google/uuid"

	"github.com/graphql-go/graphql"
)

func (s *StructGraphQLResolvers) GetBankCardQueryType(BCService *sv.BankcardService) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
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
			Name: "Mutation",
			Fields: graphql.Fields{
				"create_request_function": &graphql.Field{
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
						ctype, okctype := p.Args["card_type"].(string)

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
							return nil, fmt.Errorf("error in creating card request: %v", err)
						}
						s.log.Info("Request Created Successfully: %v", fname)

						return data, nil
					},
				},
				"create_bankcard_function": &graphql.Field{
					Type:        bcen.GetBankCardResponseEntity(),
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
						account_holder_id, okaccount_holder_id := p.Args["account_holder_id"].(string)

						if !okfname {
							return nil, fmt.Errorf("first_name argument is required and must be an string")
						}

						if !oklname {
							return nil, fmt.Errorf("last_name argument is required and must be an string")
						}

						if !okctype {
							return nil, fmt.Errorf("card_type argument is required and must be an string")
						}

						if !okaccount_holder_id {
							return nil, fmt.Errorf("account_holder_id argument is required and must be an string")
						}

						accountHolderUUID, err := uuid.Parse(account_holder_id)
						if err != nil {
							return nil, fmt.Errorf("invalid account_holder_id: %v", err)
						}
						data, err := BCService.CreateBankCardForUser(p.Context, fname, lname, ctype, accountHolderUUID)

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
