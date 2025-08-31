package infrastructure

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	// Database
	MongoDBURI string
	DBName     string
	
	// JWT
	JWTSecret             string
	AccessTokenDuration   string
	RefreshTokenDuration  string
	
	// Server
	Port string
	Environment string
	
	// Email
	EmailFrom     string
	EmailHost     string
	EmailPort     string
	EmailUsername string
	EmailPassword string
	
	// Frontend URL (for email links)
	FrontendURL string
}

var Env EnvConfig

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	Env = EnvConfig{
		MongoDBURI:           getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DBName:              getEnv("DB_NAME", "jobgen"),
		JWTSecret:           getEnv("JWT_SECRET", ""),
		AccessTokenDuration: getEnv("ACCESS_TOKEN_DURATION", "24h"),
		RefreshTokenDuration: getEnv("REFRESH_TOKEN_DURATION", "168h"), // 7 days
		Port:                getEnv("PORT", "8080"),
		Environment:         getEnv("ENVIRONMENT", "development"),
		EmailFrom:           getEnv("EMAIL_FROM", ""),
		EmailHost:           getEnv("EMAIL_HOST", ""),
		EmailPort:           getEnv("EMAIL_PORT", "587"),
		EmailUsername:       getEnv("EMAIL_USERNAME", ""),
		EmailPassword:       getEnv("EMAIL_PASSWORD", ""),
		FrontendURL:         getEnv("FRONTEND_URL", "http://localhost:3000"),
	}

	// Validate required environment variables
	if Env.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	
	if Env.MongoDBURI == "" {
		log.Fatal("MONGODB_URI is required")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
