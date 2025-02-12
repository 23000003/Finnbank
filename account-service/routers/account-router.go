package routers

import (
	"account-service/service"
	"github.com/gin-gonic/gin"
)


func AccountRouter(r *gin.RouterGroup) {

	// new instance of account-service
    accountService := &service.AccountService{}

    // routes
    r.GET("/", accountService.GetRoot)
    r.GET("/albums", accountService.GetAlbums)
    r.GET("/albums/:id", accountService.GetAlbumsById)
    r.POST("/hello", accountService.PostHello)
    r.POST("/albums", accountService.CreateAlbum)
}
