package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wetrade/tokens"
)

func Authenticate(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Token not found"})
			c.Abort()
			return
		}

		claims, err := tokens.VerifyToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			c.Abort()
			return
		}

		log.Println("User Role:", claims.Role)
		log.Println("Allowed Roles:", allowedRoles)

		if !authorizedRoles(claims.Role, allowedRoles) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to access this resource 😜"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("userId", claims.Id)
		c.Set("role", claims.Role)
		c.Next()
	}

}

func authorizedRoles(userRole string, allowedRoles []string) bool {
	for _, role := range allowedRoles {
		if role == userRole {
			return true
		}
	}
	return false
}
