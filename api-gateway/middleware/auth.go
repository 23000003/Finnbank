package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		if idToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authentication token"})
			fmt.Println("No authentication token")
			c.Abort()
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured"})
			fmt.Println("JWT secret not configured")
			c.Abort()
			return
		}

		token, err := jwt.Parse(idToken, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			fmt.Println("Invalid or expired token:", err)
			c.Abort()
			return
		}
		c.Next()
	}
}
