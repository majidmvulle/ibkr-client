package ibkr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// PlaceOrderRequest represents a request to place an order.
type PlaceOrderRequest struct {
	ConID     int     `json:"conid"`
	SecType   string  `json:"secType"`
	OrderType string  `json:"orderType"`
	Side      string  `json:"side"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price,omitempty"`
	Tif       string  `json:"tif"`
	Ticker    string  `json:"ticker"`
}

// ModifyOrderRequest represents a request to modify an order.
type ModifyOrderRequest struct {
	Quantity float64 `json:"quantity,omitempty"`
	Price    float64 `json:"price,omitempty"`
}

// OrderResponse represents an order response from the Gateway.
type OrderResponse struct {
	OrderID     string   `json:"order_id"`
	OrderStatus string   `json:"order_status"`
	EncryptedID string   `json:"id"`
	Message     []string `json:"message"`
}

// Order represents an order from the Gateway.
type Order struct {
	AcctID            string  `json:"acct"`
	ConID             int     `json:"conid"`
	OrderID           string  `json:"orderId"`
	CashCcy           string  `json:"cashCcy"`
	SizeAndFills      string  `json:"sizeAndFills"`
	OrderDesc         string  `json:"orderDesc"`
	Description1      string  `json:"description1"`
	Ticker            string  `json:"ticker"`
	SecType           string  `json:"secType"`
	ListingExchange   string  `json:"listingExchange"`
	RemainingQuantity float64 `json:"remainingQuantity"`
	FilledQuantity    float64 `json:"filledQuantity"`
	TotalSize         float64 `json:"totalSize"`
	CompanyName       string  `json:"companyName"`
	Status            string  `json:"status"`
	OrigOrderType     string  `json:"origOrderType"`
	Side              string  `json:"side"`
	Price             float64 `json:"price"`
	BgColor           string  `json:"bgColor"`
	FgColor           string  `json:"fgColor"`
}

// PlaceOrder places a new order.
func (c *Client) PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*OrderResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/v1/api/iserver/account/%s/orders", c.baseURL, c.accountID),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to place order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("place order failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var orderResp OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &orderResp, nil
}

// ModifyOrder modifies an existing order.
func (c *Client) ModifyOrder(ctx context.Context, orderID string, req *ModifyOrderRequest) (*OrderResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/v1/api/iserver/account/%s/order/%s", c.baseURL, c.accountID, orderID),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to modify order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("modify order failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var orderResp OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &orderResp, nil
}

// CancelOrder cancels an order.
func (c *Client) CancelOrder(ctx context.Context, orderID string) error {
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/v1/api/iserver/account/%s/order/%s", c.baseURL, c.accountID, orderID),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("cancel order failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// GetLiveOrders retrieves live orders.
func (c *Client) GetLiveOrders(ctx context.Context) ([]Order, error) {
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/v1/api/iserver/account/%s/orders", c.baseURL, c.accountID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("get orders failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var orders struct {
		Orders []Order `json:"orders"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return orders.Orders, nil
}
