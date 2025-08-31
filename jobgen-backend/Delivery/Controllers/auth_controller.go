package controllers

import (
	domain "jobgen-backend/Domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
    authUsecase domain.IAuthUsecase
}

func NewAuthController(authUsecase domain.IAuthUsecase) *AuthController {
    return &AuthController{authUsecase: authUsecase}
}

type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
    var req RefreshTokenRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    tokens, err := c.authUsecase.RefreshToken(ctx, req.RefreshToken)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{
        "access_token":  tokens.AccessToken,
        "refresh_token": tokens.RefreshToken,
    })
}

func (c *AuthController) Logout(ctx *gin.Context) {
    userID := ctx.GetString("user_id")
    if userID == "" {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
        return
    }
    if err := c.authUsecase.Logout(ctx, userID); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
        return
    }
    ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)
    ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
