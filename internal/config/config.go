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
	SMTPHost               string
	SMTPPort               string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	requiredEnvVars := map[string]string{
		"ORDERSPACE_CLIENT_ID":     os.Getenv("ORDERSPACE_CLIENT_ID"),
		"ORDERSPACE_CLIENT_SECRET": os.Getenv("ORDERSPACE_CLIENT_SECRET"),
	}

	for key, value := range requiredEnvVars {
		if value == "" {
			return nil, fmt.Errorf("%s is required", key)
		}
	}

	postmarkToken := os.Getenv("POSTMARK_SERVER_TOKEN")
	if smtpHost == "" && postmarkToken == "" {
		return nil, fmt.Errorf("either SMTP_HOST or POSTMARK_SERVER_TOKEN is required")
	}

	return &Config{
		OrderspaceClientID:     requiredEnvVars["ORDERSPACE_CLIENT_ID"],
		OrderspaceClientSecret: requiredEnvVars["ORDERSPACE_CLIENT_SECRET"],
		PostmarkServerToken:    postmarkToken,
		DatabaseURL:            os.Getenv("DATABASE_URL"),
		SMTPHost:               smtpHost,
		SMTPPort:               smtpPort,
	}, nil
}
