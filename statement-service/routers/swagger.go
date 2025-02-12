package routers

import (
	"github.com/gin-gonic/gin"
	"statement-service/utils"
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    "statement-service/docs"
)

func InitializeSwagger(r *gin.Engine) {
	
	logger, err := utils.NewLogger()
    if err != nil {
        panic(err)
    }

	docs.SwaggerInfo.Title = "Bank-It Statement Service Documentation API"
	docs.SwaggerInfo.Description = "This is an api Testing and Documentation for Bank-It services."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8084"
	docs.SwaggerInfo.BasePath = "/api/statement-service"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info("Swagger running on http://localhost:8084/swagger/index.html")
}