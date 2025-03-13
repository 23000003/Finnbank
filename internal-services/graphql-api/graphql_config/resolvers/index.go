package resolvers

import (
	"finnbank/common/utils"
	"finnbank/internal-services/graphql-api/graphql_config/entities"

	"github.com/graphql-go/graphql"
)

type StructGraphQLResolvers struct {
	log *utils.Logger
}

func NewGraphQLResolvers(log *utils.Logger) *StructGraphQLResolvers {
	return &StructGraphQLResolvers{
		log: log,
	}
}

var productType *graphql.Object = entities.GetProductEntityType()

// === Unused Warning ==
// var accountType *graphql.Object = entities.GetAccountEntityType()
// var bankCardType *graphql.Object = entities.GetBankCardEntityType()
// var transaction_type *graphql.Object = entities.GetTransactionEntityType()
// var notification_type *graphql.Object = entities.GetNotificationEntityType()
