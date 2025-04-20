package types

import (
	"github.com/gin-gonic/gin"
)

type ApiGatewayServices struct {
	ProductService IProductService
	AccountService IAccountService
	StatementService IStatementService
	TransactionService ITransactionService
	BankcardService	IBankcardService
	NotificationService INotificationService
	OpenedAccountService IOpenedAccountService
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
}

type IStatementService interface {
	GenerateStatement(*gin.Context)  // post req
}

type ITransactionService interface {
	GetAllTransaction(*gin.Context) 
	GetTransaction(*gin.Context) 
	GenerateTransaction(*gin.Context)  // post req
}

type IBankcardService interface {
	GetUserBankcard(*gin.Context) 
	GenerateBankcardForUser(*gin.Context)  // post req
	RenewBankcardForUser(*gin.Context)  // update req
}

type INotificationService interface {
	GetUserNotifications(*gin.Context) 
	GenerateNotification(*gin.Context)  // post req
	MarkAsReadNotification(*gin.Context)  // update req
	DeleteNotification(*gin.Context) 
}

type IOpenedAccountService interface {
	GetAllOpenedAccountsByUserId(*gin.Context)
	GetOpenedAccountOfUserById(*gin.Context)
	OpenAnAccountByAccountType(*gin.Context)
	UpdateOpenedAccountStatus(*gin.Context)
}
