package routers

import (
	"bankcard-service/service"
	"github.com/gin-gonic/gin"
)


func BankcardRouter(r *gin.RouterGroup) {

	// new instance of bankcard-service
    bankcardService := &service.BankcardService{}

    // routes
    r.GET("/", bankcardService.GetRoot)
    r.GET("/albums", bankcardService.GetAlbums)
    r.GET("/albums/:id", bankcardService.GetAlbumsById)
    r.POST("/hello", bankcardService.PostHello)
    r.POST("/albums", bankcardService.CreateAlbum)
}
