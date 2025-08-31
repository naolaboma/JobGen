package usecases

import (
	"jobgen-backend/Domain"
	"context"
	"fmt"
	"time"
)

type authUsecase struct {
	jwtService       domain.IJWTService
	userRepo         domain.IUserRepository
	refreshTokenRepo domain.IRefreshTokenRepository
	contextTimeout   time.Duration
}

func NewAuthUsecase(
	jwtService domain.IJWTService,
	userRepo domain.IUserRepository,
	refreshTokenRepo domain.IRefreshTokenRepository,
	timeout time.Duration,
) domain.IAuthUsecase {
	return &authUsecase{
		jwtService:       jwtService,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		contextTimeout:   timeout,
	}
}

func (a *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	// Validate refresh token
	payload, err := a.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Get stored token
	storedToken, err := a.refreshTokenRepo.FindByTokenID(ctx, payload.TokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}

	if storedToken == nil {
		return nil, domain.ErrInvalidToken
	}

	// Check if token is revoked
	if storedToken.RevokedAt != nil {
		return nil, domain.ErrInvalidToken
	}

	// Check if token matches
	if storedToken.Token != refreshToken {
		return nil, domain.ErrInvalidToken
	}

	// Check if token is expired
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, domain.ErrTokenExpired
	}

	// Get user
	user, err := a.userRepo.GetByID(ctx, payload.UserID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive {
		return nil, domain.ErrUserDeactivated
	}

	// Create new tokens
	newAccessToken, err := a.jwtService.CreateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create new access token: %w", err)
	}

	newRefreshToken, newPayload, err := a.jwtService.CreateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create new refresh token: %w", err)
	}

	// Store new refresh token
	newTokenEntity := &domain.RefreshToken{
		TokenID:   newPayload.TokenID,
		Token:     newRefreshToken,
		UserID:    user.ID,
		ExpiresAt: newPayload.ExpiresAt,
	}

	if err := a.refreshTokenRepo.StoreToken(ctx, newTokenEntity); err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	// Revoke old refresh token
	if err := a.refreshTokenRepo.RevokeToken(ctx, refreshToken); err != nil {
		fmt.Printf("Failed to revoke old refresh token: %v\n", err)
	}

	return &domain.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (a *authUsecase) Logout(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	// Delete all refresh tokens for the user
	return a.refreshTokenRepo.DeleteAllTokensForUser(ctx, userID)
}
