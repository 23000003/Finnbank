package routers

import (
	"account-service/service"

	"github.com/gin-gonic/gin"
)

func AccountRouter(r *gin.RouterGroup, accountService *service.AccountService) {
	r.GET("/accounts", accountService.GetAccounts)
	r.GET("/fetch-acc/:acc_num", accountService.GetAccoutById)
	r.POST("/register", accountService.AddAccount)
	r.PATCH("/update_acc/:acc_num", accountService.UpdateHasCard)
	r.DELETE("/delete-user/:email", accountService.DeleteUser)
}
