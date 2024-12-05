// config/config.go
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OrderspaceClientID     string
	OrderspaceClientSecret string
	DatabaseURL            string
	PostmarkServerToken    string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	requiredEnvVars := map[string]string{
		"ORDERSPACE_CLIENT_ID":     os.Getenv("ORDERSPACE_CLIENT_ID"),
		"ORDERSPACE_CLIENT_SECRET": os.Getenv("ORDERSPACE_CLIENT_SECRET"),
		"POSTMARK_SERVER_TOKEN":    os.Getenv("POSTMARK_SERVER_TOKEN"),
	}

	for key, value := range requiredEnvVars {
		if value == "" {
			return nil, fmt.Errorf("%s is required", key)
		}
	}

	return &Config{
		OrderspaceClientID:     requiredEnvVars["ORDERSPACE_CLIENT_ID"],
		OrderspaceClientSecret: requiredEnvVars["ORDERSPACE_CLIENT_SECRET"],
		PostmarkServerToken:    requiredEnvVars["POSTMARK_SERVER_TOKEN"],
		DatabaseURL:            os.Getenv("DATABASE_URL"),
	}, nil
}
