package domain

import (
	"context"
	"time"
)

type OTPPurpose string

const (
	PurposeEmailVerification OTPPurpose = "email_verification"
	PurposePasswordReset     OTPPurpose = "password_reset"
)

type EmailVerification struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Email     string    `json:"email" bson:"email"`
	OTP       string    `json:"otp" bson:"otp"`
	Purpose   OTPPurpose `json:"purpose" bson:"purpose"` // Added field
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Used      bool      `json:"used" bson:"used"`
}

type IEmailVerificationRepository interface {
	Store(ctx context.Context, verification *EmailVerification) error
	GetByEmailAndPurpose(ctx context.Context, email string, purpose OTPPurpose) (*EmailVerification, error) // Modified
	GetByOTPAndEmailAndPurpose(ctx context.Context, otp, email string, purpose OTPPurpose) (*EmailVerification, error) // Modified
	MarkUsed(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
	DeleteExisting(ctx context.Context, email string, purpose OTPPurpose) error // New method
}

type VerifyEmailInput struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}
