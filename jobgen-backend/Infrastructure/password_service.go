package infrastructure

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	domain "jobgen-backend/Domain"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct{}

func NewPasswordService() domain.IPasswordService {
	return &PasswordService{}
}

func (p *PasswordService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

func (p *PasswordService) Compare(hashed, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return domain.ErrInvalidCredentials
	}
	return nil
}

func (p *PasswordService) ValidateStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return fmt.Errorf("password must be less than 128 characters")
	}

	var (
		hasLower   = false
		hasUpper   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	requirements := []string{}
	if !hasLower {
		requirements = append(requirements, "at least one lowercase letter")
	}
	if !hasUpper {
		requirements = append(requirements, "at least one uppercase letter")
	}
	if !hasNumber {
		requirements = append(requirements, "at least one number")
	}
	if !hasSpecial {
		requirements = append(requirements, "at least one special character")
	}

	if len(requirements) > 0 {
		return fmt.Errorf("password must contain %s", strings.Join(requirements, ", "))
	}

	return nil
}

func (p *PasswordService) GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
