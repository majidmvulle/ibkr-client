package api

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	marketdatav1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/marketdata/v1"
)

func TestGetQuote(t *testing.T) {
	mockClient := new(MockMarketDataClient)
	handler := NewMarketDataServiceHandler(mockClient)

	ctx := middleware.SetAccountIDInContext(context.Background(), "U12345")
	req := connect.NewRequest(&marketdatav1.GetQuoteRequest{Symbol: "AAPL"})

	// Mocks
	contracts := []ibkr.Contract{{ConID: 12345, Symbol: "AAPL"}}
	mockClient.On("SearchContracts", ctx, "AAPL").Return(contracts, nil)

	snapshots := []ibkr.MarketDataSnapshot{{LastPrice: 150.0}}
	mockClient.On("GetMarketData", ctx, []int{12345}, []string(nil)).Return(snapshots, nil)

	resp, err := handler.GetQuote(ctx, req)
	if err != nil {
		t.Fatalf("GetQuote() error = %v", err)
	}

	if resp.Msg.Quote.Last != 150.0 {
		t.Errorf("Last = %v, want 150.0", resp.Msg.Quote.Last)
	}
}
