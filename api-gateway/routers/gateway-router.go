package routers

import (
	"finnbank/api-gateway/middleware"
	"finnbank/api-gateway/types"
	"finnbank/common/utils"

	"github.com/gin-gonic/gin"
)

type StructGatewayRouter struct {
	s *types.ApiGatewayServices
	r *gin.RouterGroup
	l *utils.Logger
}

func NewGatewayRouter(l *utils.Logger, r *gin.RouterGroup, s *types.ApiGatewayServices) *StructGatewayRouter {
	return &StructGatewayRouter{
		r: r,
		l: l,
		s: s,
	}
}

func (gr *StructGatewayRouter) ConfigureGatewayRouter() {

	// ======================== TEST ========================
	product := gr.r.Group("/product")
	{
		product.GET("", gr.s.ProductService.GetAllProduct)
		product.GET("/:id", gr.s.ProductService.GetByIdProduct)
		product.POST("", gr.s.ProductService.CreateProduct)
		product.PATCH("/:id", gr.s.ProductService.UpdateProduct)
		product.DELETE("/:id", gr.s.ProductService.DeleteProduct)
	}

	// ======================== Services ========================

	auth := gr.r.Group("/auth")
	{
		auth.POST("/login", gr.s.AccountService.LoginUser)
		auth.POST("/signup", gr.s.AccountService.SignupUser)
	}
	account := gr.r.Group("/account")
	account.Use(middleware.AuthMiddleware())
	{
		// These could be used for a potential "search user by" functions in the frontend
		account.GET("/get-user-by-id/:id", gr.s.AccountService.GetUserAccountById)
		account.GET("/get-user-by-email/:email", gr.s.AccountService.GetUserAccountByEmail)
		account.GET("/get-user-by-account-number/:account_number", gr.s.AccountService.GetUserAccountByAccountNumber)
		account.PATCH("/update-password", gr.s.AccountService.UpdateUserPassword)
		account.PATCH("/update-user", gr.s.AccountService.UpdateUser)
		account.PATCH("/update-user-details", gr.s.AccountService.UpdateUserDetails)
		account.PATCH("/update-account-status", gr.s.AccountService.UpdateAccountStatus)
	}

	statement := gr.r.Group("/statement")
	statement.Use(middleware.AuthMiddleware())
	{
		statement.GET("/generate-statement", gr.s.StatementService.GenerateStatement)
	}

	bankcard := gr.r.Group("/bankcard")
	bankcard.Use(middleware.AuthMiddleware())
	{
		bankcard.GET("/get-all-bankcard/:id", gr.s.BankcardService.GetAllBankCardOfUserById)
		bankcard.PATCH("/renew-bankcard/:id", gr.s.BankcardService.UpdateBankcardExpiryDateByUserId)
		bankcard.PATCH("/update-pin-number/:id/:new-pin", gr.s.BankcardService.UpdateBankcardPinNumberById)
	}

	transaction := gr.r.Group("/transaction")
	transaction.Use(middleware.AuthMiddleware())
	{
		transaction.GET("/get-all", gr.s.TransactionService.GetTransactionByOpenAccountId)
		transaction.GET("/get-all-by-timestamp", gr.s.TransactionService.GetTransactionByTimestamp)
		transaction.POST("/generate-transaction", gr.s.TransactionService.CreateTransaction)
	}

	notification := gr.r.Group("/notification")
	notification.Use(middleware.AuthMiddleware())
	{
		notification.GET("/get-all/:id/:limit", gr.s.NotificationService.GetAllNotificationByUserId)
		notification.GET("/get-all-unread/:id", gr.s.NotificationService.GetAllUnreadNotificationByUserId)
		notification.GET("/get-one/:id", gr.s.NotificationService.GetNotificationByUserId)
		notification.POST("/generate-notif", gr.s.NotificationService.GenerateNotification)
		notification.PATCH("/mark-as-read/:id", gr.s.NotificationService.ReadNotificationByUserId)
	}

	openedAccount := gr.r.Group("/opened-account")
	openedAccount.Use(middleware.AuthMiddleware())
	{
		openedAccount.GET("/get-all/:id", gr.s.OpenedAccountService.GetAllOpenedAccountsByUserId)
		openedAccount.GET("/:id", gr.s.OpenedAccountService.GetOpenedAccountOfUserById)
		openedAccount.GET("/get-user-id/:id", gr.s.OpenedAccountService.GetUserIdByOpenedAccountId)
		openedAccount.GET("/get-both-account-number/:sent_id/:receive_id", gr.s.OpenedAccountService.GetBothAccountNumberForReceipt)
		openedAccount.GET("/find-by-account-number/:acc_num", gr.s.OpenedAccountService.GetOpenedAccountIdByAccountNumber)
		openedAccount.POST("/create-account", gr.s.OpenedAccountService.OpenAnAccountByAccountType)
		openedAccount.PATCH("/update-status/:id/:status", gr.s.OpenedAccountService.UpdateOpenedAccountStatus)
	}

	realTime := gr.r.Group("/ws")
	realTime.Use(middleware.AuthMiddleware())
	{
		realTime.GET("/listen-to-notification", gr.s.RealTimeService.GetRealTimeNotification)
		realTime.GET("/listen-to-transaction", gr.s.RealTimeService.GetRealTimeTransaction)
	}

}
