package ibkr

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	// DefaultHTTPTimeout is the default timeout for HTTP requests.
	DefaultHTTPTimeout = 30 * time.Second
)

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

// Additional Gateway API methods will be implemented in Phase 4:
// - AuthStatus() - Check authentication status
// - GetAccounts() - Get available accounts
// - GetPortfolio() - Get portfolio positions
// - PlaceOrder() - Place an order
// - GetMarketData() - Get market data.
