package ibkr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// MarketDataSnapshot represents a market data snapshot.
type MarketDataSnapshot struct {
	ConID     int     `json:"conid"`
	ConIDEx   string  `json:"conidEx"`
	LastPrice float64 `json:"31,omitempty"`   // Last price.
	Symbol    string  `json:"55,omitempty"`   // Symbol.
	Bid       float64 `json:"84,omitempty"`   // Bid.
	Ask       float64 `json:"86,omitempty"`   // Ask.
	Volume    int64   `json:"87,omitempty"`   // Volume.
	High      float64 `json:"70,omitempty"`   // High.
	Low       float64 `json:"71,omitempty"`   // Low.
	Close     float64 `json:"82,omitempty"`   // Close.
	Open      float64 `json:"7295,omitempty"` // Open.
	ServerID  string  `json:"_updated"`
}

// HistoricalDataResponse represents historical market data.
type HistoricalDataResponse struct {
	ServerID           string          `json:"serverId"`
	Symbol             string          `json:"symbol"`
	Text               string          `json:"text"`
	PriceFactor        int             `json:"priceFactor"`
	StartTime          string          `json:"startTime"`
	High               string          `json:"high"`
	Low                string          `json:"low"`
	TimePeriod         string          `json:"timePeriod"`
	BarLength          int             `json:"barLength"`
	MdAvailability     string          `json:"mdAvailability"`
	MktDataDelay       int             `json:"mktDataDelay"`
	OutsideRth         bool            `json:"outsideRth"`
	TradingDayDuration int             `json:"tradingDayDuration"`
	VolumeFactor       int             `json:"volumeFactor"`
	PriceDisplayRule   int             `json:"priceDisplayRule"`
	PriceDisplayValue  string          `json:"priceDisplayValue"`
	NegativeCapable    bool            `json:"negativeCapable"`
	MessageVersion     int             `json:"messageVersion"`
	Data               []HistoricalBar `json:"data"`
	Points             int             `json:"points"`
}

// HistoricalBar represents a single historical bar.
type HistoricalBar struct {
	Time   int64   `json:"t"` // Timestamp.
	Open   float64 `json:"o"`
	Close  float64 `json:"c"`
	High   float64 `json:"h"`
	Low    float64 `json:"l"`
	Volume int64   `json:"v"`
}

// Contract represents a contract search result.
type Contract struct {
	ConID         int               `json:"conid"`
	CompanyHeader string            `json:"companyHeader"`
	CompanyName   string            `json:"companyName"`
	Symbol        string            `json:"symbol"`
	Description   string            `json:"description"`
	Restricted    string            `json:"restricted"`
	Fop           string            `json:"fop"`
	Opt           string            `json:"opt"`
	War           string            `json:"war"`
	Sections      []ContractSection `json:"sections"`
}

// ContractSection represents a section in contract search results.
type ContractSection struct {
	SecType    string `json:"secType"`
	Months     string `json:"months"`
	Symbol     string `json:"symbol"`
	Exchange   string `json:"exchange"`
	LegSecType string `json:"legSecType,omitempty"`
}

// GetMarketData retrieves market data snapshot for a contract.
func (c *Client) GetMarketData(ctx context.Context, conIDs []int, fields []string) ([]MarketDataSnapshot, error) {
	conIDsStr := make([]string, 0, len(conIDs))
	for _, id := range conIDs {
		conIDsStr = append(conIDsStr, fmt.Sprintf("%d", id))
	}

	params := url.Values{}
	params.Set("conids", strings.Join(conIDsStr, ","))

	if len(fields) > 0 {
		params.Set("fields", strings.Join(fields, ","))
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/v1/api/iserver/marketdata/snapshot?%s", c.baseURL, params.Encode()),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get market data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("get market data failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var snapshots []MarketDataSnapshot
	if err := json.NewDecoder(resp.Body).Decode(&snapshots); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return snapshots, nil
}

// GetHistoricalData retrieves historical market data.
func (c *Client) GetHistoricalData(
	ctx context.Context,
	conID int,
	period, barSize string,
) (*HistoricalDataResponse, error) {
	params := url.Values{}
	params.Set("conid", fmt.Sprintf("%d", conID))
	params.Set("period", period)
	params.Set("bar", barSize)

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/v1/api/iserver/marketdata/history?%s", c.baseURL, params.Encode()),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("get historical data failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var histData HistoricalDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&histData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &histData, nil
}

// SearchContracts searches for contracts by symbol.
func (c *Client) SearchContracts(ctx context.Context, symbol string) ([]Contract, error) {
	params := url.Values{}
	params.Set("symbol", symbol)

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/v1/api/iserver/secdef/search?%s", c.baseURL, params.Encode()),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to search contracts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("search contracts failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var contracts []Contract
	if err := json.NewDecoder(resp.Body).Decode(&contracts); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return contracts, nil
}
