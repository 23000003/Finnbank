package handlers

/*
	Refactor to http handler file
*/

import (
	"finnbank/services/account/middleware"
	"finnbank/services/account/service"
	"github.com/gin-gonic/gin"
)

func AccountRouter(r *gin.RouterGroup, accountService *service.AccountService) {
	r.GET("/accounts", accountService.GetAccounts)
	r.GET("/fetch-acc/:acc_num", accountService.GetAccoutById)

	// Adding new User
	r.POST("/register", accountService.AddAccount)

	// Update has_card field
	r.PATCH("/update_acc/:acc_num", accountService.UpdateHasCard)

	// Delete User
	r.DELETE("/delete-user/:email",
		middleware.FetchUserUUID(),
		middleware.DeleteAuthUser(),
		accountService.DeleteUser,
	)
}
