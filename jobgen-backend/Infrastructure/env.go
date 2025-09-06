package infrastructure

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	// Database
	MongoDBURI string
	DBName     string

	// JWT
	JWTSecret            string
	AccessTokenDuration  string
	RefreshTokenDuration string

	// Server
	Port        string
	Environment string

	// Email
	EmailFrom     string
	EmailHost     string
	EmailPort     string
	EmailUsername string
	EmailPassword string
	AdminEmail    string // New: Admin email for notifications

	// Frontend URL (for email links)
	FrontendURL string

	// File storage
	AccessKey          string
	SecretKey          string
	FileStorageURL     string
	MaxAllowedFileSize int64
	MaxFileUrlLife     int64 // maximum time before firl url expiring in seconds

	// AI (Gemini API)
	GeminiAPIKey string
	GeminiModel  string
}

var Env EnvConfig

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
	maxSizeStr := getEnv("MAX_ALLOWED_FILE_SIZE", "3000000")
	maxSize, err := strconv.ParseInt(maxSizeStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid MAX_ALLOWED_FILE_SIZE: %v", err)
	}
	maxFileUrlLifeStr := getEnv("MAX_FILE_URL_LIFE", "300")
	maxFileUrlLife, err := strconv.ParseInt(maxFileUrlLifeStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid MAX_ALLOWED_FILE_SIZE: %v", err)
	}
	Env = EnvConfig{
		MongoDBURI:           getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DBName:               getEnv("DB_NAME", "jobgen"),
		JWTSecret:            getEnv("JWT_SECRET", ""),
		AccessTokenDuration:  getEnv("ACCESS_TOKEN_DURATION", "24h"),
		RefreshTokenDuration: getEnv("REFRESH_TOKEN_DURATION", "168h"), // 7 days
		Port:                 getEnv("PORT", "8080"),
		Environment:          getEnv("ENVIRONMENT", "development"),
		EmailFrom:            getEnv("EMAIL_FROM", ""),
		EmailHost:            getEnv("EMAIL_HOST", ""),
		EmailPort:            getEnv("EMAIL_PORT", "587"),
		EmailUsername:        getEnv("EMAIL_USERNAME", ""),
		EmailPassword:        getEnv("EMAIL_PASSWORD", ""),
		AdminEmail:           getEnv("ADMIN_EMAIL", ""), // New: Get admin email from env
		FrontendURL:          getEnv("FRONTEND_URL", "http://localhost:3000"),
		AccessKey:            getEnv("STORAGE_ACCESS_KEY", ""),
		SecretKey:            getEnv("STORAGE_SECRET_KEY", ""),
		FileStorageURL:       getEnv("FILE_STORAGE_URL", ""),
		MaxFileUrlLife:       maxFileUrlLife,
		MaxAllowedFileSize:   maxSize,
		GeminiAPIKey:         getEnv("GEMINI_API_KEY", ""),
		GeminiModel:          getEnv("GEMINI_MODEL", "gemini-1.5-flash"),
	}

	// Validate required environment variables
	if Env.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	if Env.MongoDBURI == "" {
		log.Fatal("MONGODB_URI is required")
	}
	// AdminEmail is not strictly required but good to warn if not set for contact forms
	if Env.AdminEmail == "" {
		log.Println("Warning: ADMIN_EMAIL is not set. Contact form submissions will not be emailed to an administrator.")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
