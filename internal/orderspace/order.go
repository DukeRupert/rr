package orderspace

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/DukeRupert/rr/internal/models"
)

// OrderRequest represents the structure of the order creation request
type OrderRequest struct {
	Order OrderRequestBody `json:"order"`
}

// OrderRequestBody represents the body of the order creation request
type OrderRequestBody struct {
	CustomerID       string             `json:"customer_id"`
	DeliveryDate     string             `json:"delivery_date"` // Format: "2006-01-02"
	Reference        string             `json:"reference,omitempty"`
	InternalNote     string             `json:"internal_note,omitempty"`
	CustomerPONumber string             `json:"customer_po_number,omitempty"`
	CustomerNote     string             `json:"customer_note,omitempty"`
	ShippingAddress  models.Address     `json:"shipping_address"`
	BillingAddress   models.Address     `json:"billing_address"`
	OrderLines       []OrderLineRequest `json:"order_lines"`
}

// OrderLineRequest represents an order line in the creation request
type OrderLineRequest struct {
	SKU       string   `json:"sku,omitempty"`
	Name      string   `json:"name,omitempty"`
	Quantity  int      `json:"quantity"`
	UnitPrice *float64 `json:"unit_price,omitempty"`
	Shipping  *bool    `json:"shipping,omitempty"`
	TaxRateID *string  `json:"tax_rate_id,omitempty"`
}

// OrderError represents an error response from the API
type OrderError struct {
	Message string `json:"message"`
}

// CreateOrder sends a request to create a new order
func (c *Client) CreateOrder(req *OrderRequest) (*models.Order, error) {
	// Convert request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Make authenticated request
	resp, err := c.MakeAuthenticatedRequest("POST", "/orders", jsonData)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Handle different response status codes
	switch resp.StatusCode {
	case http.StatusOK:
		var successResp struct {
			Order models.Order `json:"order"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&successResp); err != nil {
			return nil, fmt.Errorf("error decoding response: %w", err)
		}
		return &successResp.Order, nil
	case http.StatusUnprocessableEntity:
		var errorResp OrderError
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("error decoding error response: %w", err)
		}
		return nil, fmt.Errorf("order validation failed: %s", errorResp.Message)
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

// Helper function to create a shipping line item
func NewShippingLine(name string, unitPrice float64) OrderLineRequest {
	shipping := true
	return OrderLineRequest{
		Name:      name,
		UnitPrice: &unitPrice,
		Quantity:  1,
		Shipping:  &shipping,
	}
}

// Helper function to create a product line item
func NewProductLine(sku string, quantity int) OrderLineRequest {
	return OrderLineRequest{
		SKU:      sku,
		Quantity: quantity,
	}
}

// Helper function to create a custom line item
func NewCustomLine(name string, quantity int, unitPrice float64, taxRateID string) OrderLineRequest {
	return OrderLineRequest{
		Name:      name,
		Quantity:  quantity,
		UnitPrice: &unitPrice,
		TaxRateID: &taxRateID,
	}
}

type OrderListParams struct {
	StartingAfter      string     `url:"starting_after,omitempty"`
	Limit              int        `url:"limit,omitempty"`
	CreatedSince       *time.Time `url:"-"`
	CreatedBefore      *time.Time `url:"-"`
	DeliveryDateSince  *time.Time `url:"-"`
	DeliveryDateBefore *time.Time `url:"-"`
	Number             int        `url:"number,omitempty"`
	Status             string     `url:"status,omitempty"`
	Reference          string     `url:"reference,omitempty"`
	CustomerID         string     `url:"customer_id,omitempty"`
	StandingOrderID    string     `url:"standing_order_id,omitempty"`
}

type OrderListResponse struct {
	Orders  []models.Order `json:"orders"`
	HasMore bool           `json:"has_more"`
}

func (c *Client) ListOrders(params *OrderListParams) (*OrderListResponse, error) {
	basePath := "/orders"

	if params != nil {
		u, err := url.Parse(basePath)
		if err != nil {
			return nil, fmt.Errorf("error parsing URL: %w", err)
		}

		q := u.Query()
		if params.StartingAfter != "" {
			q.Set("starting_after", params.StartingAfter)
		}
		if params.Limit > 0 {
			q.Set("limit", fmt.Sprintf("%d", params.Limit))
		}
		if params.CreatedSince != nil {
			q.Set("created_since", params.CreatedSince.Format(time.RFC3339))
		}
		if params.CreatedBefore != nil {
			q.Set("created_before", params.CreatedBefore.Format(time.RFC3339))
		}
		if params.DeliveryDateSince != nil {
			q.Set("delivery_date_since", params.DeliveryDateSince.Format("2006-01-02"))
		}
		if params.DeliveryDateBefore != nil {
			q.Set("delivery_date_before", params.DeliveryDateBefore.Format("2006-01-02"))
		}
		if params.Number > 0 {
			q.Set("number", fmt.Sprintf("%d", params.Number))
		}
		if params.Status != "" {
			q.Set("status", params.Status)
		}
		if params.Reference != "" {
			q.Set("reference", params.Reference)
		}
		if params.CustomerID != "" {
			q.Set("customer_id", params.CustomerID)
		}
		if params.StandingOrderID != "" {
			q.Set("standing_order_id", params.StandingOrderID)
		}

		basePath = fmt.Sprintf("%s?%s", basePath, q.Encode())
	}

	resp, err := c.MakeAuthenticatedRequest("GET", basePath, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result OrderListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}
