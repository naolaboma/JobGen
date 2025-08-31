package controllers

import (
	domain "jobgen-backend/Domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase domain.IUserUsecase
}

func NewUserController(userUsecase domain.IUserUsecase) *UserController {
	return &UserController{
		userUsecase: userUsecase,
	}
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Email           string   `json:"email" binding:"required,email"`
	Username        string   `json:"username" binding:"required,min=3,max=30"`
	Password        string   `json:"password" binding:"required,min=8"`
	FullName        string   `json:"full_name" binding:"required,min=1"`
	PhoneNumber     string   `json:"phone_number,omitempty"`
	Location        string   `json:"location,omitempty"`
	Skills          []string `json:"skills,omitempty"`
	ExperienceYears int      `json:"experience_years,omitempty"`
	Bio             string   `json:"bio,omitempty"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// VerifyEmailRequest represents the email verification request
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

// UpdateProfileRequest represents the profile update request
type UpdateProfileRequest struct {
	FullName        *string   `json:"full_name,omitempty"`
	PhoneNumber     *string   `json:"phone_number,omitempty"`
	Location        *string   `json:"location,omitempty"`
	Skills          *[]string `json:"skills,omitempty"`
	ExperienceYears *int      `json:"experience_years,omitempty"`
	Bio             *string   `json:"bio,omitempty"`
	ProfilePicture  *string   `json:"profile_picture,omitempty"`
}

// ChangePasswordRequest represents the password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// RequestPasswordResetRequest represents the password reset request
type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the password reset request
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UpdateUserRoleRequest represents the admin role update request
type UpdateUserRoleRequest struct {
	Role domain.Role `json:"role" binding:"required,oneof=user admin"`
}

// ResendOTPRequest represents the request to resend OTP
type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// @Summary Register a new user
// @Description Register a new user account with email verification
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Router /api/v1/auth/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Email:           req.Email,
		Username:        req.Username,
		Password:        req.Password,
		FullName:        req.FullName,
		PhoneNumber:     req.PhoneNumber,
		Location:        req.Location,
		Skills:          req.Skills,
		ExperienceYears: req.ExperienceYears,
		Bio:             req.Bio,
	}

	if err := c.userUsecase.Register(ctx, user); err != nil {
		if err == domain.ErrEmailTaken {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		if err == domain.ErrUsernameTaken {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please check your email for verification code.",
		"email":   req.Email,
	})
}

// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful with tokens"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Router /api/v1/auth/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := c.userUsecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		case domain.ErrUserNotVerified:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Please verify your email before logging in"})
		case domain.ErrUserDeactivated:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Your account has been deactivated"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		}
		return
	}

	// Set refresh token as HTTP-only cookie
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(
		"refresh_token",
		tokens.RefreshToken,
		7*24*60*60, // 7 days
		"/",
		"",
		ctx.GetHeader("X-Forwarded-Proto") == "https", // secure in production
		true, // httpOnly
	)

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"access_token": tokens.AccessToken,
	})
}

// @Summary Verify email address
// @Description Verify user email with OTP code
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body VerifyEmailRequest true "Verification details"
// @Success 200 {object} map[string]interface{} "Email verified successfully"
// @Failure 400 {object} map[string]interface{} "Bad request or invalid OTP"
// @Router /api/v1/auth/verify-email [post]
func (c *UserController) VerifyEmail(ctx *gin.Context) {
	var req VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := domain.VerifyEmailInput{
		Email: req.Email,
		OTP:   req.OTP,
	}

	if err := c.userUsecase.VerifyEmail(ctx, input); err != nil {
		if err == domain.ErrInvalidOTP {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired OTP"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Verification failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully. You can now login.",
	})
}

// @Summary Get user profile
// @Description Get current user's profile information
// @Tags User Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.User "User profile"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /api/v1/users/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	user, err := c.userUsecase.GetProfile(ctx, userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// @Summary Update user profile
// @Description Update current user's profile information
// @Tags User Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Profile updates"
// @Success 200 {object} domain.User "Updated user profile"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/users/profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profileUpdates := domain.UserUpdateInput{
		FullName:        req.FullName,
		PhoneNumber:     req.PhoneNumber,
		Location:        req.Location,
		Skills:          req.Skills,
		ExperienceYears: req.ExperienceYears,
		Bio:             req.Bio,
		ProfilePicture:  req.ProfilePicture,
	}

	updatedUser, err := c.userUsecase.UpdateProfile(ctx, userID, profileUpdates)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    updatedUser,
	})
}

// @Summary Request password reset
// @Description Request a password reset link to be sent to the user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RequestPasswordResetRequest true "Email address"
// @Success 200 {object} map[string]interface{} "Password reset requested"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /api/v1/auth/request-password-reset [post]
func (c *UserController) RequestPasswordReset(ctx *gin.Context) {
    var req RequestPasswordResetRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, err := c.userUsecase.RequestPasswordReset(ctx, req.Email)
    if err != nil {
        if err == domain.ErrUserNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to request password reset"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Password reset link has been sent to your email",
    })
}

// @Summary Reset password
// @Description Reset the user's password using the token from the password reset email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Password reset details"
// @Success 200 {object} map[string]interface{} "Password reset successful"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /api/v1/auth/reset-password [post]
func (c *UserController) ResetPassword(ctx *gin.Context) {
    var req ResetPasswordRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    input := domain.ResetPasswordInput{
        Token:       req.Token,
        NewPassword: req.NewPassword,
    }

    err := c.userUsecase.ResetPassword(ctx, input)
    if err != nil {
        if err == domain.ErrInvalidToken {
            ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Password reset successful. You can now login with your new password.",
    })
}

// @Summary Change password
// @Description Change the user's password while logged in
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Password change details"
// @Success 200 {object} map[string]interface{} "Password changed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/auth/change-password [post]
func (c *UserController) ChangePassword(ctx *gin.Context) {
	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	if err := c.userUsecase.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword); err != nil {
		if err == domain.ErrInvalidCredentials {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// @Summary Delete account
// @Description Delete the user's account permanently
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Account deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/auth/delete-account [delete]
func (c *UserController) DeleteAccount(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	if err := c.userUsecase.DeleteAccount(ctx, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}

// @Summary Get users
// @Description Get a list of users (admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.User "List of users"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/admin/users [get]
func (c *UserController) GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Get users endpoint"})
}

// @Summary Update user role
// @Description Update the role of a user (admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateUserRoleRequest true "Role update details"
// @Success 200 {object} map[string]interface{} "User role updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/admin/users/role [put]
func (c *UserController) UpdateUserRole(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Update user role endpoint"})
}

// @Summary Toggle user status
// @Description Activate or deactivate a user account (admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User status toggled successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/admin/users/status/{user_id} [post]
func (c *UserController) ToggleUserStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Toggle user status endpoint"})
}

// @Summary Delete user
// @Description Delete a user account (admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/admin/users/{user_id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Delete user endpoint"})
}

// @Summary Resend OTP
// @Description Resend the OTP verification code to the user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ResendOTPRequest true "Email address"
// @Success 200 {object} map[string]interface{} "OTP resent successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /api/v1/auth/resend-otp [post]
func (c *UserController) ResendOTP(ctx *gin.Context) {
    var req ResendOTPRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    err := c.userUsecase.ResendOTP(ctx, req.Email)
    if err != nil {
        if err == domain.ErrUserNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resend OTP"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "Verification code resent to your email"})
}
