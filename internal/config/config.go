package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OrderspaceClientID     string
	OrderspaceClientSecret string
	DatabaseURL            string // Adding this for future database configuration
}

func Load() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Get environment variables
	clientID := os.Getenv("ORDERSPACE_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("ORDERSPACE_CLIENT_ID is required")
	}

	clientSecret := os.Getenv("ORDERSPACE_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, fmt.Errorf("ORDERSPACE_CLIENT_SECRET is required")
	}

	// Create config struct
	config := &Config{
		OrderspaceClientID:     clientID,
		OrderspaceClientSecret: clientSecret,
		DatabaseURL:            os.Getenv("DATABASE_URL"), // Default to empty string if not set
	}

	return config, nil
}
