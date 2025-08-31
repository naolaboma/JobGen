package usecases

import (
	"context"
	"fmt"
	domain "jobgen-backend/Domain"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo              domain.IUserRepository
	emailVerificationRepo domain.IEmailVerificationRepository
	refreshTokenRepo      domain.IRefreshTokenRepository
	passwordResetRepo     domain.IPasswordResetRepository
	jwtService            domain.IJWTService
	passwordService       domain.IPasswordService
	emailService          domain.IEmailService
	contextTimeout        time.Duration
}

func NewUserUsecase(
	userRepo domain.IUserRepository,
	emailVerificationRepo domain.IEmailVerificationRepository,
	refreshTokenRepo domain.IRefreshTokenRepository,
	passwordResetRepo domain.IPasswordResetRepository,
	jwtService domain.IJWTService,
	passwordService domain.IPasswordService,
	emailService domain.IEmailService,
	timeout time.Duration,
) domain.IUserUsecase {
	return &userUsecase{
		userRepo:              userRepo,
		emailVerificationRepo: emailVerificationRepo,
		refreshTokenRepo:      refreshTokenRepo,
		passwordResetRepo:     passwordResetRepo,
		jwtService:            jwtService,
		passwordService:       passwordService,
		emailService:          emailService,
		contextTimeout:        timeout,
	}
}

func (u *userUsecase) Register(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Validate required fields
	if user.Email == "" || user.Username == "" || user.Password == "" || user.FullName == "" {
		return domain.ErrInvalidInput
	}

	// Validate password strength
	if err := u.passwordService.ValidateStrength(user.Password); err != nil {
		return err
	}

	// Hash password
	hashedPassword, err := u.passwordService.Hash(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Set user defaults
	user.Password = hashedPassword
	user.Role = domain.RoleUser
	user.IsVerified = false
	user.IsActive = true
	user.Skills = []string{}
	user.ExperienceYears = 0

	// Create user
	if err := u.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// Generate OTP
	otp := u.generateOTP()
	verification := &domain.EmailVerification{
		Email:     user.Email,
		OTP:       otp,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		Used:      false,
	}

	// Store verification
	if err := u.emailVerificationRepo.Store(ctx, verification); err != nil {
		return fmt.Errorf("failed to store email verification: %w", err)
	}

	// Send welcome email with OTP
	if err := u.emailService.SendWelcomeEmail(ctx, user, otp); err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	return nil
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (*domain.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, domain.ErrUserDeactivated
	}

	// Check if email is verified
	if !user.IsVerified {
		return nil, domain.ErrUserNotVerified
	}

	// Verify password
	if err := u.passwordService.Compare(user.Password, password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Create tokens
	accessToken, err := u.jwtService.CreateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, refreshPayload, err := u.jwtService.CreateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	// Store refresh token
	refreshTokenEntity := &domain.RefreshToken{
		TokenID:   refreshPayload.TokenID,
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: refreshPayload.ExpiresAt,
	}

	if err := u.refreshTokenRepo.StoreToken(ctx, refreshTokenEntity); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Update last login
	if err := u.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *userUsecase) VerifyEmail(ctx context.Context, input domain.VerifyEmailInput) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// fetch verification by otp & email
	verification, err := u.emailVerificationRepo.GetByOTP(ctx, input.OTP, input.Email)
	if err != nil {
		return err
	}
	if verification == nil {
		return domain.ErrInvalidOTP
	}

	// check if OTP was already used
	if verification.Used {
		return domain.ErrInvalidOTP
	}

	// check expiration: prefer ExpiresAt if present, otherwise CreatedAt + 15min
	now := time.Now()
	expiry := verification.ExpiresAt
	if expiry.IsZero() && !verification.CreatedAt.IsZero() {
		expiry = verification.CreatedAt.Add(15 * time.Minute)
	}
	if expiry.IsZero() || now.After(expiry) {
		return fmt.Errorf("verification code expired")
	}

	// load user
	user, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// if already verified, consume OTP and return
	if user.IsVerified {
		if err := u.emailVerificationRepo.MarkUsed(ctx, verification.ID); err != nil {
			return fmt.Errorf("failed to mark verification as used: %w", err)
		}
		return nil
	}

	// set verified and persist
	user.IsVerified = true
	if err := u.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	// mark verification as used. If this fails, attempt best-effort user verification
	if err := u.emailVerificationRepo.MarkUsed(ctx, verification.ID); err != nil {
		fmt.Printf("Warning: failed to mark verification as used: %v\n", err)
	}

	return nil
}

func (u *userUsecase) UpdateProfile(ctx context.Context, userID string, updates domain.UserUpdateInput) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Get current user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if updates.FullName != nil {
		user.FullName = *updates.FullName
	}
	if updates.PhoneNumber != nil {
		user.PhoneNumber = *updates.PhoneNumber
	}
	if updates.Location != nil {
		user.Location = *updates.Location
	}
	if updates.Skills != nil {
		user.Skills = *updates.Skills
	}
	if updates.ExperienceYears != nil {
		user.ExperienceYears = *updates.ExperienceYears
	}
	if updates.Bio != nil {
		user.Bio = *updates.Bio
	}
	if updates.ProfilePicture != nil {
		user.ProfilePicture = *updates.ProfilePicture
	}

	// Update user
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return user, nil
}

func (u *userUsecase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Get user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := u.passwordService.Compare(user.Password, oldPassword); err != nil {
		return domain.ErrInvalidCredentials
	}

	// Validate new password strength
	if err := u.passwordService.ValidateStrength(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := u.passwordService.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	if err := u.userRepo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all refresh tokens for security
	if err := u.refreshTokenRepo.DeleteAllTokensForUser(ctx, userID); err != nil {
		fmt.Printf("Failed to revoke refresh tokens: %v\n", err)
	}

	return nil
}

func (u *userUsecase) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Get user by email
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists or not
		return "", nil
	}

	// Generate reset token
	token, err := u.passwordService.GenerateRandomToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Hash the token for storage
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash reset token: %w", err)
	}

	// Store reset token
	resetToken := &domain.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: string(hashedToken),
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	if err := u.passwordResetRepo.Store(ctx, resetToken); err != nil {
		return "", fmt.Errorf("failed to store reset token: %w", err)
	}

	// Send reset email
	if err := u.emailService.SendPasswordResetEmail(ctx, user, token); err != nil {
		return "", fmt.Errorf("failed to send reset email: %w", err)
	}

	return token, nil
}

func (u *userUsecase) ResetPassword(ctx context.Context, input domain.ResetPasswordInput) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Validate new password
	if err := u.passwordService.ValidateStrength(input.NewPassword); err != nil {
		return err
	}

	// Get reset token
	resetToken, err := u.passwordResetRepo.GetByTokenHash(ctx, input.Token)
	if err != nil {
		return err
	}

	// Mark token as used
	if err := u.passwordResetRepo.MarkUsed(ctx, resetToken.ID); err != nil {
		return fmt.Errorf("failed to mark reset token as used: %w", err)
	}

	// Hash new password
	hashedPassword, err := u.passwordService.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	if err := u.userRepo.UpdatePassword(ctx, resetToken.UserID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all refresh tokens for security
	if err := u.refreshTokenRepo.DeleteAllTokensForUser(ctx, resetToken.UserID); err != nil {
		fmt.Printf("Failed to revoke refresh tokens: %v\n", err)
	}

	return nil
}

func (u *userUsecase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Clear sensitive information
	user.Password = ""

	return user, nil
}

func (u *userUsecase) DeleteAccount(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Delete all refresh tokens
	if err := u.refreshTokenRepo.DeleteAllTokensForUser(ctx, userID); err != nil {
		fmt.Printf("Failed to delete user tokens: %v\n", err)
	}

	// Delete user
	return u.userRepo.Delete(ctx, userID)
}

// Admin operations
func (u *userUsecase) GetUsers(ctx context.Context, filter domain.UserFilter) (*domain.PaginatedUsersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 10
	}

	users, total, err := u.userRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	// Clear passwords
	for i := range users {
		users[i].Password = ""
	}

	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))

	response := &domain.PaginatedUsersResponse{
		Users:      users,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    filter.Page < totalPages,
		HasPrev:    filter.Page > 1,
	}

	return response, nil
}

func (u *userUsecase) UpdateUserRole(ctx context.Context, adminUserID, targetUserID string, role domain.Role) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Prevent self-role changes
	if adminUserID == targetUserID {
		return domain.ErrForbidden
	}

	// Get target user for notification
	user, err := u.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	// Update role
	if err := u.userRepo.UpdateRole(ctx, targetUserID, role); err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	// Send notification email
	if err := u.emailService.SendRoleChangeNotification(ctx, user, role); err != nil {
		fmt.Printf("Failed to send role change notification: %v\n", err)
	}

	return nil
}

func (u *userUsecase) ToggleUserStatus(ctx context.Context, adminUserID, targetUserID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Prevent self-status changes
	if adminUserID == targetUserID {
		return domain.ErrForbidden
	}

	// Get user current status for notification
	user, err := u.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	// Toggle status
	if err := u.userRepo.ToggleActiveStatus(ctx, targetUserID); err != nil {
		return fmt.Errorf("failed to toggle user status: %w", err)
	}

	// Send notification if user was deactivated
	if user.IsActive {
		if err := u.emailService.SendAccountDeactivationEmail(ctx, user); err != nil {
			fmt.Printf("Failed to send deactivation notification: %v\n", err)
		}
	}

	// Revoke all refresh tokens if user is deactivated
	if user.IsActive {
		if err := u.refreshTokenRepo.DeleteAllTokensForUser(ctx, targetUserID); err != nil {
			fmt.Printf("Failed to revoke user tokens: %v\n", err)
		}
	}

	return nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, adminUserID, targetUserID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Prevent self-deletion
	if adminUserID == targetUserID {
		return domain.ErrForbidden
	}

	// Delete all refresh tokens first
	if err := u.refreshTokenRepo.DeleteAllTokensForUser(ctx, targetUserID); err != nil {
		fmt.Printf("Failed to delete user tokens: %v\n", err)
	}

	// Delete user
	return u.userRepo.Delete(ctx, targetUserID)
}

// Helper function to generate 6-digit OTP
func (u *userUsecase) generateOTP() string {
	return strconv.Itoa(rand.Intn(900000) + 100000)
}

func (u *userUsecase) ResendOTP(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return domain.ErrUserNotFound
	}
	if user.IsVerified {
		return nil // Already verified, no need to resend
	}

	otp := u.generateOTP()
	verification := &domain.EmailVerification{
		Email:     email,
		OTP:       otp,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		Used:      false,
	}
	if err := u.emailVerificationRepo.Store(ctx, verification); err != nil {
		return fmt.Errorf("failed to store email verification: %w", err)
	}
	if err := u.emailService.SendWelcomeEmail(ctx, user, otp); err != nil {
		return fmt.Errorf("failed to send welcome email: %w", err)
	}
	return nil
}
