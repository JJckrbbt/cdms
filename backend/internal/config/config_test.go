package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name        string
		setup       func(t *testing.T) // Function to set up environment for the test
		expectError bool
		check       func(t *testing.T, cfg *Config) // Function to check the results
	}{
		{
			name: "Happy Path - all variables set",
			setup: func(t *testing.T) {
				t.Setenv("DATABASE_URL", "postgres://user:pass@host/db")
				t.Setenv("GCS_BUCKET_NAME", "my-test-bucket")
				t.Setenv("APP_ENV", "production")
			},
			expectError: false,
			check: func(t *testing.T, cfg *Config) {
				if cfg.DatabaseURL != "postgres://user:pass@host/db" {
					t.Errorf("expected DatabaseURL to be 'postgres://user:pass@host/db', got '%s'", cfg.DatabaseURL)
				}
				if cfg.GCSBucketName != "my-test-bucket" {
					t.Errorf("expected GCSBucketName to be 'my-test-bucket', got '%s'", cfg.GCSBucketName)
				}
				if cfg.AppEnv != "production" {
					t.Errorf("expected AppEnv to be 'production', got '%s'", cfg.AppEnv)
				}
			},
		},
		{
			name: "Failure - DATABASE_URL is missing",
			setup: func(t *testing.T) {
				t.Setenv("GCS_BUCKET_NAME", "my-test-bucket")
			},
			expectError: true,
		},
		{
			name: "Failure - GCS_BUCKET_NAME is missing",
			setup: func(t *testing.T) {
				t.Setenv("DATABASE_URL", "postgres://user:pass@host/db")
			},
			expectError: true,
		},
		{
			name: "Defaulting - APP_ENV defaults to development",
			setup: func(t *testing.T) {
				t.Setenv("DATABASE_URL", "postgres://user:pass@host/db")
				t.Setenv("GCS_BUCKET_NAME", "my-test-bucket")
			},
			expectError: false,
			check: func(t *testing.T, cfg *Config) {
				if cfg.AppEnv != "development" {
					t.Errorf("expected AppEnv to default to 'development', got '%s'", cfg.AppEnv)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(t)

			cfg, err := LoadConfig()

			if tc.expectError {
				if err == nil {
					t.Error("expected an error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error, but got: %v", err)
				}
				if tc.check != nil {
					tc.check(t, cfg)
				}
			}
		})
	}
}
