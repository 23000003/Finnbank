package routers

import (
	"statement-service/service"
	"github.com/gin-gonic/gin"
)


func StatementRouter(r *gin.RouterGroup) {

	// new instance of statement-service
    statementService := &service.StatementService{}

    // routes
    r.GET("/", statementService.GetRoot)
    r.GET("/albums", statementService.GetAlbums)
    r.GET("/albums/:id", statementService.GetAlbumsById)
    r.POST("/hello", statementService.PostHello)
    r.POST("/albums", statementService.CreateAlbum)
}
