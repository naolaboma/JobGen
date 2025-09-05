package domain

import (
	"context"
	"time"
)

type PasswordResetToken struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id"`
	TokenHash string    `json:"token_hash" bson:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Used      bool      `json:"used" bson:"used"`
}

type IPasswordResetRepository interface {
	Store(ctx context.Context, token *PasswordResetToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*PasswordResetToken, error)
	MarkUsed(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}

type RequestPasswordResetInput struct {
	Email string `json:"email"`
}

// ResetPasswordInput can remain the same or be adapted
type ResetPasswordInput struct {
	Email       string `json:"email" binding:"required,email"` // Added email
	OTP         string `json:"otp" binding:"required,len=6"` // Changed from Token
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
