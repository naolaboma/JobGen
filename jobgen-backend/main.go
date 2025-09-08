package main

import (
	"context"
	controllers "jobgen-backend/Delivery/Controllers"
	router "jobgen-backend/Delivery/Router"
	domain "jobgen-backend/Domain"
	infrastructure "jobgen-backend/Infrastructure"
	"jobgen-backend/Infrastructure/services"
	repositories "jobgen-backend/Repositories"
	usecases "jobgen-backend/Usecases"
	worker "jobgen-backend/Worker"
	_ "jobgen-backend/docs" // This line is important for swagger
	"log"
	"time"

	"github.com/go-redis/redis/v8"
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
	contactRepo := repositories.NewContactRepository(db)
	jobRepo := repositories.NewJobRepository(db)

	// Initialize job-related services
	jobAggregationService := services.NewJobAggregationService(jobRepo)
	jobMatchingService := services.NewJobMatchingService(jobRepo, userRepo)

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
	contactUsecase := usecases.NewContactUsecase(
		contactRepo,
		emailService,
		contextTimeout,
	)
	jobUsecase := usecases.NewJobUsecase(
		jobRepo,
		userRepo,
		jobAggregationService,
		jobMatchingService,
		contextTimeout,
	)

	// Initialize AI Service for chatbot
	aiService, err := infrastructure.NewAIService()
	if err != nil {
		log.Fatalf("Failed to initialize AI service: %v", err)
	}

	// Initialize Chat components
	chatRepo := repositories.NewChatRepository(db)
	chatUsecase := usecases.NewChatUsecase(chatRepo, aiService)
	chatController := controllers.NewChatController(chatUsecase)

	// Initialize controllers
	userController := controllers.NewUserController(userUsecase)
	authController := controllers.NewAuthController(authUsecase)
	contactController := controllers.NewContactController(contactUsecase)
	jobController := controllers.NewJobController(jobUsecase)

	// MinIO Setup
	minioURL := infrastructure.Env.FileStorageURL
	minioAccessKey := infrastructure.Env.AccessKey
	minioSecretKey := infrastructure.Env.SecretKey
	maxFileSize := infrastructure.Env.MaxAllowedFileSize
	maxUrlLife := infrastructure.Env.MaxFileUrlLife

	minioService, err := infrastructure.NewFileService(minioURL, minioAccessKey, minioSecretKey, maxFileSize, maxUrlLife)
	if err != nil {
		log.Fatal("MinIO setup error:", err)
	}

	// File Repository and Usecase
	fileRepo := repositories.NewFileRepository(db)
	fileUsecase := usecases.NewFileUsecase(fileRepo, minioService)
	fileController := controllers.NewFileController(fileUsecase)

	// --- Initialize Repositories ---
	cvRepo, err := repositories.NewCVRepository(db) // New CV Repo
	if err != nil {
		log.Fatalf("Could not create CV Repository: %v", err)
	}

	// --- Initialize Infrastructure & Services ---
	cvParserService := infrastructure.NewCVParserService() // New CV Parser

	// Initialize Queue service: try Redis, fall back to in-memory if not reachable
	var queueService infrastructure.QueueService
	{
		redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
		// Ping to check connectivity with a tiny timeout
		if err := redisClient.Ping(ctxWithTimeout(2 * time.Second)).Err(); err != nil {
			log.Printf("⚠️ Redis not available (%v). Falling back to in-memory queue.", err)
			queueService = infrastructure.NewInMemoryQueueService(200)
		} else {
			queueService = infrastructure.NewQueueService(redisClient, "cv_processing_queue")
		}
	}
	aiServiceClient := infrastructure.NewAIServiceClient() // Gemini AI Client

	// CV storage: prefer MinIO when configured; fallback to local disk for dev
	var cvStorage infrastructure.FileStorageService
	var cvDomainStorage domain.FileStorageService
	if infrastructure.Env.FileStorageURL != "" && infrastructure.Env.AccessKey != "" && infrastructure.Env.SecretKey != "" {
		// Use the same bucket as document uploads by default
		bucket := "documents"
		mstore, err := infrastructure.NewMinioCVFileStorageService(
			infrastructure.Env.FileStorageURL,
			infrastructure.Env.AccessKey,
			infrastructure.Env.SecretKey,
			bucket,
		)
		if err != nil {
			log.Printf("⚠️ MinIO CV storage init failed (%v). Falling back to local storage.", err)
			local := infrastructure.NewLocalCVFileStorageService("./data/cv")
			cvStorage = local
			cvDomainStorage = local
		} else {
			cvStorage = mstore
			cvDomainStorage = mstore
		}
	} else {
		local := infrastructure.NewLocalCVFileStorageService("./data/cv")
		cvStorage = local
		cvDomainStorage = local
	}

	// --- Initialize Usecases ---
	cvUsecase := usecases.NewCVUsecase(cvRepo, queueService, cvDomainStorage) // New CV Usecase

	// --- Initialize Controllers ---
	cvController := controllers.NewCVController(cvUsecase) // New CV Controller

	// --- Start Background Worker ---
	cvProcessor := worker.NewCVProcessor(queueService, cvRepo, cvParserService, cvStorage, aiServiceClient)
	go cvProcessor.Start() // Run the worker in a separate goroutine

	// Setup router (match parameter order defined in router.SetupRouter)
	router := router.SetupRouter(
		userController,
		authController,
		jobController,
		authMiddleware,
		fileController,
		cvController,
		contactController,
		chatController,
	)

	// Start server
	port := infrastructure.Env.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting JobGen API server on port %s", port)
	log.Printf("Environment: %s", infrastructure.Env.Environment)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/index.html", port)
	log.Printf("AI Chatbot endpoints available at: /api/v1/chat/*")

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

// ctxWithTimeout is a tiny helper to avoid repeating context.WithTimeout boilerplate.
func ctxWithTimeout(d time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), d)
	return ctx
}
