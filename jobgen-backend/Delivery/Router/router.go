package router

import (
	controllers "jobgen-backend/Delivery/Controllers"
	infrastructure "jobgen-backend/Infrastructure"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	userController *controllers.UserController,
	authController *controllers.AuthController,
	authMiddleware *infrastructure.AuthMiddleware,
	fileController *controllers.FileController,
) *gin.Engine {
	r := gin.Default()

	// CORS setup (production: restrict origins)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
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

		admin := api.Group("/admin")
		admin.Use(authMiddleware.RequireAdmin())
		{
			admin.GET("/users", userController.GetUsers)
			admin.PUT("/users/:user_id/role", userController.UpdateUserRole)
			admin.PUT("/users/:user_id/toggle-status", userController.ToggleUserStatus)
			admin.DELETE("/users/:user_id", userController.DeleteUser)
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
