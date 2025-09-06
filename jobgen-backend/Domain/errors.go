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
	ErrAlreadyVerified    = errors.New("email already verified")
	
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
	
	// Validation errors
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrInvalidUsername   = errors.New("invalid username format")
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters")
	ErrFullNameRequired  = errors.New("full name is required")
	
	// File related errors
	ErrFileTooBig 		 = errors.New("file size exceeds the allowed limits")
	ErrInvalidFileFormat = errors.New("file size exceeds the allowed limits")
	ErrUnknownFileType   = errors.New("file size exceeds the allowed limits")
	ErrFileNotFound      = errors.New("file size exceeds the allowed limits")

	// Job-related errors
	ErrJobNotFound    = errors.New("job not found")
	ErrJobExists      = errors.New("job already exists")
	ErrInvalidJobData = errors.New("invalid job data")
	ErrNotFound       = errors.New("resource not found")

	// Scraping errors
	ErrScrapingFailed     = errors.New("scraping failed")
	ErrRateLimitExceeded  = errors.New("rate limit exceeded")
	ErrSourceUnavailable  = errors.New("job source unavailable")

	// Matching errors
	ErrNoMatchingJobs     = errors.New("no matching jobs found")
	ErrInvalidPreferences = errors.New("invalid user preferences")
)
