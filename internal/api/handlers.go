package api

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DukeRupert/rr/internal/email"
	"github.com/DukeRupert/rr/internal/orderspace"
	"github.com/labstack/echo/v4"
)

type AdHocEmailRequest struct {
	Subject  string `json:"subject"`
	HtmlBody string `json:"htmlBody"`
	TextBody string `json:"textBody"`
}

type AdHocEmailResponse struct {
	Sent    int      `json:"sent"`
	Failed  int      `json:"failed"`
	Skipped int      `json:"skipped"`
	Details []string `json:"details"`
}

type Handler struct {
	client *orderspace.Client
	email  email.Sender
	db     *sql.DB
}

func NewHandler(client *orderspace.Client, emailClient email.Sender, db *sql.DB) *Handler {
	return &Handler{client: client, email: emailClient, db: db}
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

func (h *Handler) SendAdHocEmail(c echo.Context) error {
	var req AdHocEmailRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	if req.Subject == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "subject is required")
	}
	if req.HtmlBody == "" && req.TextBody == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "htmlBody or textBody is required")
	}

	sixWeeksAgo := time.Now().AddDate(0, 0, -42)
	params := &orderspace.CustomerListParams{
		UpdatedSince: &sixWeeksAgo,
	}

	resp, err := h.client.ListCustomers(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch customers: "+err.Error())
	}

	result := AdHocEmailResponse{
		Details: []string{},
	}

	for _, customer := range resp.Customers {
		var notifyDays bool
		err := h.db.QueryRow(`
			SELECT COALESCE(
				(SELECT email_notify_days FROM customer_notifications WHERE customer_id = ?),
				true
			)
		`, customer.ID).Scan(&notifyDays)
		if err != nil {
			log.Printf("ERROR checking notification preference for %s: %v", customer.CompanyName, err)
			result.Failed++
			result.Details = append(result.Details, "ERROR: "+customer.CompanyName+" (failed to check preferences)")
			continue
		}

		if !notifyDays {
			result.Skipped++
			result.Details = append(result.Details, "SKIPPED: "+customer.CompanyName+" (notifications disabled)")
			continue
		}

		adHocEmail := email.Email{
			From:     "info@rockabillyroasting.com",
			To:       customer.EmailAddresses.Orders,
			Subject:  req.Subject,
			HtmlBody: req.HtmlBody,
			TextBody: req.TextBody,
		}

		_, err = h.email.SendEmail(adHocEmail)
		if err != nil {
			log.Printf("ERROR sending ad-hoc email to %s: %v", customer.CompanyName, err)
			result.Failed++
			result.Details = append(result.Details, "ERROR: "+customer.CompanyName+" ("+err.Error()+")")
		} else {
			log.Printf("SUCCESS sent ad-hoc email to %s (%s)", customer.CompanyName, customer.EmailAddresses.Orders)
			result.Sent++
			result.Details = append(result.Details, "SUCCESS: "+customer.CompanyName+" ("+customer.EmailAddresses.Orders+")")
		}
	}

	return c.JSON(http.StatusOK, result)
}
