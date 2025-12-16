package api

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	portfoliov1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/portfolio/v1"
)

func TestGetPortfolio(t *testing.T) {
	mockClient := new(MockPortfolioClient)
	handler := NewPortfolioServiceHandler(mockClient)

	ctx := middleware.SetAccountIDInContext(context.Background(), "U12345")
	req := connect.NewRequest(&portfoliov1.GetPortfolioRequest{AccountId: "U12345"})

	// Mocks
	positions := []ibkr.Position{
		{ContractDesc: "AAPL", Position: 10, MktValue: 1500, Currency: "USD"},
	}
	mockClient.On("GetPortfolio", ctx).Return(positions, nil)

	summary := &ibkr.AccountSummary{
		NetLiquidation: 10000,
		TotalCashValue: 5000,
		Currency:       "USD",
	}
	mockClient.On("GetAccountSummary", ctx).Return(summary, nil)

	resp, err := handler.GetPortfolio(ctx, req)
	if err != nil {
		t.Fatalf("GetPortfolio() error = %v", err)
	}

	if len(resp.Msg.Portfolio.Positions) != 1 {
		t.Errorf("Expected 1 position, got %d", len(resp.Msg.Portfolio.Positions))
	}
}

func TestGetAccountSummary(t *testing.T) {
	mockClient := new(MockPortfolioClient)
	handler := NewPortfolioServiceHandler(mockClient)

	ctx := middleware.SetAccountIDInContext(context.Background(), "U12345")
	req := connect.NewRequest(&portfoliov1.GetAccountSummaryRequest{})

	summary := &ibkr.AccountSummary{
		NetLiquidation: 10000,
		Currency:       "USD",
	}
	mockClient.On("GetAccountSummary", ctx).Return(summary, nil)

	resp, err := handler.GetAccountSummary(ctx, req)
	if err != nil {
		t.Fatalf("GetAccountSummary() error = %v", err)
	}

	if resp.Msg.AccountSummary.NetLiquidation.Units != 10000 {
		t.Errorf("NetLiquidation units = %v, want 10000", resp.Msg.AccountSummary.NetLiquidation.Units)
	}
}
