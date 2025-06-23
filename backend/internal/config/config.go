// internal/config/config.go
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration settings.
type Config struct {
	DatabaseURL   string
	AppEnv        string
	GCSBucketName string
	// Add any other configuration variables here (e.g., API keys, service endpoints)
}

// LoadConfig loads configuration from environment variables or .env file.
func LoadConfig() (*Config, error) {
	// Load .env file in development. In production, env vars are set directly.
	// It's fine if .env doesn't exist (e.g., in production deployments).
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development" // Default to development if not set
	}

	gcsBucketName := os.Getenv("GCS_BUCKET_NAME")
	if gcsBucketName == "" {
		return nil, fmt.Errorf("GCS_BUCKET_NAME environment variable not set")
	}

	return &Config{
		DatabaseURL:   dbURL,
		AppEnv:        appEnv,
		GCSBucketName: gcsBucketName,
	}, nil
}
