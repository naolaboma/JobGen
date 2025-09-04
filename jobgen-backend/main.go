package main

import (
	controllers "jobgen-backend/Delivery/Controllers"
	router "jobgen-backend/Delivery/Router"
	infrastructure "jobgen-backend/Infrastructure"
	repositories "jobgen-backend/Repositories"
	usecases "jobgen-backend/Usecases"
	_ "jobgen-backend/docs" // This line is important for swagger
	"log"
	"time"
)

// @title JobGen API
// @version 1.0
// @description AI-Powered Remote Job Finder & CV Optimizer API
// @termsOfService http://swagger.io/terms/

// @contact.name JobGen Support
// @contact.url http://www.jobgen.io/support
// @contact.email support@jobgen.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	infrastructure.LoadEnv()

	// Initialize database
	mongoClient := repositories.NewMongoClient()
	db := repositories.GetDatabase(mongoClient)

	// Initialize infrastructure services
	jwtService := infrastructure.NewJWTService()
	passwordService := infrastructure.NewPasswordService()
	emailService := infrastructure.NewEmailService()
	authMiddleware := infrastructure.NewAuthMiddleware(jwtService)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)
	emailVerificationRepo := repositories.NewEmailVerificationRepository(db)
	passwordResetRepo := repositories.NewPasswordResetRepository(db)

	// Initialize use cases
	contextTimeout := 30 * time.Second
	userUsecase := usecases.NewUserUsecase(
		userRepo,
		emailVerificationRepo,
		refreshTokenRepo,
		passwordResetRepo,
		jwtService,
		passwordService,
		emailService,
		contextTimeout,
	)
	authUsecase := usecases.NewAuthUsecase(
		jwtService,
		userRepo,
		refreshTokenRepo,
		contextTimeout,
	)

	// Initialize controllers
	userController := controllers.NewUserController(userUsecase)
	authController := controllers.NewAuthController(authUsecase)

	// --- MinIO Setup ---
	minioURL := infrastructure.Env.FileStorageURL
	minioAccessKey := infrastructure.Env.AccessKey
	minioSecretKey := infrastructure.Env.SecretKey
	maxFileSize := infrastructure.Env.MaxAllowedFileSize
	maxUrlLife := infrastructure.Env.MaxFileUrlLife

	minioService, err := infrastructure.NewFileService(minioURL, minioAccessKey, minioSecretKey, maxFileSize, maxUrlLife)
	if err != nil {
		log.Fatal("MinIO setup error:", err)
	}

	// --- Repository ---
	fileRepo := repositories.NewFileRepository(db)
	fileUsecase := usecases.NewFileUsecase(fileRepo, minioService)

	// --- Controller ---
	fileController := controllers.NewFileController(fileUsecase)
	// ----------------------------

	// Setup router
	router := router.SetupRouter(userController, authController, authMiddleware, fileController)

	// Start server
	port := infrastructure.Env.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting JobGen API server on port %s", port)
	log.Printf("Environment: %s", infrastructure.Env.Environment)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/index.html", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Additional model definitions for Swagger
type StandardResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string      `json:"code" example:"VALIDATION_ERROR"`
	Message string      `json:"message" example:"Invalid input provided"`
	Details interface{} `json:"details,omitempty"`
}
