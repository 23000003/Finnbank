package types

import (
	"github.com/gin-gonic/gin"
)

type ApiGatewayServices struct {
	ProductService       IProductService
	AccountService       IAccountService
	StatementService     IStatementService
	TransactionService   ITransactionService
	BankcardService      IBankcardService
	NotificationService  INotificationService
	OpenedAccountService IOpenedAccountService
	RealTimeService     IRealTimeService
}

type IProductService interface {
	GetAllProduct(*gin.Context)
	GetByIdProduct(*gin.Context)
	CreateProduct(*gin.Context)
	UpdateProduct(*gin.Context)
	DeleteProduct(*gin.Context)
}

type IAccountService interface {
	LoginUser(*gin.Context)
	SignupUser(*gin.Context)
	GetUserAccountById(*gin.Context)
	GetUserAccountByAccountNumber(*gin.Context)
	GetUserAccountByEmail(*gin.Context)
	UpdateUserPassword(*gin.Context)
	UpdateUserDetails(*gin.Context)
	UpdateUser(c *gin.Context)
}

type IStatementService interface {
	GenerateStatement(*gin.Context) // post req
}

type ITransactionService interface {
	GetTransactionByOpenAccountId(*gin.Context)
	CreateTransaction(*gin.Context) // post req
}

type IBankcardService interface {
	GetAllBankCardOfUserById(*gin.Context)
	UpdateBankcardExpiryDateByUserId(*gin.Context) // update req
	UpdateBankcardPinNumberById(*gin.Context)
}

type INotificationService interface {
	GetAllNotificationByUserId(*gin.Context)
	GetAllUnreadNotificationByUserId(*gin.Context)
	GetNotificationByUserId(*gin.Context)
	GenerateNotification(*gin.Context)     // post req
	ReadNotificationByUserId(*gin.Context) // update req
}

type IOpenedAccountService interface {
	GetAllOpenedAccountsByUserId(*gin.Context)
	GetOpenedAccountOfUserById(*gin.Context)
	GetOpenedAccountIdByAccountNumber(*gin.Context)
	GetBothAccountNumberForReceipt(*gin.Context)
	OpenAnAccountByAccountType(*gin.Context)
	UpdateOpenedAccountStatus(*gin.Context)
	GetUserIdByOpenedAccountId(*gin.Context)
}

type IRealTimeService interface {
	GetRealTimeTransaction(*gin.Context)
	GetRealTimeNotification(*gin.Context)
	Shutdown()
}