package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	Auth0Domain   string
	Auth0Audience string
	AppEnv        string
	GCSBucketName string
	SentryDSN     string `env:"SENTRY_DSN"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	auth0Domain := os.Getenv("AUTH0_DOMAIN")
	if auth0Domain == "" {
		return nil, fmt.Errorf("AUTH0_DOMAIN environment variable not set")
	}

	auth0Audience := os.Getenv("AUTH0_AUDIENCE")
	if auth0Domain == "" {
		return nil, fmt.Errorf("AUTH0_AUDIENCE environment variable not set")
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	gcsBucketName := os.Getenv("GCS_BUCKET_NAME")
	if gcsBucketName == "" {
		return nil, fmt.Errorf("GCS_BUCKET_NAME environment variable not set")
	}

	SentryDSN := os.Getenv("SENTRY_DSN")
	if SentryDSN == "" {
		return nil, fmt.Errorf("SENTRY_DSN environment variable not set")
	}

	return &Config{
		DatabaseURL:   dbURL,
		Auth0Domain:   auth0Domain,
		Auth0Audience: auth0Audience,
		AppEnv:        appEnv,
		GCSBucketName: gcsBucketName,
		SentryDSN:     SentryDSN,
	}, nil
}
