package domain

import "context"

type IEmailService interface {
	SendWelcomeEmail(ctx context.Context, user *User, otp string) error
	SendPasswordResetEmail(ctx context.Context, user *User, resetToken string) error
	SendAccountDeactivationEmail(ctx context.Context, user *User) error
	SendRoleChangeNotification(ctx context.Context, user *User, newRole Role) error
}
