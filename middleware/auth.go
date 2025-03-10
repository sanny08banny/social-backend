package middleware

import (
	"net/http"
	"social-backend/auth"
	"strings"
	"github.com/gin-gonic/gin"
)

// JWT Middleware
// func JWTAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
// 			c.Abort()
// 			return
// 		}

// 		// Extract token from "Bearer <token>"
// 		tokenParts := strings.Split(authHeader, " ")
// 		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
// 			c.Abort()
// 			return
// 		}
// 		tokenString := tokenParts[1]

// 		// Validate token
// 		claims, err := auth.ValidateJWT(tokenString)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
// 			c.Abort()
// 			return
// 		}

// 		// Store user ID in context
// 		c.Set("user_id", uint(claims["user_id"].(float64)))
// 		c.Next()
// 	}
// }
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// If no token is provided, allow access as a guest
		if authHeader == "" {
			c.Set("user_id", nil) // No user authenticated
			c.Next()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		tokenString := tokenParts[1]

		// Validate token
		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store user ID in context if token is valid
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Next()
	}
}

