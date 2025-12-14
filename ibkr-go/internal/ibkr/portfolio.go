package ibkr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Position represents a portfolio position.
type Position struct {
	AcctID        string   `json:"acctId"`
	ConID         int      `json:"conid"`
	ContractDesc  string   `json:"contractDesc"`
	Position      float64  `json:"position"`
	MktPrice      float64  `json:"mktPrice"`
	MktValue      float64  `json:"mktValue"`
	Currency      string   `json:"currency"`
	AvgCost       float64  `json:"avgCost"`
	AvgPrice      float64  `json:"avgPrice"`
	RealizedPnl   float64  `json:"realizedPnl"`
	UnrealizedPnl float64  `json:"unrealizedPnl"`
	ExcRate       float64  `json:"excRate"`
	ExpDate       string   `json:"expDate"`
	PutOrCall     string   `json:"putOrCall"`
	Multiplier    float64  `json:"multiplier"`
	Strike        float64  `json:"strike"`
	ExerciseStyle string   `json:"exerciseStyle"`
	ConExchMap    []string `json:"conExchMap"`
	AssetClass    string   `json:"assetClass"`
	UndConID      int      `json:"undConid"`
}

// AccountSummary represents account summary information.
type AccountSummary struct {
	AccountID           string  `json:"accountcode"`
	AccountType         string  `json:"accounttype"`
	NetLiquidation      float64 `json:"netliquidation"`
	TotalCashValue      float64 `json:"totalcashvalue"`
	SettledCash         float64 `json:"settledcash"`
	AccruedCash         float64 `json:"accruedcash"`
	BuyingPower         float64 `json:"buyingpower"`
	EquityWithLoanValue float64 `json:"equitywithloanvalue"`
	PreviousDayEquity   float64 `json:"previousdayequitywithloanvalue"`
	GrossPositionValue  float64 `json:"grosspositionvalue"`
	RegTEquity          float64 `json:"regtequity"`
	RegTMargin          float64 `json:"regtmargin"`
	SMA                 float64 `json:"sma"`
	Currency            string  `json:"currency"`
}

// GetPortfolio retrieves portfolio positions.
func (c *Client) GetPortfolio(ctx context.Context) ([]Position, error) {
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/v1/api/portfolio/%s/positions/0", c.baseURL, c.accountID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("get portfolio failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var positions []Position
	if err := json.NewDecoder(resp.Body).Decode(&positions); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return positions, nil
}

// GetAccountSummary retrieves account summary information.
func (c *Client) GetAccountSummary(ctx context.Context) (*AccountSummary, error) {
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/v1/api/portfolio/%s/summary", c.baseURL, c.accountID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get account summary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("get account summary failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var summary AccountSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &summary, nil
}
