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

// @Summary Refresh access token
// @Description Refresh the access token using a valid refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} StandardResponse "Tokens refreshed successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Invalid or expired refresh token"
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
    var req RefreshTokenRequest
    
    // First try to get refresh token from request body
    if err := ctx.ShouldBindJSON(&req); err != nil {
        // If no body, try to get from cookie
        refreshTokenFromCookie, err := ctx.Cookie("refresh_token")
        if err != nil || refreshTokenFromCookie == "" {
            ValidationErrorResponse(ctx, err)
            return
        }
        req.RefreshToken = refreshTokenFromCookie
    }
    
    tokens, err := c.authUsecase.RefreshToken(ctx, req.RefreshToken)
    if err != nil {
        switch err {
        case domain.ErrInvalidToken, domain.ErrTokenExpired:
            UnauthorizedResponse(ctx, "Invalid or expired refresh token")
        default:
            InternalErrorResponse(ctx, "Failed to refresh token")
        }
        return
    }
    
    // Update refresh token cookie
    ctx.SetSameSite(http.SameSiteNoneMode)
    ctx.SetCookie(
        "refresh_token",
        tokens.RefreshToken,
        7*24*60*60, // 7 days
        "/",
        "",
        ctx.GetHeader("X-Forwarded-Proto") == "https",
        true, // httpOnly
    )
    
    SuccessResponse(ctx, http.StatusOK, "Tokens refreshed successfully", gin.H{
        "access_token": tokens.AccessToken,
    })
}

// @Summary Logout user
// @Description Logout the current user and invalidate all tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StandardResponse "Logged out successfully"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
    userID := ctx.GetString("user_id")
    if userID == "" {
        UnauthorizedResponse(ctx, "User ID not found")
        return
    }
    
    if err := c.authUsecase.Logout(ctx, userID); err != nil {
        InternalErrorResponse(ctx, "Logout failed")
        return
    }
    
    // Clear refresh token cookie
    ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)
    
    SuccessResponse(ctx, http.StatusOK, "Logged out successfully", nil)
}
