package orderspace

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// AuthResponse represents the OAuth token response
type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// TokenInfo stores token data with expiration
type TokenInfo struct {
	Token     string
	CreatedAt time.Time
}

// Client represents the API client with auth capabilities
type Client struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client
	DB           *sql.DB
}

// NewClient creates a new API client with database connection
func NewClient(id, secret string, db *sql.DB) (*Client, error) {
	client := &Client{
		BaseURL:      "https://api.orderspace.com/v1",
		ClientID:     id,
		ClientSecret: secret,
		HTTPClient:   &http.Client{Timeout: time.Second * 30},
		DB:           db,
	}

	// Initialize the tokens table
	err := client.initDB()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	return client, nil
}

// initDB creates the tokens table if it doesn't exist
func (c *Client) initDB() error {
	_, err := c.DB.Exec(`
        CREATE TABLE IF NOT EXISTS tokens (
            id INTEGER PRIMARY KEY,
            access_token TEXT NOT NULL,
            created_at DATETIME NOT NULL
        )
    `)
	return err
}

// GetValidToken returns a valid token or obtains a new one if necessary
func (c *Client) GetValidToken() (string, error) {
	// Try to get existing valid token
	var token TokenInfo
	err := c.DB.QueryRow(`
        SELECT access_token, created_at 
        FROM tokens 
        ORDER BY created_at DESC 
        LIMIT 1
    `).Scan(&token.Token, &token.CreatedAt)

	if err == nil {
		// Check if token is still valid (less than 25 minutes old to add buffer)
		if time.Since(token.CreatedAt) < 25*time.Minute {
			return token.Token, nil
		}
	}

	// Get new token if none exists or current one is expired
	return c.refreshToken()
}

// refreshToken obtains a new access token from the auth endpoint
func (c *Client) refreshToken() (string, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST",
		"https://identity.orderspace.com/oauth/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth request failed with status: %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	// Store new token in database
	_, err = c.DB.Exec(`
        INSERT INTO tokens (access_token, created_at) 
        VALUES (?, ?)
    `, authResp.AccessToken, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to store token: %v", err)
	}

	return authResp.AccessToken, nil
}

// MakeAuthenticatedRequest makes a request with the current valid token
func (c *Client) MakeAuthenticatedRequest(method, path string, body []byte) (*http.Response, error) {
	token, err := c.GetValidToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get valid token: %v", err)
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}
