package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	controllers "jobgen-backend/Delivery/Controllers"
	router "jobgen-backend/Delivery/Router"
	domain "jobgen-backend/Domain"
	infrastructure "jobgen-backend/Infrastructure"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Test Suite
type APITestSuite struct {
	suite.Suite
	router         *gin.Engine
	userUsecase    *MockUserUsecase
	authUsecase    *MockAuthUsecase
	fileUsecase    *MockFileUsecase
	jwtService     *MockJWTService
	userController *controllers.UserController
	authController *controllers.AuthController
	fileController *controllers.FileController
	authMiddleware *infrastructure.AuthMiddleware
}

// Mock implementations
type MockFileUsecase struct {
	mock.Mock
}

type MockUserUsecase struct {
	mock.Mock
}

type MockAuthUsecase struct {
	mock.Mock
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockFileUsecase) Upload(ctx context.Context, file io.Reader, metaData *domain.File) (*domain.File, error) {
	args := m.Called(ctx, file, metaData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.File), args.Error(1)
}

func (m *MockFileUsecase) Download(ctx context.Context, fileID string, userID string) (string, error) {
	args := m.Called(ctx, fileID, userID)
	return args.String(0), args.Error(1)
}

func (m *MockFileUsecase) Delete(ctx context.Context, ID string, userID string) error {
	args := m.Called(ctx, ID, userID)
	return args.Error(0)
}

func (m *MockFileUsecase) Exists(ctx context.Context, ID string) (bool, error) {
	args := m.Called(ctx, ID)
	return args.Bool(0), args.Error(1)
}

func (m *MockFileUsecase) GetProfilePictureByUserID(ctx context.Context, ID string) (string, error) {
	args := m.Called(ctx, ID)
	return args.String(0), args.Error(1)
}

// Mock UserUsecase methods
func (m *MockUserUsecase) Register(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUsecase) Login(ctx context.Context, email, password string) (*domain.TokenResponse, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenResponse), args.Error(1)
}

func (m *MockUserUsecase) VerifyEmail(ctx context.Context, input domain.VerifyEmailInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockUserUsecase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUsecase) UpdateProfile(ctx context.Context, userID string, updates domain.UserUpdateInput) (*domain.User, error) {
	args := m.Called(ctx, userID, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUsecase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	args := m.Called(ctx, userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockUserUsecase) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) ResetPassword(ctx context.Context, input domain.ResetPasswordInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockUserUsecase) DeleteAccount(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserUsecase) GetUsers(ctx context.Context, filter domain.UserFilter) (*domain.PaginatedUsersResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PaginatedUsersResponse), args.Error(1)
}

func (m *MockUserUsecase) UpdateUserRole(ctx context.Context, adminUserID, targetUserID string, role domain.Role) error {
	args := m.Called(ctx, adminUserID, targetUserID, role)
	return args.Error(0)
}

func (m *MockUserUsecase) ToggleUserStatus(ctx context.Context, adminUserID, targetUserID string) error {
	args := m.Called(ctx, adminUserID, targetUserID)
	return args.Error(0)
}

func (m *MockUserUsecase) DeleteUser(ctx context.Context, adminUserID, targetUserID string) error {
	args := m.Called(ctx, adminUserID, targetUserID)
	return args.Error(0)
}

func (m *MockUserUsecase) ResendOTP(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

// Mock AuthUsecase methods
func (m *MockAuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenResponse), args.Error(1)
}

func (m *MockAuthUsecase) Logout(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Mock JWTService methods
func (m *MockJWTService) CreateAccessToken(user *domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) CreateRefreshToken(user *domain.User) (string, *domain.RefreshTokenPayload, error) {
	args := m.Called(user)
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*domain.RefreshTokenPayload), args.Error(2)
}

func (m *MockJWTService) ValidateAccessToken(tokenString string) (*domain.AccessTokenPayload, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AccessTokenPayload), args.Error(1)
}

func (m *MockJWTService) ValidateRefreshToken(tokenString string) (*domain.RefreshTokenPayload, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshTokenPayload), args.Error(1)
}

// Setup and teardown
func (suite *APITestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	// Initialize mocks
	suite.userUsecase = new(MockUserUsecase)
	suite.authUsecase = new(MockAuthUsecase)
	suite.jwtService = new(MockJWTService)
	suite.fileUsecase = new(MockFileUsecase)

	// Initialize controllers
	suite.userController = controllers.NewUserController(suite.userUsecase)
	suite.authController = controllers.NewAuthController(suite.authUsecase)
	suite.fileController = controllers.NewFileController(suite.fileUsecase)
	suite.authMiddleware = infrastructure.NewAuthMiddleware(suite.jwtService)

	// Setup router
suite.router = router.SetupRouter(suite.userController, suite.authController, suite.authMiddleware, nil, nil)
}

func (suite *APITestSuite) TearDownTest() {
	// Assert that all expectations were met
	suite.userUsecase.AssertExpectations(suite.T())
	suite.authUsecase.AssertExpectations(suite.T())
	suite.jwtService.AssertExpectations(suite.T())
	suite.fileUsecase.AssertExpectations(suite.T())
}

// Helper functions
func (suite *APITestSuite) makeRequest(method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	req, _ := http.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	return w
}

func (suite *APITestSuite) makeAuthenticatedRequest(method, path string, body interface{}, userID string, role domain.Role) *httptest.ResponseRecorder {
	// Mock valid token validation
	suite.jwtService.On("ValidateAccessToken", "valid-token").Return(&domain.AccessTokenPayload{
		UserID: userID,
		Email:  "test@example.com",
		Role:   role,
	}, nil).Maybe()

	headers := map[string]string{
		"Authorization": "Bearer valid-token",
	}

	return suite.makeRequest(method, path, body, headers)
}

func (suite *APITestSuite) createTestUser() *domain.User {
	return &domain.User{
		ID:              "test-user-id",
		Email:           "test@example.com",
		Username:        "testuser",
		FullName:        "Test User",
		Role:            domain.RoleUser,
		IsVerified:      true,
		IsActive:        true,
		Skills:          []string{},
		ExperienceYears: 0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func (suite *APITestSuite) createTestAdmin() *domain.User {
	admin := suite.createTestUser()
	admin.ID = "admin-user-id"
	admin.Email = "admin@example.com"
	admin.Username = "admin"
	admin.FullName = "Admin User"
	admin.Role = domain.RoleAdmin
	return admin
}

// Standard response structure for assertions
type TestResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   *struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details"`
	} `json:"error"`
}

func (suite *APITestSuite) parseResponse(w *httptest.ResponseRecorder) *TestResponse {
	var response TestResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Require().NoError(err)
	return &response
}

// Run the test suite
func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
