package api

import (
	"database/sql"

	"github.com/DukeRupert/rr/internal/email"
	"github.com/DukeRupert/rr/internal/orderspace"
	"github.com/labstack/echo/v4"
)

// routes.go
func SetupRoutes(e *echo.Echo, client *orderspace.Client, email *email.Client, db *sql.DB) {
	h := NewHandler(client, email, db)
	e.GET("/api/customers", h.GetCustomers)
	e.GET("/api/orders", h.GetOrders)
}
