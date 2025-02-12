package routers

import (
	"graphql-service/service"
	"github.com/gin-gonic/gin"
)


func GraphqlRouter(r *gin.RouterGroup) {

	// new instance of graphql-service
    graphqlService := &service.GraphqlService{}

    // routes
    r.GET("/", graphqlService.GetRoot)
    r.GET("/albums", graphqlService.GetAlbums)
    r.GET("/albums/:id", graphqlService.GetAlbumsById)
    r.POST("/hello", graphqlService.PostHello)
    r.POST("/albums", graphqlService.CreateAlbum)
}
