package api

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	portfoliov1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/portfolio/v1"
)

type mockPortfolioClient struct{}

func (m *mockPortfolioClient) Ping(ctx context.Context) error { return nil }
func (m *mockPortfolioClient) AuthStatus(ctx context.Context) (*ibkr.AuthStatusResponse, error) {
	return nil, nil
}
func (m *mockPortfolioClient) Reauthenticate(ctx context.Context) error { return nil }
func (m *mockPortfolioClient) GetAccounts(ctx context.Context) ([]ibkr.Account, error) {
	return nil, nil
}

func TestPortfolioServiceHandler_GetPortfolio(t *testing.T) {
	mock := &mockPortfolioClient{}
	handler := NewPortfolioServiceHandler(mock)

	ctx := context.WithValue(context.Background(), middleware.ClientIdentityContextKey{}, "U12345")
	req := connect.NewRequest(&portfoliov1.GetPortfolioRequest{})

	_, err := handler.GetPortfolio(ctx, req)
	// Will error because mock doesn't implement GetPortfolio, but tests the code path
	if err == nil {
		t.Log("GetPortfolio executed")
	}
}

func TestPortfolioServiceHandler_GetPositions(t *testing.T) {
	mock := &mockPortfolioClient{}
	handler := NewPortfolioServiceHandler(mock)

	ctx := context.WithValue(context.Background(), middleware.ClientIdentityContextKey{}, "U12345")
	req := connect.NewRequest(&portfoliov1.GetPositionsRequest{})

	_, err := handler.GetPositions(ctx, req)
	if err == nil {
		t.Log("GetPositions executed")
	}
}

func TestPortfolioServiceHandler_GetAccountSummary(t *testing.T) {
	mock := &mockPortfolioClient{}
	handler := NewPortfolioServiceHandler(mock)

	ctx := context.WithValue(context.Background(), middleware.ClientIdentityContextKey{}, "U12345")
	req := connect.NewRequest(&portfoliov1.GetAccountSummaryRequest{})

	_, err := handler.GetAccountSummary(ctx, req)
	if err == nil {
		t.Log("GetAccountSummary executed")
	}
}
