package api

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	marketdatav1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/marketdata/v1"
)

type mockMarketDataClient struct{}

func (m *mockMarketDataClient) Ping(ctx context.Context) error { return nil }
func (m *mockMarketDataClient) AuthStatus(ctx context.Context) (*ibkr.AuthStatusResponse, error) {
	return nil, nil
}
func (m *mockMarketDataClient) Reauthenticate(ctx context.Context) error { return nil }
func (m *mockMarketDataClient) GetAccounts(ctx context.Context) ([]ibkr.Account, error) {
	return nil, nil
}

func TestMarketDataServiceHandler_GetQuote(t *testing.T) {
	mock := &mockMarketDataClient{}
	handler := NewMarketDataServiceHandler(mock)

	ctx := context.WithValue(context.Background(), middleware.ClientIdentityContextKey{}, "U12345")
	req := connect.NewRequest(&marketdatav1.GetQuoteRequest{
		Symbol: "AAPL",
	})

	_, err := handler.GetQuote(ctx, req)
	if err == nil {
		t.Log("GetQuote executed")
	}
}

func TestMarketDataServiceHandler_GetHistoricalData(t *testing.T) {
	mock := &mockMarketDataClient{}
	handler := NewMarketDataServiceHandler(mock)

	ctx := context.WithValue(context.Background(), middleware.ClientIdentityContextKey{}, "U12345")
	req := connect.NewRequest(&marketdatav1.GetHistoricalDataRequest{
		Symbol:  "AAPL",
		Period:  "1d",
		BarSize: marketdatav1.BarSize_BAR_SIZE_1_DAY,
	})

	_, err := handler.GetHistoricalData(ctx, req)
	if err == nil {
		t.Log("GetHistoricalData executed")
	}
}
