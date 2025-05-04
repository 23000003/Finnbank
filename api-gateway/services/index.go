package services

import (
	"finnbank/common/utils"
	"finnbank/api-gateway/types"
)

type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors interface{} `json:"errors"`
}

func NewApiGatewayServices(log *utils.Logger) *types.ApiGatewayServices {
	return &types.ApiGatewayServices{
		ProductService:      newProductService(log),
		AccountService:      newAccountService(log),
		StatementService:    newStatementService(log),
		TransactionService:  newTransactionService(log),
		BankcardService:     newBankcardService(log),
		NotificationService: newNotificationService(log),
		OpenedAccountService: newOpenedAccountService(log),
		RealTimeService: 		 newRealTimeService(log),
	}
}
