package routers

import (
	"github.com/gin-gonic/gin"
	"graphql-service/utils"
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    "graphql-service/docs"
)

func InitializeSwagger(r *gin.Engine) {
	
	logger, err := utils.NewLogger()
    if err != nil {
        panic(err)
    }

	docs.SwaggerInfo.Title = "Bank-It GraphQL-DB Documentation API"
	docs.SwaggerInfo.Description = "This is an api Testing and Documentation for Bank-It services."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8083"
	docs.SwaggerInfo.BasePath = "/api/graphql-service"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info("Swagger running on http://localhost:8083/swagger/index.html")
}