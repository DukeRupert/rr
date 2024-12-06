package api

import (
	"database/sql"
	"net/http"

	"github.com/DukeRupert/rr/internal/email"
	"github.com/DukeRupert/rr/internal/orderspace"
	"github.com/DukeRupert/rr/internal/services"

	"github.com/labstack/echo/v4"
)

// routes.go
func SetupRoutes(e *echo.Echo, client *orderspace.Client, email *email.Client, db *sql.DB) {
	h := NewHandler(client, email, db)
	e.GET("/api/customers", h.GetCustomers)
	e.GET("/api/orders", h.GetOrders)
	e.GET("/api/email/preview-reminders", func(c echo.Context) error {
		if err := services.PreviewOrderReminders(db, client, email); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "preview sent"})
	})
}
