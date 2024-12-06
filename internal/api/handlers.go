package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/DukeRupert/rr/internal/email"
	"github.com/DukeRupert/rr/internal/orderspace"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	client *orderspace.Client
	email  *email.Client
	db     *sql.DB
}

func NewHandler(client *orderspace.Client, email *email.Client, db *sql.DB) *Handler {
	return &Handler{client: client, email: email, db: db}
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

func (h *Handler) GetOrders(c echo.Context) error {
	params := &orderspace.OrderListParams{}

	// Parse query parameters
	if limit := c.QueryParam("limit"); limit != "" {
		if n, err := strconv.Atoi(limit); err == nil {
			params.Limit = n
		}
	}

	params.StartingAfter = c.QueryParam("starting_after")
	params.Status = c.QueryParam("status")
	params.Reference = c.QueryParam("reference")
	params.CustomerID = c.QueryParam("customer_id")
	params.StandingOrderID = c.QueryParam("standing_order_id")

	// Parse date parameters
	if created_since := c.QueryParam("created_since"); created_since != "" {
		if t, err := time.Parse(time.RFC3339, created_since); err == nil {
			params.CreatedSince = &t
		}
	}

	if created_before := c.QueryParam("created_before"); created_before != "" {
		if t, err := time.Parse(time.RFC3339, created_before); err == nil {
			params.CreatedBefore = &t
		}
	}

	if delivery_date_since := c.QueryParam("delivery_date_since"); delivery_date_since != "" {
		if t, err := time.Parse("2006-01-02", delivery_date_since); err == nil {
			params.DeliveryDateSince = &t
		}
	}

	if delivery_date_before := c.QueryParam("delivery_date_before"); delivery_date_before != "" {
		if t, err := time.Parse("2006-01-02", delivery_date_before); err == nil {
			params.DeliveryDateBefore = &t
		}
	}

	// Call the OrderSpace API
	response, err := h.client.ListOrders(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, response)
}
