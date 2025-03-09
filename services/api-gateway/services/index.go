package services

import (
	"finnbank/services/api-gateway/types"
	"finnbank/services/common/utils"
)

type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors interface{} `json:"errors"`
}

func NewApiGatewayServices (log *utils.Logger) *types.ApiGatewayServices {
	return &types.ApiGatewayServices{
		ProductService: NewProductService(log),
		AccountService: NewAccountService(log),
		StatementService: NewStatementService(log),
		TransactionService: NewTransactionService(log),
		BankcardService: NewBankcardService(log),
		NotificationService: NewNotificationService(log),
	}
}
