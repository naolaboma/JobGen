package domain

import (
	"context"
	"time"
)

// JWT Service Interface
type IJWTService interface {
	CreateAccessToken(user *User) (string, error)
	CreateRefreshToken(user *User) (string, *RefreshTokenPayload, error)
	ValidateAccessToken(tokenString string) (*AccessTokenPayload, error)
	ValidateRefreshToken(tokenString string) (*RefreshTokenPayload, error)
}

// Password Service Interface
type IPasswordService interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) error
	ValidateStrength(password string) error
	GenerateRandomToken() (string, error)
}

// Auth Use case Interface
type IAuthUsecase interface {
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	Logout(ctx context.Context, userID string) error
}

// Token structures
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccessTokenPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   Role   `json:"role"`
}

type RefreshTokenPayload struct {
	TokenID   string    `json:"token_id"`
	UserID    string    `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Refresh Token Entity
type RefreshToken struct {
	ID        string     `json:"id" bson:"_id,omitempty"`
	TokenID   string     `json:"token_id" bson:"token_id"`
	Token     string     `json:"token" bson:"token"`
	UserID    string     `json:"user_id" bson:"user_id"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	ExpiresAt time.Time  `json:"expires_at" bson:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" bson:"revoked_at,omitempty"`
}

type IRefreshTokenRepository interface {
	StoreToken(ctx context.Context, token *RefreshToken) error
	FindByTokenID(ctx context.Context, tokenID string) (*RefreshToken, error)
	RevokeToken(ctx context.Context, tokenString string) error
	DeleteAllTokensForUser(ctx context.Context, userID string) error
	CleanupExpiredTokens(ctx context.Context) error
}
