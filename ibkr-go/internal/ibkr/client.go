package ibkr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// DefaultHTTPTimeout is the default timeout for HTTP requests.
	DefaultHTTPTimeout = 30 * time.Second
)

// AuthStatusResponse represents the authentication status response.
type AuthStatusResponse struct {
	Authenticated bool   `json:"authenticated"`
	Competing     bool   `json:"competing"`
	Connected     bool   `json:"connected"`
	Message       string `json:"message"`
	MAC           string `json:"MAC"`
	ServerInfo    struct {
		ServerName    string `json:"serverName"`
		ServerVersion string `json:"serverVersion"`
	} `json:"serverInfo"`
}

// Account represents an IBKR account.
type Account struct {
	ID             string `json:"id"`
	AccountID      string `json:"accountId"`
	AccountVan     string `json:"accountVan"`
	AccountTitle   string `json:"accountTitle"`
	DisplayName    string `json:"displayName"`
	AccountAlias   string `json:"accountAlias"`
	AccountStatus  string `json:"accountStatus"`
	Currency       string `json:"currency"`
	Type           string `json:"type"`
	TradingType    string `json:"tradingType"`
	Faclient       bool   `json:"faclient"`
	ClearingStatus string `json:"clearingStatus"`
	Covestor       bool   `json:"covestor"`
	Parent         struct {
		MmcID       int    `json:"mmc"`
		AccountID   string `json:"accountId"`
		IsMParent   bool   `json:"isMParent"`
		IsMChild    bool   `json:"isMChild"`
		IsMultiplex bool   `json:"isMultiplex"`
	} `json:"parent"`
}

// Client is an HTTP client for the IBKR Client Portal Gateway API.
type Client struct {
	baseURL    string
	httpClient *http.Client
	accountID  string
}

// NewClient creates a new IBKR Gateway client.
func NewClient(baseURL, accountID string) *Client {
	return &Client{
		baseURL:   baseURL,
		accountID: accountID,
		httpClient: &http.Client{
			Timeout: DefaultHTTPTimeout,
		},
	}
}

// Ping checks if the Gateway is accessible.
func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/api/tickle", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to ping gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gateway returned status %d", resp.StatusCode)
	}

	return nil
}

// AuthStatus checks the current authentication status.
func (c *Client) AuthStatus(ctx context.Context) (*AuthStatusResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/api/iserver/auth/status", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to check auth status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("auth status check failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authStatus AuthStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&authStatus); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &authStatus, nil
}

// Reauthenticate triggers reauthentication.
func (c *Client) Reauthenticate(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/api/iserver/reauthenticate", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to reauthenticate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("reauthentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetAccounts retrieves the list of accounts.
func (c *Client) GetAccounts(ctx context.Context) ([]Account, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/api/portfolio/accounts", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("get accounts failed with status %d: %s", resp.StatusCode, string(body))
	}

	var accounts []Account
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return accounts, nil
}
