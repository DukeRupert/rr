package main

import (
	"log"

	"github.com/DukeRupert/rr/internal/api"
	"github.com/DukeRupert/rr/internal/config"
	"github.com/DukeRupert/rr/internal/database"
	"github.com/DukeRupert/rr/internal/email"
	"github.com/DukeRupert/rr/internal/orderspace"
	"github.com/DukeRupert/rr/internal/services"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize Echo
	e := echo.New()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Postmark client
	postmarkClient := email.NewClient(cfg.PostmarkServerToken)

	// Initialize Orderspace client
	orderspaceClient, err := orderspace.NewClient(cfg.OrderspaceClientID, cfg.OrderspaceClientSecret, db)
	if err != nil {
		log.Fatal(err)
	}

	// Setup routes
	api.SetupRoutes(e, orderspaceClient, postmarkClient, db)

	// Initialize reminder service
	reminderService, err := services.NewReminderScheduler(db, orderspaceClient, postmarkClient)
	if err != nil {
		log.Fatalf("Failed to create reminder service: %v", err)
	}
	reminderService.Start()
	defer reminderService.Shutdown()

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
