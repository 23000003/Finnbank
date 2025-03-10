package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        token := ctx.GetHeader("Authorization")
        
		// get supabase jwt secret key and use jwt verify token

		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No authentication token"})
		}

        ctx.Next()
    }
}