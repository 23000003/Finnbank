package resolvers

import (
	"github.com/graphql-go/graphql"
	"finnbank/graphql-api/services"
)

func (s *StructGraphQLResolvers) GetStatementQueryType(ST *services.StatementService) *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				// http://localhost:8083/graphql/statement?query={generate_statement(credit_id:1, debit_id:2, savings_id:3, start_date:"HEY", end_date:"HEY"){pdf_buffer}}
				"generate_statement": &graphql.Field{
					Type: statementType,
					Description: "Get all bank card of user",
					Args: graphql.FieldConfigArgument{
						"credit_id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"debit_id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"savings_id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"start_date": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"end_date": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						creditId := p.Args["credit_id"].(int)
						debitId := p.Args["debit_id"].(int)
						savingsId := p.Args["savings_id"].(int)
						startDate := p.Args["start_date"].(string)
						endDate := p.Args["end_date"].(string)

						data, err := ST.GenerateStatementForUser(creditId, debitId, savingsId, startDate, endDate, &p.Context)
						if err != nil {
							return nil, err
						}
						return map[string]interface{}{
							"pdf_buffer": data,
						}, nil
					},
				},
			},
		},
	)
}
