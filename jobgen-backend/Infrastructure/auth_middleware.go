package infrastructure

import (
	"fmt"
	domain "jobgen-backend/Domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService domain.IJWTService
}

func NewAuthMiddleware(jwtService domain.IJWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		payload, err := a.jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Set user context
		// todo the payload isn't working for the thing that I didn't know check it out at some time
		fmt.Println("user id", payload, payload.UserID, payload.Email, payload.Role)
		c.Set("user_id", payload.UserID)
		c.Set("user_email", payload.Email)
		c.Set("user_role", string(payload.Role))

		c.Next()
	}
}

func (a *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		a.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		role := c.GetString("user_role")
		if role != string(domain.RoleAdmin) {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			c.Abort()
			return
		}

		c.Next()
	})
}

func (a *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				payload, err := a.jwtService.ValidateAccessToken(parts[1])
				if err == nil {
					c.Set("user_id", payload.UserID)
					c.Set("user_email", payload.Email)
					c.Set("user_role", string(payload.Role))
				}
			}
		}
		c.Next()
	}
}
