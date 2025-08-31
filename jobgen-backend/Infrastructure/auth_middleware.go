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

// StandardResponse for middleware errors
type MiddlewareErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, MiddlewareErrorResponse{
				Success: false,
				Message: "Authentication required",
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "MISSING_AUTH_HEADER",
					Message: "Authorization header is required",
				},
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, MiddlewareErrorResponse{
				Success: false,
				Message: "Invalid authorization format",
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "INVALID_AUTH_FORMAT",
					Message: "Authorization header must be in format: Bearer <token>",
				},
			})
			c.Abort()
			return
		}

		token := strings.TrimSpace(parts[1])
		if token == "" {
			c.JSON(http.StatusUnauthorized, MiddlewareErrorResponse{
				Success: false,
				Message: "Empty token",
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "EMPTY_TOKEN",
					Message: "Token cannot be empty",
				},
			})
			c.Abort()
			return
		}

		payload, err := a.jwtService.ValidateAccessToken(token)
		if err != nil {
			errorMessage := "Invalid or expired token"
			errorCode := "INVALID_TOKEN"
			
			switch err {
			case domain.ErrTokenExpired:
				errorMessage = "Token has expired"
				errorCode = "TOKEN_EXPIRED"
			case domain.ErrInvalidToken:
				errorMessage = "Invalid token"
				errorCode = "INVALID_TOKEN"
			}
			
			c.JSON(http.StatusUnauthorized, MiddlewareErrorResponse{
				Success: false,
				Message: "Authentication failed",
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    errorCode,
					Message: errorMessage,
				},
			})
			c.Abort()
			return
		}

		// Validate payload
		if payload == nil || payload.UserID == "" {
			c.JSON(http.StatusUnauthorized, MiddlewareErrorResponse{
				Success: false,
				Message: "Invalid token payload",
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "INVALID_TOKEN_PAYLOAD",
					Message: "Token payload is invalid or missing user information",
				},
			})
			c.Abort()
			return
		}

		// Set user context
		fmt.Println("Payload Data", payload)
		c.Set("user_id", payload.UserID)
		c.Set("user_email", payload.Email)
		c.Set("user_role", string(payload.Role))

		c.Next()
	}
}

func (a *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If the request was already aborted upstream, do nothing.
		if c.IsAborted() {
			return
		}
		// Try to get role from context (maybe RequireAuth ran earlier)
		role := c.GetString("user_role")
		fmt.Println("Admin required route (initial role):", role)

		// If role is empty, attempt to extract & validate token (do not rely on RequireAuth)
		if role == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, MiddlewareErrorResponse{
					Success: false,
					Message: "Authentication required",
					Error: struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					}{
						Code:    "MISSING_AUTH_HEADER",
						Message: "Authorization header is required",
					},
				})
				c.Abort()
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				c.JSON(http.StatusUnauthorized, MiddlewareErrorResponse{
					Success: false,
					Message: "Invalid authorization format",
					Error: struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					}{
						Code:    "INVALID_AUTH_FORMAT",
						Message: "Authorization header must be in format: Bearer <token>",
					},
				})
				c.Abort()
				return
			}

			token := strings.TrimSpace(parts[1])
			if token == "" {
				c.JSON(http.StatusUnauthorized, MiddlewareErrorResponse{
					Success: false,
					Message: "Empty token",
					Error: struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					}{
						Code:    "EMPTY_TOKEN",
						Message: "Token cannot be empty",
					},
				})
				c.Abort()
				return
			}

			payload, err := a.jwtService.ValidateAccessToken(token)
			if err != nil {
				errorMessage := "Invalid or expired token"
				errorCode := "INVALID_TOKEN"

				switch err {
				case domain.ErrTokenExpired:
					errorMessage = "Token has expired"
					errorCode = "TOKEN_EXPIRED"
				case domain.ErrInvalidToken:
					errorMessage = "Invalid token"
					errorCode = "INVALID_TOKEN"
				}

				c.JSON(http.StatusUnauthorized, MiddlewareErrorResponse{
					Success: false,
					Message: "Authentication failed",
					Error: struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					}{
						Code:    errorCode,
						Message: errorMessage,
					},
				})
				c.Abort()
				return
			}

			// Validate payload
			if payload == nil || payload.UserID == "" {
				c.JSON(http.StatusUnauthorized, MiddlewareErrorResponse{
					Success: false,
					Message: "Invalid token payload",
					Error: struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					}{
						Code:    "INVALID_TOKEN_PAYLOAD",
						Message: "Token payload is invalid or missing user information",
					},
				})
				c.Abort()
				return
			}

			// Set user info into context so subsequent handlers/middlewares can use it
			c.Set("user_id", payload.UserID)
			c.Set("user_email", payload.Email)
			c.Set("user_role", string(payload.Role))

			role = string(payload.Role)
			fmt.Println("Admin required route (role from token):", role)
		}

		// Now check role
		if role != string(domain.RoleAdmin) {
			c.JSON(http.StatusForbidden, MiddlewareErrorResponse{
				Success: false,
				Message: "Admin access required",
				Error: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "INSUFFICIENT_PERMISSIONS",
					Message: "This endpoint requires admin privileges",
				},
			})
			c.Abort()
			return
		}

		// Role is admin — continue the chain
		c.Next()
	}
}


func (a *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				token := strings.TrimSpace(parts[1])
				if token != "" {
					payload, err := a.jwtService.ValidateAccessToken(token)
					if err == nil && payload != nil {
						c.Set("user_id", payload.UserID)
						c.Set("user_email", payload.Email)
						c.Set("user_role", string(payload.Role))
					}
				}
			}
		}
		c.Next()
	}
}
