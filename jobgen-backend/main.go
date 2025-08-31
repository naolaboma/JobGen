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

	// Setup router
	router := router.SetupRouter(userController, authController, authMiddleware)

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
