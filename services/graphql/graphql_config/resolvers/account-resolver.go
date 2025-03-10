package resolvers

import (
	"finnbank/services/common/grpc/account"

	"github.com/graphql-go/graphql"
)

// FUTURE: GET THIS DONE
func (s *StructGraphQLResolvers) GetAccountQueryType(accServer account.AccountServiceClient) *graphql.Object {

}

func (s *StructGraphQLResolvers) GetAccountMutationType(accServer account.AccountServiceClient) *graphql.Object {

}
