package api

import (
	"context"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	"github.com/stretchr/testify/mock"
)

type MockOrderClient struct {
	mock.Mock
}

func (m *MockOrderClient) PlaceOrder(ctx context.Context, req *ibkr.PlaceOrderRequest) (*ibkr.OrderResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ibkr.OrderResponse), args.Error(1)
}

func (m *MockOrderClient) ModifyOrder(ctx context.Context, orderID string, req *ibkr.ModifyOrderRequest) (*ibkr.OrderResponse, error) {
	args := m.Called(ctx, orderID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ibkr.OrderResponse), args.Error(1)
}

func (m *MockOrderClient) CancelOrder(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

func (m *MockOrderClient) GetLiveOrders(ctx context.Context) ([]ibkr.Order, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ibkr.Order), args.Error(1)
}

type MockMarketDataClient struct {
	mock.Mock
}

func (m *MockMarketDataClient) GetMarketData(ctx context.Context, conIDs []int, fields []string) ([]ibkr.MarketDataSnapshot, error) {
	args := m.Called(ctx, conIDs, fields)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ibkr.MarketDataSnapshot), args.Error(1)
}

func (m *MockMarketDataClient) GetHistoricalData(ctx context.Context, conID int, period, barSize string) (*ibkr.HistoricalDataResponse, error) {
	args := m.Called(ctx, conID, period, barSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ibkr.HistoricalDataResponse), args.Error(1)
}

func (m *MockMarketDataClient) SearchContracts(ctx context.Context, symbol string) ([]ibkr.Contract, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ibkr.Contract), args.Error(1)
}

type MockPortfolioClient struct {
	mock.Mock
}

func (m *MockPortfolioClient) GetPortfolio(ctx context.Context) ([]ibkr.Position, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ibkr.Position), args.Error(1)
}

func (m *MockPortfolioClient) GetAccountSummary(ctx context.Context) (*ibkr.AccountSummary, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ibkr.AccountSummary), args.Error(1)
}
