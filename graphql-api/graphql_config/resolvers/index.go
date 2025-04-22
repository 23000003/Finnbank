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

var productType *graphql.Object = entities.GetProductEntityType()
var accountType *graphql.Object = entities.GetAccountEntityType()
var openedAccountType *graphql.Object = entities.GetOpenedAccountEntityType()
<<<<<<< HEAD
var bankCardType *graphql.Object = entities.GetBankCardEntityType()
=======
var notificationType *graphql.Object = entities.GetNotificationEntityType()
var transactionType *graphql.Object = entities.GetTransactionEntityType()
var transactionInputType *graphql.InputObject = entities.GetTransactionInputType()
>>>>>>> ef2c68d1eaf28510b3055a948a19f52f82245ad6

// === Unused Warning ==
// var transaction_type *graphql.Object = entities.GetTransactionEntityType()
