package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/DukeRupert/rr/internal/orderspace"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	client *orderspace.Client
	db     *sql.DB
}

func NewHandler(client *orderspace.Client, db *sql.DB) *Handler {
	return &Handler{client: client, db: db}
}

func (h *Handler) GetCustomers(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 50 // default limit
	}

	params := &orderspace.CustomerListParams{
		Limit:         limit,
		StartingAfter: c.QueryParam("starting_after"),
	}

	customers, err := h.client.ListCustomers(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch customers: "+err.Error())
	}

	return c.JSON(http.StatusOK, customers)
}

// routes.go
func SetupRoutes(e *echo.Echo, client *orderspace.Client, db *sql.DB) {
	h := NewHandler(client, db)
	e.GET("/api/customers", h.GetCustomers)
}
