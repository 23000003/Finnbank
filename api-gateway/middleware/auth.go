package middleware

import (
	// "fmt"
	// "net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // token := ctx.GetHeader("Authorization")
        
        // // get supabase jwt secret key and use jwt verify token
        // fmt.Println("Token: ", token)
        // if token == "" || token == "null" || token == "Bearer null" {
        //     ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No authentication token"})
        //     fmt.Println("No authentication token")
        //     ctx.Abort()
        //     return 
        // }

        ctx.Next()
    }
}