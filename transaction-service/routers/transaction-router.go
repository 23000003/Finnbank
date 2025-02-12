package routers

import (
	"transaction-service/service"
	"github.com/gin-gonic/gin"
)


func TransactionRouter(r *gin.RouterGroup) {

	// new instance of transaction-service
    transactionService := &service.TransactionService{}

    // routes
    r.GET("/", transactionService.GetRoot)
    r.GET("/albums", transactionService.GetAlbums)
    r.GET("/albums/:id", transactionService.GetAlbumsById)
    r.POST("/hello", transactionService.PostHello)
    r.POST("/albums", transactionService.CreateAlbum)
}
