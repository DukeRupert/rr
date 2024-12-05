package orderspace

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/DukeRupert/rr/internal/models"
)

type CustomerListParams struct {
	StartingAfter string     `url:"starting_after,omitempty"`
	Limit         int        `url:"limit,omitempty"`
	CreatedSince  *time.Time `url:"-"`
	UpdatedSince  *time.Time `url:"-"`
	Status        string     `url:"status,omitempty"`
}

type CustomerResponse struct {
	Customers []models.Customer `json:"customers"`
	HasMore   bool              `json:"has_more"`
}

func (c *Client) ListCustomers(params *CustomerListParams) (*CustomerResponse, error) {
	// Build query parameters
	baseURL := "/customers"
	if params != nil {
		u, err := url.Parse(baseURL)
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
		if params.Status != "" {
			q.Set("status", params.Status)
		}

		baseURL = fmt.Sprintf("%s?%s", baseURL, q.Encode())
	}

	// Make authenticated request
	resp, err := c.MakeAuthenticatedRequest("GET", baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result CustomerResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}
