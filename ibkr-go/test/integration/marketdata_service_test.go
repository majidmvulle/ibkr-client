//go:build integration

package integration

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/api"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	marketdatav1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/marketdata/v1"
)

func TestIntegration_MarketDataService_GetQuote(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create test session
	token := CreateTestSession(t, testCtx.Config.IBKRAccountID)
	defer DeleteTestSession(t, token)

	// Add account ID to context
	ctx = middleware.SetAccountIDInContext(ctx, testCtx.Config.IBKRAccountID)

	// Create market data service handler
	handler := api.NewMarketDataServiceHandler(testCtx.IBKRClient)

	// Create get quote request
	req := connect.NewRequest(&marketdatav1.GetQuoteRequest{
		Symbol: "AAPL",
	})

	// Get quote
	resp, err := handler.GetQuote(ctx, req)
	if err != nil {
		t.Fatalf("GetQuote failed: %v", err)
	}

	// Verify response
	if resp.Msg.Quote == nil {
		t.Fatal("Expected quote to be set")
	}

	quote := resp.Msg.Quote

	if quote.Symbol != "AAPL" {
		t.Errorf("Expected symbol AAPL, got %s", quote.Symbol)
	}

	if quote.Last == 0 {
		t.Error("Expected last price to be set")
	}

	if quote.Bid == 0 {
		t.Error("Expected bid price to be set")
	}

	if quote.Ask == 0 {
		t.Error("Expected ask price to be set")
	}

	t.Logf("Quote retrieved: Symbol=%s, Last=%.2f, Bid=%.2f, Ask=%.2f",
		quote.Symbol, quote.Last, quote.Bid, quote.Ask)
}

func TestIntegration_MarketDataService_GetHistoricalData(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create test session
	token := CreateTestSession(t, testCtx.Config.IBKRAccountID)
	defer DeleteTestSession(t, token)

	// Add account ID to context
	ctx = middleware.SetAccountIDInContext(ctx, testCtx.Config.IBKRAccountID)

	// Create market data service handler
	handler := api.NewMarketDataServiceHandler(testCtx.IBKRClient)

	// Create get historical data request
	req := connect.NewRequest(&marketdatav1.GetHistoricalDataRequest{
		Symbol:  "AAPL",
		Period:  "1d",
		BarSize: "5min",
	})

	// Get historical data
	resp, err := handler.GetHistoricalData(ctx, req)
	if err != nil {
		t.Fatalf("GetHistoricalData failed: %v", err)
	}

	// Verify response
	if resp.Msg.Bars == nil {
		t.Error("Expected bars list to be initialized")
	}

	t.Logf("Retrieved %d historical bars", len(resp.Msg.Bars))

	// Verify bar structure if any exist
	for i, bar := range resp.Msg.Bars {
		if bar.Timestamp == "" {
			t.Errorf("Bar %d: expected timestamp to be set", i)
		}
		if bar.Open == 0 {
			t.Errorf("Bar %d: expected open price to be set", i)
		}
		if bar.High == 0 {
			t.Errorf("Bar %d: expected high price to be set", i)
		}
		if bar.Low == 0 {
			t.Errorf("Bar %d: expected low price to be set", i)
		}
		if bar.Close == 0 {
			t.Errorf("Bar %d: expected close price to be set", i)
		}
	}
}

func TestIntegration_MarketDataService_InvalidSymbol(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create test session
	token := CreateTestSession(t, testCtx.Config.IBKRAccountID)
	defer DeleteTestSession(t, token)

	// Add account ID to context
	ctx = middleware.SetAccountIDInContext(ctx, testCtx.Config.IBKRAccountID)

	// Create market data service handler
	handler := api.NewMarketDataServiceHandler(testCtx.IBKRClient)

	// Create get quote request with invalid symbol
	req := connect.NewRequest(&marketdatav1.GetQuoteRequest{
		Symbol: "INVALID_SYMBOL_12345",
	})

	// Get quote (should fail or handle gracefully)
	_, err := handler.GetQuote(ctx, req)

	// We expect this might fail with the mock, but it shouldn't panic
	if err != nil {
		t.Logf("Expected error for invalid symbol: %v", err)
	}
}
