package routers

import (
	"security-service/service"
	"github.com/gin-gonic/gin"
)


func SecurityRouter(r *gin.RouterGroup) {

	// new instance of security-service
    securityService := &service.SecurityService{}

    // routes
    r.GET("/", securityService.GetRoot)
    r.GET("/albums", securityService.GetAlbums)
    r.GET("/albums/:id", securityService.GetAlbumsById)
    r.POST("/hello", securityService.PostHello)
    r.POST("/albums", securityService.CreateAlbum)
}
