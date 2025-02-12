package routers

import (
	"account-service/service"

	"github.com/gin-gonic/gin"
)

func AccountRouter(r *gin.RouterGroup, accountService *service.AccountService) {
	r.GET("/accounts", accountService.GetAccounts)
	r.POST("/register-acc", accountService.AddAccount)
}
