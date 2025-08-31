package infrastructure

import (
	"fmt"
	domain "jobgen-backend/Domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secretKey string
	accessTokenDuration time.Duration
	refreshTokenDuration time.Duration
}

func NewJWTService() domain.IJWTService {
	accessDuration, _ := time.ParseDuration(Env.AccessTokenDuration)
	refreshDuration, _ := time.ParseDuration(Env.RefreshTokenDuration)
	
	return &JWTService{
		secretKey: Env.JWTSecret,
		accessTokenDuration: accessDuration,
		refreshTokenDuration: refreshDuration,
	}
}

type JwtAccessClaims struct {
	UserID string      `json:"user_id"`
	Email  string      `json:"email"`
	Role   domain.Role `json:"role"`
	Type   string      `json:"type"`
	jwt.RegisteredClaims
}

type JwtRefreshClaims struct {
	TokenID string `json:"token_id"`
	UserID  string `json:"user_id"`
	Type    string `json:"type"`
	jwt.RegisteredClaims
}

func (j *JWTService) CreateAccessToken(user *domain.User) (string, error) {
	claims := JwtAccessClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenDuration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "jobgen-api",
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to create access token: %w", err)
	}

	return tokenString, nil
}

func (j *JWTService) CreateRefreshToken(user *domain.User) (string, *domain.RefreshTokenPayload, error) {
	tokenID := uuid.NewString()
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(j.refreshTokenDuration)

	claims := JwtRefreshClaims{
		TokenID: tokenID,
		UserID:  user.ID,
		Type:    "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
			Issuer:    "jobgen-api",
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	payload := &domain.RefreshTokenPayload{
		TokenID:   tokenID,
		UserID:    user.ID,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	return tokenString, payload, nil
}

func (j *JWTService) ValidateAccessToken(tokenString string) (*domain.AccessTokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtAccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*JwtAccessClaims)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	if claims.Type != "access" {
		return nil, domain.ErrInvalidToken
	}

	return &domain.AccessTokenPayload{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}

func (j *JWTService) ValidateRefreshToken(tokenString string) (*domain.RefreshTokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtRefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*JwtRefreshClaims)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	if claims.Type != "refresh" {
		return nil, domain.ErrInvalidToken
	}

	return &domain.RefreshTokenPayload{
		TokenID:   claims.TokenID,
		UserID:    claims.UserID,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}
