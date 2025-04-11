package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

func (s *authService) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			log.Println("Missing authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// More robust Bearer token extraction
		tokenString := ""
		if strings.HasPrefix(authHeader, BearerPrefix) {
			tokenString = strings.TrimSpace(authHeader[len(BearerPrefix):])
		} else {
			log.Println("Invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header must start with 'Bearer '"})
			return
		}

		// Token validation
		claims, err := s.ParseToken(tokenString)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":  "invalid token",
				"detail": err.Error(), // Only include in development
			})
			return
		}

		log.Printf("Authenticated user: %s (ID: %s)", claims.Email, claims.UserID)
		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
