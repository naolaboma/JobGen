package domain

import (
	"context"
	"time"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	Email           string    `json:"email" bson:"email"`
	Username        string    `json:"username" bson:"username"`
	Password        string    `json:"-" bson:"password"`
	FullName        string    `json:"full_name" bson:"full_name"`
	PhoneNumber     string    `json:"phone_number" bson:"phone_number"`
	Location        string    `json:"location" bson:"location"`
	Skills          []string  `json:"skills" bson:"skills"`
	ExperienceYears int       `json:"experience_years" bson:"experience_years"`
	Bio             string    `json:"bio" bson:"bio"`
	ProfilePicture  string    `json:"profile_picture" bson:"profile_picture"`
	Role            Role      `json:"role" bson:"role"`
	IsVerified      bool      `json:"is_verified" bson:"is_verified"`
	IsActive        bool      `json:"is_active" bson:"is_active"`
	LastLoginAt     *time.Time `json:"last_login_at" bson:"last_login_at"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

// Repository interfaces
type IUserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, userID, newPasswordHash string) error
	UpdateLastLogin(ctx context.Context, userID string) error
	Delete(ctx context.Context, userID string) error
	List(ctx context.Context, filter UserFilter) ([]User, int64, error)
	UpdateRole(ctx context.Context, userID string, role Role) error
	ToggleActiveStatus(ctx context.Context, userID string) error
}

// Use case interfaces
type IUserUsecase interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, email, password string) (*TokenResponse, error)
	VerifyEmail(ctx context.Context, input VerifyEmailInput) error
	UpdateProfile(ctx context.Context, userID string, updates UserUpdateInput) (*User, error)
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
	RequestPasswordReset(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, input ResetPasswordInput) error
	GetProfile(ctx context.Context, userID string) (*User, error)
	DeleteAccount(ctx context.Context, userID string) error
	// Admin operations
	GetUsers(ctx context.Context, filter UserFilter) (*PaginatedUsersResponse, error)
	UpdateUserRole(ctx context.Context, adminUserID, targetUserID string, role Role) error
	ToggleUserStatus(ctx context.Context, adminUserID, targetUserID string) error
	DeleteUser(ctx context.Context, adminUserID, targetUserID string) error
	ResendOTP(ctx context.Context, email string) error
}

// DTOs and filters
type UserFilter struct {
	Role      *Role  `json:"role,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"`
	Search    string `json:"search,omitempty"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

type UserUpdateInput struct {
	FullName        *string   `json:"full_name,omitempty"`
	PhoneNumber     *string   `json:"phone_number,omitempty"`
	Location        *string   `json:"location,omitempty"`
	Skills          *[]string `json:"skills,omitempty"`
	ExperienceYears *int      `json:"experience_years,omitempty"`
	Bio             *string   `json:"bio,omitempty"`
	ProfilePicture  *string   `json:"profile_picture,omitempty"`
}

type PaginatedUsersResponse struct {
	Users      []User `json:"users"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Total      int64  `json:"total"`
	TotalPages int    `json:"total_pages"`
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
}
