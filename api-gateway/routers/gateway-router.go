package routers

import (
	"gateway/service"
	"github.com/gin-gonic/gin"
)


func GatewayRouter(r *gin.RouterGroup) {

	// new instance of gateway-service
    gatewayService := &service.GatewayService{}

    // routes
    r.GET("/", gatewayService.GetRoot)
    r.GET("/albums", gatewayService.GetAlbums)
    r.GET("/albums/:id", gatewayService.GetAlbumsById)
    r.POST("/hello", gatewayService.PostHello)
    r.POST("/albums", gatewayService.CreateAlbum)
}
