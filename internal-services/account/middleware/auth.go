package middleware

import (
	"finnbank/internal-services/account/helpers"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// These middleware code can be reused for later

func FetchUserUUID() gin.HandlerFunc {
	return func(c *gin.Context) {
		encodedEmail := c.Param("email")
		email, err := url.QueryUnescape(encodedEmail)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			c.Abort()
			return
		}

		uuid, err := helpers.GetUserUUIDByEmail(email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Storing UUID in Gin context for later use
		c.Set("email", email)
		c.Set("uuid", uuid)

		c.Next()
	}
}

func DeleteAuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieving UUID from context
		uuid, exists := c.Get("uuid")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID not found in context"})
			c.Abort()
			return
		}

		err := helpers.DeleteUserByUUID(uuid.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user from Auth: " + err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
