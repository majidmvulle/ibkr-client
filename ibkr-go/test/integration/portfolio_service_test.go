//go:build integration

package integration

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/api"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	portfoliov1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/portfolio/v1"
)

func TestIntegration_PortfolioService_GetPortfolio(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create test session
	token := CreateTestSession(t, testCtx.Config.IBKRAccountID)
	defer DeleteTestSession(t, token)

	// Add account ID to context
	ctx = middleware.SetAccountIDInContext(ctx, testCtx.Config.IBKRAccountID)

	// Create portfolio service handler
	handler := api.NewPortfolioServiceHandler(testCtx.IBKRClient)

	// Create get portfolio request
	req := connect.NewRequest(&portfoliov1.GetPortfolioRequest{})

	// Get portfolio
	resp, err := handler.GetPortfolio(ctx, req)
	if err != nil {
		t.Fatalf("GetPortfolio failed: %v", err)
	}

	// Verify response
	if resp.Msg.Portfolio == nil {
		t.Fatal("Expected portfolio to be set")
	}

	portfolio := resp.Msg.Portfolio

	if portfolio.AccountId != testCtx.Config.IBKRAccountID {
		t.Errorf("Expected account ID %s, got %s", testCtx.Config.IBKRAccountID, portfolio.AccountId)
	}

	if portfolio.TotalValue == nil {
		t.Error("Expected total value to be set")
	}

	if portfolio.CashBalance == nil {
		t.Error("Expected cash balance to be set")
	}

	t.Logf("Portfolio retrieved: Account=%s, Positions=%d",
		portfolio.AccountId, len(portfolio.Positions))
}

func TestIntegration_PortfolioService_GetPositions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create test session
	token := CreateTestSession(t, testCtx.Config.IBKRAccountID)
	defer DeleteTestSession(t, token)

	// Add account ID to context
	ctx = middleware.SetAccountIDInContext(ctx, testCtx.Config.IBKRAccountID)

	// Create portfolio service handler
	handler := api.NewPortfolioServiceHandler(testCtx.IBKRClient)

	// Create get positions request
	req := connect.NewRequest(&portfoliov1.GetPositionsRequest{})

	// Get positions
	resp, err := handler.GetPositions(ctx, req)
	if err != nil {
		t.Fatalf("GetPositions failed: %v", err)
	}

	// Verify response
	if resp.Msg.Positions == nil {
		t.Error("Expected positions list to be initialized")
	}

	t.Logf("Retrieved %d positions", len(resp.Msg.Positions))

	// Verify position structure if any exist
	for _, pos := range resp.Msg.Positions {
		if pos.Symbol == "" {
			t.Error("Expected position symbol to be set")
		}
		if pos.Quantity == 0 {
			t.Error("Expected position quantity to be non-zero")
		}
		if pos.MarketValue == nil {
			t.Error("Expected market value to be set")
		}
	}
}

func TestIntegration_PortfolioService_GetAccountSummary(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create test session
	token := CreateTestSession(t, testCtx.Config.IBKRAccountID)
	defer DeleteTestSession(t, token)

	// Add account ID to context
	ctx = middleware.SetAccountIDInContext(ctx, testCtx.Config.IBKRAccountID)

	// Create portfolio service handler
	handler := api.NewPortfolioServiceHandler(testCtx.IBKRClient)

	// Create get account summary request
	req := connect.NewRequest(&portfoliov1.GetAccountSummaryRequest{})

	// Get account summary
	resp, err := handler.GetAccountSummary(ctx, req)
	if err != nil {
		t.Fatalf("GetAccountSummary failed: %v", err)
	}

	// Verify response
	if resp.Msg.AccountSummary == nil {
		t.Fatal("Expected account_summary to be set")
	}

	summary := resp.Msg.AccountSummary

	if summary.AccountId != testCtx.Config.IBKRAccountID {
		t.Errorf("Expected account ID %s, got %s", testCtx.Config.IBKRAccountID, summary.AccountId)
	}

	if summary.NetLiquidation == nil {
		t.Error("Expected net liquidation to be set")
	}

	if summary.TotalCash == nil {
		t.Error("Expected total cash to be set")
	}

	t.Logf("Account summary retrieved: Account=%s", summary.AccountId)
}

func TestIntegration_PortfolioService_MoneyConversion(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create test session
	token := CreateTestSession(t, testCtx.Config.IBKRAccountID)
	defer DeleteTestSession(t, token)

	// Add account ID to context
	ctx = middleware.SetAccountIDInContext(ctx, testCtx.Config.IBKRAccountID)

	// Create portfolio service handler
	handler := api.NewPortfolioServiceHandler(testCtx.IBKRClient)

	// Get portfolio to test money conversion
	req := connect.NewRequest(&portfoliov1.GetPortfolioRequest{})
	resp, err := handler.GetPortfolio(ctx, req)
	if err != nil {
		t.Fatalf("GetPortfolio failed: %v", err)
	}

	portfolio := resp.Msg.Portfolio

	// Verify Money proto structure
	if portfolio.TotalValue != nil {
		if portfolio.TotalValue.CurrencyCode == "" {
			t.Error("Expected currency code to be set")
		}
		// Units and Nanos should be set (can be 0)
		t.Logf("Total Value: %d.%09d %s",
			portfolio.TotalValue.Units,
			portfolio.TotalValue.Nanos,
			portfolio.TotalValue.CurrencyCode)
	}

	if portfolio.CashBalance != nil {
		if portfolio.CashBalance.CurrencyCode == "" {
			t.Error("Expected currency code to be set")
		}
		t.Logf("Cash Balance: %d.%09d %s",
			portfolio.CashBalance.Units,
			portfolio.CashBalance.Nanos,
			portfolio.CashBalance.CurrencyCode)
	}
}
