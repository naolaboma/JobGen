package usecases

import (
	"context"
	"fmt"
	domain "jobgen-backend/Domain"
	"math/rand"
	"strconv"
	"strings"
	"time"
	// Still used for password hashing, but not for reset token hashing
)

type userUsecase struct {
	userRepo              domain.IUserRepository
	emailVerificationRepo domain.IEmailVerificationRepository
	refreshTokenRepo      domain.IRefreshTokenRepository
	passwordResetRepo     domain.IPasswordResetRepository // This repo might become redundant or its methods will change
	jwtService            domain.IJWTService
	passwordService       domain.IPasswordService
	emailService          domain.IEmailService
	contextTimeout        time.Duration
}

func NewUserUsecase(
	userRepo domain.IUserRepository,
	emailVerificationRepo domain.IEmailVerificationRepository,
	refreshTokenRepo domain.IRefreshTokenRepository,
	passwordResetRepo domain.IPasswordResetRepository, // Keep for now, but its usage will change or be removed
	jwtService domain.IJWTService,
	passwordService domain.IPasswordService,
	emailService domain.IEmailService,
	timeout time.Duration,
) domain.IUserUsecase {
	return &userUsecase{
		userRepo:              userRepo,
		emailVerificationRepo: emailVerificationRepo,
		refreshTokenRepo:      refreshTokenRepo,
		passwordResetRepo:     passwordResetRepo, // Consider removing this if no longer needed
		jwtService:            jwtService,
		passwordService:       passwordService,
		emailService:          emailService,
		contextTimeout:        timeout,
	}
}

func (u *userUsecase) Register(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if err := u.validateUserInput(user); err != nil {
		return err
	}

	if err := u.passwordService.ValidateStrength(user.Password); err != nil {
		return err
	}

	hashedPassword, err := u.passwordService.Hash(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	user.Role = domain.RoleUser
	user.IsVerified = false
	user.IsActive = true
	if user.Skills == nil {
		user.Skills = []string{}
	}
	if user.ExperienceYears < 0 {
		user.ExperienceYears = 0
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return err
	}

	otp := u.generateOTP()
	verification := &domain.EmailVerification{
		Email:     user.Email,
		OTP:       otp,
		Purpose:   domain.PurposeEmailVerification, // Set purpose
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
		Used:      false,
	}

	if err := u.emailVerificationRepo.Store(ctx, verification); err != nil {
		u.userRepo.Delete(ctx, user.ID)
		return fmt.Errorf("failed to store email verification: %w", err)
	}

	if err := u.emailService.SendWelcomeEmail(ctx, user, otp); err != nil {
		fmt.Printf("Failed to send welcome email: %v\n", err)
	}

	return nil
}

func (u *userUsecase) validateUserInput(user *domain.User) error {
	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(user.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(user.Password) == "" {
		return fmt.Errorf("password is required")
	}
	if strings.TrimSpace(user.FullName) == "" {
		return fmt.Errorf("full name is required")
	}

	if len(user.Username) < 3 || len(user.Username) > 30 {
		return fmt.Errorf("username must be between 3 and 30 characters")
	}

	if len(user.FullName) < 1 || len(user.FullName) > 100 {
		return fmt.Errorf("full name must be between 1 and 100 characters")
	}

	return nil
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (*domain.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, domain.ErrUserDeactivated
	}

	if !user.IsVerified {
		return nil, domain.ErrUserNotVerified
	}

	if err := u.passwordService.Compare(user.Password, password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := u.jwtService.CreateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, refreshPayload, err := u.jwtService.CreateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	refreshTokenEntity := &domain.RefreshToken{
		TokenID:   refreshPayload.TokenID,
		Token:     refreshToken,
		UserID:    user.ID,
		CreatedAt: refreshPayload.IssuedAt,
		UpdatedAt: refreshPayload.IssuedAt,
		ExpiresAt: refreshPayload.ExpiresAt,
	}

	if err := u.refreshTokenRepo.StoreToken(ctx, refreshTokenEntity); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	if err := u.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
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

	if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.OTP) == "" {
		return domain.ErrInvalidInput
	}

	verification, err := u.emailVerificationRepo.GetByOTPAndEmailAndPurpose(ctx, input.OTP, input.Email, domain.PurposeEmailVerification)

	fmt.Println("Get OTp Response", verification, err)
	if err != nil {
		return fmt.Errorf("failed to get verification: %w", err)
	}
	if verification == nil {
		return domain.ErrInvalidOTP
	}

	if verification.Used {
		return domain.ErrOTPUsed
	}

	if time.Now().After(verification.ExpiresAt) {
		return domain.ErrOTPExpired
	}

	user, err := u.userRepo.GetByEmail(ctx, input.Email)
	fmt.Println("Get user by Email", user)
	if err != nil {
		fmt.Println("Get by Email", err)
		return domain.ErrUserNotFound
	}

	if user.IsVerified {
		if err := u.emailVerificationRepo.MarkUsed(ctx, verification.ID); err != nil {
			fmt.Printf("Warning: failed to mark verification as used: %v\n", err)
		}
		return nil
	}

	if err := u.emailVerificationRepo.MarkUsed(ctx, verification.ID); err != nil {
		return fmt.Errorf("failed to mark verification as used: %w", err)
	}

	user.IsVerified = true
	if err := u.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	return nil
}

func (u *userUsecase) UpdateProfile(ctx context.Context, userID string, updates domain.UserUpdateInput) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if updates.FullName != nil && strings.TrimSpace(*updates.FullName) != "" {
		user.FullName = strings.TrimSpace(*updates.FullName)
	}
	if updates.PhoneNumber != nil {
		user.PhoneNumber = strings.TrimSpace(*updates.PhoneNumber)
	}
	if updates.Location != nil {
		user.Location = strings.TrimSpace(*updates.Location)
	}
	if updates.Skills != nil {
		user.Skills = *updates.Skills
	}
	if updates.ExperienceYears != nil && *updates.ExperienceYears >= 0 {
		user.ExperienceYears = *updates.ExperienceYears
	}
	if updates.Bio != nil {
		user.Bio = strings.TrimSpace(*updates.Bio)
	}
	if updates.ProfilePicture != nil {
		user.ProfilePicture = strings.TrimSpace(*updates.ProfilePicture)
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	user.Password = ""
	
	return user, nil
}

func (u *userUsecase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := u.passwordService.Compare(user.Password, oldPassword); err != nil {
		return domain.ErrInvalidCredentials
	}

	if err := u.passwordService.ValidateStrength(newPassword); err != nil {
		return err
	}

	hashedPassword, err := u.passwordService.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	if err := u.userRepo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if err := u.refreshTokenRepo.DeleteAllTokensForUser(ctx, userID); err != nil {
		fmt.Printf("Failed to revoke refresh tokens: %v\n", err)
	}

	return nil
}

func (u *userUsecase) RequestPasswordResetOTP(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// For security, don't reveal if email exists or not, but proceed to send a dummy email
		// or at least log the attempt if the email doesn't exist to prevent enumeration attacks.
		// For this implementation, we will still return ErrUserNotFound to the caller,
		// but a real-world system might just return success and not send an email.
		return domain.ErrUserNotFound
	}

	otp := u.generateOTP()
	passwordResetVerification := &domain.EmailVerification{
		Email:     email,
		OTP:       otp,
		Purpose:   domain.PurposePasswordReset, // Set purpose for password reset
		ExpiresAt: time.Now().Add(15 * time.Minute), // OTPs usually have shorter life
		CreatedAt: time.Now(),
		Used:      false,
	}

	if err := u.emailVerificationRepo.Store(ctx, passwordResetVerification); err != nil {
		return fmt.Errorf("failed to store password reset OTP: %w", err)
	}

	// Send password reset email with OTP
	if err := u.emailService.SendPasswordResetEmail(ctx, user, otp); err != nil { // Modified to send OTP
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}

func (u *userUsecase) ResetPassword(ctx context.Context, input domain.ResetPasswordInput) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.OTP) == "" || strings.TrimSpace(input.NewPassword) == "" {
		return domain.ErrInvalidInput
	}

	if err := u.passwordService.ValidateStrength(input.NewPassword); err != nil {
		return err
	}

	// Get password reset verification by OTP and email for password reset purpose
	verification, err := u.emailVerificationRepo.GetByOTPAndEmailAndPurpose(ctx, input.OTP, input.Email, domain.PurposePasswordReset)
	if err != nil {
		return fmt.Errorf("failed to get password reset verification: %w", err)
	}
	if verification == nil {
		return domain.ErrInvalidOTP // Use a more specific error if needed
	}

	if verification.Used {
		return domain.ErrOTPUsed
	}

	if time.Now().After(verification.ExpiresAt) {
		return domain.ErrOTPExpired
	}

	// Mark OTP as used FIRST
	if err := u.emailVerificationRepo.MarkUsed(ctx, verification.ID); err != nil {
		return fmt.Errorf("failed to mark password reset OTP as used: %w", err)
	}

	// Get user by ID from verification record
	user, err := u.userRepo.GetByEmail(ctx, input.Email) // We have email in verification, so retrieve user by email
	if err != nil {
		return domain.ErrUserNotFound // Should not happen if verification exists
	}

	hashedPassword, err := u.passwordService.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	if err := u.userRepo.UpdatePassword(ctx, user.ID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if err := u.refreshTokenRepo.DeleteAllTokensForUser(ctx, user.ID); err != nil {
		fmt.Printf("Failed to revoke refresh tokens: %v\n", err)
	}

	return nil
}

func (u *userUsecase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrInvalidInput
	}

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}

func (u *userUsecase) DeleteAccount(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if err := u.refreshTokenRepo.DeleteAllTokensForUser(ctx, userID); err != nil {
		fmt.Printf("Failed to delete user tokens: %v\n", err)
	}

	return u.userRepo.Delete(ctx, userID)
}

func (u *userUsecase) GetUsers(ctx context.Context, filter domain.UserFilter) (*domain.PaginatedUsersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

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

	if adminUserID == targetUserID {
		return domain.ErrForbidden
	}

	user, err := u.userRepo.GetByID(ctx, targetUserID)
	fmt.Println("User found that needed to be updated from Admin", user, targetUserID)
	if err != nil {
		return err
	}

	if err := u.userRepo.UpdateRole(ctx, targetUserID, role); err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	if err := u.emailService.SendRoleChangeNotification(ctx, user, role); err != nil {
		fmt.Printf("Failed to send role change notification: %v\n", err)
	}

	return nil
}

func (u *userUsecase) ToggleUserStatus(ctx context.Context, adminUserID, targetUserID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if adminUserID == targetUserID {
		return domain.ErrForbidden
	}

	user, err := u.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	if err := u.userRepo.ToggleActiveStatus(ctx, targetUserID); err != nil {
		return fmt.Errorf("failed to toggle user status: %w", err)
	}

	if user.IsActive {
		if err := u.emailService.SendAccountDeactivationEmail(ctx, user); err != nil {
			fmt.Printf("Failed to send deactivation notification: %v\n", err)
		}
		
		if err := u.refreshTokenRepo.DeleteAllTokensForUser(ctx, targetUserID); err != nil {
			fmt.Printf("Failed to revoke user tokens: %v\n", err)
		}
	}

	return nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, adminUserID, targetUserID string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if adminUserID == targetUserID {
		return domain.ErrForbidden
	}

	if err := u.refreshTokenRepo.DeleteAllTokensForUser(ctx, targetUserID); err != nil {
		fmt.Printf("Failed to delete user tokens: %v\n", err)
	}

	return u.userRepo.Delete(ctx, targetUserID)
}

func (u *userUsecase) generateOTP() string {
	return strconv.Itoa(rand.Intn(900000) + 100000)
}
func (u *userUsecase) ResendOTP(ctx context.Context, email string, purpose domain.OTPPurpose) error {
    ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()

    user, err := u.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return domain.ErrUserNotFound
    }

    // Only enforce verification check for email verification
    if purpose == domain.PurposeEmailVerification {
        if user.IsVerified {
            return domain.ErrAlreadyVerified
        }
    }

    otp := u.generateOTP()
    verification := &domain.EmailVerification{
        Email:     email,
        OTP:       otp,
        Purpose:   purpose,
        ExpiresAt: time.Now().Add(15 * time.Minute),
        CreatedAt: time.Now(),
        Used:      false,
    }

    if err := u.emailVerificationRepo.Store(ctx, verification); err != nil {
        return fmt.Errorf("failed to store email verification: %w", err)
    }

    // Decide which email to send
    switch purpose {
    case domain.PurposeEmailVerification:
        if err := u.emailService.SendWelcomeEmail(ctx, user, otp); err != nil {
            return fmt.Errorf("failed to send verification email: %w", err)
        }
    case domain.PurposePasswordReset:
        if err := u.emailService.SendPasswordResetEmail(ctx, user, otp); err != nil {
            return fmt.Errorf("failed to send password reset email: %w", err)
        }
    default:
        return fmt.Errorf("invalid OTP purpose: %s", purpose)
    }

    return nil
}
