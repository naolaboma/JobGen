package router

import (
	controllers "jobgen-backend/Delivery/Controllers"
	infrastructure "jobgen-backend/Infrastructure"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	userController *controllers.UserController,
	authController *controllers.AuthController,
	jobController *controllers.JobController,
	authMiddleware *infrastructure.AuthMiddleware,
	fileController *controllers.FileController,
	cvController *controllers.CVController,
	contactController *controllers.ContactController,
	chatController *controllers.ChatController, // Add this parameter
) *gin.Engine {
	r := gin.New()

	// Configure trusted proxies per Gin security guidance
	// https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	if tp := os.Getenv("TRUSTED_PROXIES"); strings.TrimSpace(tp) == "" {
		if err := r.SetTrustedProxies(nil); err != nil {
			log.Printf("failed to set trusted proxies: %v", err)
		}
	} else {
		parts := strings.Split(tp, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		if err := r.SetTrustedProxies(parts); err != nil {
			log.Printf("failed to set trusted proxies: %v", err)
		}
	}

	// Prevent automatic trailing slash redirects (fix 307 for POST)
	r.RedirectTrailingSlash = false

	// CORS setup - allow all origins, methods, headers (development)
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // allow everything
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // must be false when AllowAllOrigins = true
		MaxAge:           12 * time.Hour,
	}))

	// Respond to preflight OPTIONS requests
	r.OPTIONS("/*cors", func(c *gin.Context) {
		c.AbortWithStatus(204)
	})

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		// Public routes
		api.POST("/contact", contactController.SubmitContactForm)

		auth := api.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
			auth.POST("/verify-email", userController.VerifyEmail)
			auth.POST("/forgot-password", userController.RequestPasswordReset)
			auth.POST("/reset-password", userController.ResetPassword)
			auth.POST("/refresh", authController.RefreshToken)
			auth.POST("/logout", authMiddleware.RequireAuth(), authController.Logout)
			auth.POST("/resend-otp", userController.ResendOTP)
			auth.POST("/change-password", authMiddleware.RequireAuth(), userController.ChangePassword)
		}

		users := api.Group("/users")
		users.Use(authMiddleware.RequireAuth())
		{
			users.GET("/profile", userController.GetProfile)
			users.PUT("/profile", userController.UpdateProfile)
			users.DELETE("/account", userController.DeleteAccount)
		}

		// Job routes
		jobs := api.Group("/jobs")
		{
			jobs.GET("", jobController.GetJobs)
			jobs.GET("/", jobController.GetJobs)
			jobs.GET("/:id", jobController.GetJobByID)
			jobs.GET("/trending", jobController.GetTrendingJobs)
			jobs.GET("/stats", jobController.GetJobStats)
			jobs.GET("/sources", jobController.GetJobSources)
			jobs.GET("/search-by-skills", jobController.SearchJobsBySkills)
			jobs.GET("/search", authMiddleware.OptionalAuth(), jobController.SearchJobs)

			authenticated := jobs.Group("/")
			authenticated.Use(authMiddleware.RequireAuth())
			{
				authenticated.GET("/matched", jobController.GetMatchedJobs)
			}
		}

		admin := api.Group("/admin")
		admin.Use(authMiddleware.RequireAdmin())
		{
			admin.GET("/users", userController.GetUsers)
			admin.PUT("/users/:user_id/role", userController.UpdateUserRole)
			admin.PUT("/users/:user_id/toggle-status", userController.ToggleUserStatus)
			admin.DELETE("/users/:user_id", userController.DeleteUser)

			jobAdmin := admin.Group("/jobs")
			{
				jobAdmin.POST("/aggregate", jobController.TriggerJobAggregation)
				jobAdmin.POST("/", jobController.CreateJob)
				jobAdmin.PUT("/:id", jobController.UpdateJob)
				jobAdmin.DELETE("/:id", jobController.DeleteJob)
			}
		}

		files := api.Group("/files")
		files.Use(authMiddleware.RequireAuth())
		{
			files.GET("/:id", fileController.DownloadFile)
			files.GET("/profile-picture/:id", fileController.GetProfilePicture)
			files.GET("/profile-picture/me", fileController.GetMyProfilePicture)
			files.POST("/upload/profile", fileController.UploadProfile)
			files.POST("/upload/document", fileController.UploadDocument)
			files.DELETE("/:id", fileController.DeleteFile)
		}

		// CV routes (updated to match Users/Files style)
		cv := api.Group("/cv")
		cv.Use(authMiddleware.RequireAuth())
		{
			cv.POST("/", cvController.StartParsingJobFromRef)
			cv.POST("/parse", cvController.StartParsingJobHandler)
			cv.GET("/parse/:jobId/status", cvController.GetParsingJobStatusHandler)
			cv.GET("/:id", cvController.GetParsingJobStatusHandler)
		}

		// Chat routes - moved inside the api group
		chatRoutes := api.Group("/chat")
		chatRoutes.Use(authMiddleware.RequireAuth())
		{
			chatRoutes.POST("/message", chatController.SendMessage)
			chatRoutes.GET("/sessions", chatController.GetUserSessions)
			chatRoutes.GET("/session/:session_id", chatController.GetSessionHistory)
			chatRoutes.DELETE("/session/:session_id", chatController.DeleteSession)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "Service is healthy",
			"data": gin.H{
				"status":    "ok",
				"timestamp": time.Now().UTC(),
				"version":   "1.0.0",
			},
		})
	})

	return r
}
