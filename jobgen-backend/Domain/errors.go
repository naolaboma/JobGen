package domain

import "errors"

var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmailTaken        = errors.New("email is already registered")
	ErrUsernameTaken     = errors.New("username is already taken")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotVerified   = errors.New("user email is not verified")
	ErrUserDeactivated   = errors.New("user account is deactivated")
	
	// Authentication errors
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token has expired")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	
	// Password errors
	ErrWeakPassword      = errors.New("password is too weak")
	ErrInvalidOTP        = errors.New("invalid or expired OTP")
	ErrInvalidResetToken = errors.New("invalid or expired reset token")
	
	// System errors
	ErrInternal          = errors.New("internal server error")
	ErrServiceUnavailable = errors.New("service temporarily unavailable")
	
	// OTP errors
	ErrOTPExpired        = errors.New("OTP expired")
	ErrOTPUsed           = errors.New("OTP already used")
)
