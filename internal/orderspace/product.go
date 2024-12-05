package orderspace

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/DukeRupert/rr/internal/models"
)

// ProductsResponse represents the API response for listing products
type ProductsResponse struct {
	Products []models.Product `json:"products"`
	HasMore  bool             `json:"has_more"`
}

// ProductListParams contains all possible parameters for listing products
type ProductListParams struct {
	StartingAfter string     `url:"starting_after,omitempty"`
	Limit         int        `url:"limit,omitempty"`
	CreatedSince  *time.Time `url:"-"` // Handled separately due to formatting
	UpdatedSince  *time.Time `url:"-"` // Handled separately due to formatting
	Code          string     `url:"code,omitempty"`
	Name          string     `url:"name,omitempty"`
	Active        *bool      `url:"active,omitempty"`
	CategoryID    string     `url:"category_id,omitempty"`
}

// ListProducts retrieves a list of products with optional filtering
func (c *Client) ListProducts(params *ProductListParams) (*ProductsResponse, error) {
	// Build base path and query parameters
	basePath := "/products"

	// Add query parameters if they exist
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
		if params.UpdatedSince != nil {
			q.Set("updated_since", params.UpdatedSince.Format(time.RFC3339))
		}
		if params.Code != "" {
			q.Set("code", params.Code)
		}
		if params.Name != "" {
			q.Set("name", params.Name)
		}
		if params.Active != nil {
			q.Set("active", fmt.Sprintf("%v", *params.Active))
		}
		if params.CategoryID != "" {
			q.Set("category_id", params.CategoryID)
		}

		basePath = fmt.Sprintf("%s?%s", basePath, q.Encode())
	}

	// Make authenticated request
	resp, err := c.MakeAuthenticatedRequest("GET", basePath, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var result ProductsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// Helper function to create boolean pointer
func BoolPtr(b bool) *bool {
	return &b
}

// Helper function to create time pointer
func TimePtr(t time.Time) *time.Time {
	return &t
}
