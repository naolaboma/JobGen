package controllers

import (
	"fmt"
	domain "jobgen-backend/Domain"
	"net/http"
	"strconv"

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
// @Success 201 {object} StandardResponse "User registered successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 409 {object} StandardResponse "User already exists"
// @Router /auth/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(ctx, err)
		return
	}

	// Validate required fields explicitly
	if req.FullName == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Full name is required", nil)
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
		switch err {
		case domain.ErrEmailTaken:
			ConflictResponse(ctx, "Email already registered")
		case domain.ErrUsernameTaken:
			ConflictResponse(ctx, "Username already taken")
		default:
			ErrorResponse(ctx, http.StatusBadRequest, "REGISTRATION_ERROR", err.Error(), nil)
		}
		return
	}

	SuccessResponse(ctx, http.StatusCreated, "User registered successfully. Please check your email for verification code.", gin.H{
		"email": req.Email,
	})
}

// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} StandardResponse "Login successful with tokens"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Invalid credentials"
// @Router /auth/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(ctx, err)
		return
	}

	tokens, err := c.userUsecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			UnauthorizedResponse(ctx, "Invalid email or password")
		case domain.ErrUserNotVerified:
			UnauthorizedResponse(ctx, "Please verify your email before logging in")
		case domain.ErrUserDeactivated:
			UnauthorizedResponse(ctx, "Your account has been deactivated")
		default:
			InternalErrorResponse(ctx, "Login failed")
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

	SuccessResponse(ctx, http.StatusOK, "Login successful", gin.H{
		"access_token": tokens.AccessToken,
	})
}

// @Summary Verify email address
// @Description Verify user email with OTP code
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body VerifyEmailRequest true "Verification details"
// @Success 200 {object} StandardResponse "Email verified successfully"
// @Failure 400 {object} StandardResponse "Bad request or invalid OTP"
// @Router /auth/verify-email [post]
func (c *UserController) VerifyEmail(ctx *gin.Context) {
	var req VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(ctx, err)
		return
	}

	input := domain.VerifyEmailInput{
		Email: req.Email,
		OTP:   req.OTP,
	}

	if err := c.userUsecase.VerifyEmail(ctx, input); err != nil {
		if err == domain.ErrInvalidOTP {
			ErrorResponse(ctx, http.StatusBadRequest, "INVALID_OTP", "Invalid or expired OTP", nil)
			return
		}
		fmt.Println("error in verify Email", err)
		InternalErrorResponse(ctx, "Verification failed")
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Email verified successfully. You can now login.", nil)
}

// @Summary Get user profile
// @Description Get current user's profile information
// @Tags User Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StandardResponse "User profile"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "User not found"
// @Router /users/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		UnauthorizedResponse(ctx, "User ID not found in token")
		return
	}

	user, err := c.userUsecase.GetProfile(ctx, userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			NotFoundResponse(ctx, "User not found")
			return
		}
		InternalErrorResponse(ctx, "Failed to get profile")
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Profile retrieved successfully", gin.H{
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
// @Success 200 {object} StandardResponse "Updated user profile"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Router /users/profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		UnauthorizedResponse(ctx, "User ID not found in token")
		return
	}

	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(ctx, err)
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
		InternalErrorResponse(ctx, "Failed to update profile")
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Profile updated successfully", gin.H{
		"user": updatedUser,
	})
}

// @Summary Request password reset
// @Description Request a password reset link to be sent to the user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RequestPasswordResetRequest true "Email address"
// @Success 200 {object} StandardResponse "Password reset requested"
// @Failure 400 {object} StandardResponse "Bad request"
// @Router /auth/forgot-password [post]
func (c *UserController) RequestPasswordReset(ctx *gin.Context) {
	var req RequestPasswordResetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(ctx, err)
		return
	}

	_, err := c.userUsecase.RequestPasswordReset(ctx, req.Email)
	if err != nil {
		// For security, don't reveal if user exists or not
		SuccessResponse(ctx, http.StatusOK, "If an account with that email exists, a password reset link has been sent", nil)
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Password reset link has been sent to your email", nil)
}

// @Summary Reset password
// @Description Reset the user's password using the token from the password reset email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Password reset details"
// @Success 200 {object} StandardResponse "Password reset successful"
// @Failure 400 {object} StandardResponse "Bad request"
// @Router /auth/reset-password [post]
func (c *UserController) ResetPassword(ctx *gin.Context) {
	var req ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(ctx, err)
		return
	}

	input := domain.ResetPasswordInput{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	}

	err := c.userUsecase.ResetPassword(ctx, input)
	if err != nil {
		switch err {
		case domain.ErrInvalidToken, domain.ErrInvalidResetToken:
			ErrorResponse(ctx, http.StatusBadRequest, "INVALID_TOKEN", "Invalid or expired reset token", nil)
		default:
			InternalErrorResponse(ctx, "Failed to reset password")
		}
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Password reset successful. You can now login with your new password.", nil)
}

// @Summary Change password
// @Description Change the user's password while logged in
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Password change details"
// @Success 200 {object} StandardResponse "Password changed successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Router /auth/change-password [post]
func (c *UserController) ChangePassword(ctx *gin.Context) {
	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(ctx, err)
		return
	}

	userID := ctx.GetString("user_id")
	if userID == "" {
		UnauthorizedResponse(ctx, "User ID not found in token")
		return
	}

	if err := c.userUsecase.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword); err != nil {
		if err == domain.ErrInvalidCredentials {
			UnauthorizedResponse(ctx, "Invalid old password")
			return
		}
		InternalErrorResponse(ctx, "Failed to change password")
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Password changed successfully", nil)
}

// @Summary Delete account
// @Description Delete the user's account permanently
// @Tags User Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StandardResponse "Account deleted successfully"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Router /users/account [delete]
func (c *UserController) DeleteAccount(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		UnauthorizedResponse(ctx, "User ID not found in token")
		return
	}

	if err := c.userUsecase.DeleteAccount(ctx, userID); err != nil {
		InternalErrorResponse(ctx, "Failed to delete account")
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Account deleted successfully", nil)
}

// @Summary Get users
// @Description Get a list of users (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param role query string false "Filter by role" Enums(user, admin)
// @Param active query bool false "Filter by active status"
// @Param search query string false "Search in email, username, or full name"
// @Param sort_by query string false "Sort field" default(created_at)
// @Param sort_order query string false "Sort order" Enums(asc, desc) default(desc)
// @Success 200 {object} StandardResponse "List of users"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 403 {object} StandardResponse "Forbidden"
// @Router /admin/users [get]
func (c *UserController) GetUsers(ctx *gin.Context) {
	if ctx.IsAborted() {
		return
	}
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	
	filter := domain.UserFilter{
		Page:      page,
		Limit:     limit,
		Search:    ctx.Query("search"),
		SortBy:    ctx.DefaultQuery("sort_by", "created_at"),
		SortOrder: ctx.DefaultQuery("sort_order", "desc"),
	}
	
	if role := ctx.Query("role"); role != "" {
		r := domain.Role(role)
		filter.Role = &r
	}
	
	if active := ctx.Query("active"); active != "" {
		isActive := active == "true"
		filter.IsActive = &isActive
	}

	result, err := c.userUsecase.GetUsers(ctx, filter)
	if err != nil {
		InternalErrorResponse(ctx, "Failed to get users")
		return
	}

	paginatedData := &PaginatedResponse{
		Items:      result.Users,
		Page:       result.Page,
		Limit:      result.Limit,
		Total:      result.Total,
		TotalPages: result.TotalPages,
		HasNext:    result.HasNext,
		HasPrev:    result.HasPrev,
	}

	PaginatedSuccessResponse(ctx, http.StatusOK, "Users retrieved successfully", paginatedData)
}

// @Summary Update user role
// @Description Update the role of a user (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Param request body UpdateUserRoleRequest true "Role update details"
// @Success 200 {object} StandardResponse "User role updated successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 403 {object} StandardResponse "Forbidden"
// @Router /admin/users/{user_id}/role [put]
func (c *UserController) UpdateUserRole(ctx *gin.Context) {
	adminUserID := ctx.GetString("user_id")
	targetUserID := ctx.Param("user_id")
	
	if adminUserID == "" {
		UnauthorizedResponse(ctx, "Admin user ID not found in token")
		return
	}
	
	if targetUserID == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "User ID is required", nil)
		return
	}

	var req UpdateUserRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(ctx, err)
		return
	}

	if err := c.userUsecase.UpdateUserRole(ctx, adminUserID, targetUserID, req.Role); err != nil {
		switch err {
		case domain.ErrForbidden:
			ForbiddenResponse(ctx, "Cannot change your own role")
		case domain.ErrUserNotFound:
			NotFoundResponse(ctx, "User not found")
		default:
			InternalErrorResponse(ctx, "Failed to update user role")
		}
		return
	}

	SuccessResponse(ctx, http.StatusOK, "User role updated successfully", nil)
}

// @Summary Toggle user status
// @Description Activate or deactivate a user account (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} StandardResponse "User status toggled successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 403 {object} StandardResponse "Forbidden"
// @Router /admin/users/{user_id}/toggle-status [put]
func (c *UserController) ToggleUserStatus(ctx *gin.Context) {
	adminUserID := ctx.GetString("user_id")
	targetUserID := ctx.Param("user_id")
	
	if adminUserID == "" {
		UnauthorizedResponse(ctx, "Admin user ID not found in token")
		return
	}
	
	if targetUserID == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "User ID is required", nil)
		return
	}

	if err := c.userUsecase.ToggleUserStatus(ctx, adminUserID, targetUserID); err != nil {
		switch err {
		case domain.ErrForbidden:
			ForbiddenResponse(ctx, "Cannot change your own status")
		case domain.ErrUserNotFound:
			NotFoundResponse(ctx, "User not found")
		default:
			InternalErrorResponse(ctx, "Failed to toggle user status")
		}
		return
	}

	SuccessResponse(ctx, http.StatusOK, "User status updated successfully", nil)
}

// @Summary Delete user
// @Description Delete a user account (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} StandardResponse "User deleted successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 403 {object} StandardResponse "Forbidden"
// @Router /admin/users/{user_id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	adminUserID := ctx.GetString("user_id")
	targetUserID := ctx.Param("user_id")
	
	if adminUserID == "" {
		UnauthorizedResponse(ctx, "Admin user ID not found in token")
		return
	}
	
	if targetUserID == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "User ID is required", nil)
		return
	}

	if err := c.userUsecase.DeleteUser(ctx, adminUserID, targetUserID); err != nil {
		switch err {
		case domain.ErrForbidden:
			ForbiddenResponse(ctx, "Cannot delete your own account")
		case domain.ErrUserNotFound:
			NotFoundResponse(ctx, "User not found")
		default:
			InternalErrorResponse(ctx, "Failed to delete user")
		}
		return
	}

	SuccessResponse(ctx, http.StatusOK, "User deleted successfully", nil)
}

// @Summary Resend OTP
// @Description Resend the OTP verification code to the user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ResendOTPRequest true "Email address"
// @Success 200 {object} StandardResponse "OTP resent successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 404 {object} StandardResponse "User not found"
// @Router /auth/resend-otp [post]
func (c *UserController) ResendOTP(ctx *gin.Context) {
    var req ResendOTPRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ValidationErrorResponse(ctx, err)
        return
    }
    
    err := c.userUsecase.ResendOTP(ctx, req.Email)
    if err != nil {
        if err == domain.ErrUserNotFound {
            NotFoundResponse(ctx, "User not found")
            return
        }
        if err == domain.ErrAlreadyVerified {
            ErrorResponse(ctx, http.StatusBadRequest, "ALREADY_VERIFIED", "Email is already verified", nil)
            return
        }
        InternalErrorResponse(ctx, "Failed to resend OTP")
        return
    }
    
    SuccessResponse(ctx, http.StatusOK, "Verification code resent to your email", nil)
}
