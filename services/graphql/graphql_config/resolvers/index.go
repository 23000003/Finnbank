package resolvers

import (
	"github.com/graphql-go/graphql"
)


func GetQueryType() *graphql.Object {
	return queryType
}
func GetMutationType() *graphql.Object {
	return mutationType
}