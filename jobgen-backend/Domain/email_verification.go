package domain

import (
	"context"
	"time"
)

type EmailVerification struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Email     string    `json:"email" bson:"email"`
	OTP       string    `json:"otp" bson:"otp"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Used      bool      `json:"used" bson:"used"`
}

type IEmailVerificationRepository interface {
	Store(ctx context.Context, verification *EmailVerification) error
	GetByEmail(ctx context.Context, email string) (*EmailVerification, error)
	GetByOTP(ctx context.Context, otp, email string) (*EmailVerification, error)
	MarkUsed(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}

type VerifyEmailInput struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}
