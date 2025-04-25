package resolvers

import (
	"finnbank/common/utils"
	"finnbank/graphql-api/graphql_config/entities"

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

var bankCardType *graphql.Object = entities.GetBankCardEntity()
var productType *graphql.Object = entities.GetProductEntityType()
var accountType *graphql.Object = entities.GetAccountEntityType()
var transactionType *graphql.Object = entities.GetTransactionEntityType()
var notificationType *graphql.Object = entities.GetNotificationEntityType()
var openedAccountType *graphql.Object = entities.GetOpenedAccountEntityType()
var transactionInputType *graphql.InputObject = entities.GetTransactionInputType()
