package ibkr

import "context"

// BasicClient defines basic operations.
type BasicClient interface {
	Ping(ctx context.Context) error
	AuthStatus(ctx context.Context) (*AuthStatusResponse, error)
	Reauthenticate(ctx context.Context) error
	GetAccounts(ctx context.Context) ([]Account, error)
}

// MarketDataClient defines market data operations.
type MarketDataClient interface {
	GetMarketData(ctx context.Context, conIDs []int, fields []string) ([]MarketDataSnapshot, error)
	GetHistoricalData(ctx context.Context, conID int, period, barSize string) (*HistoricalDataResponse, error)
	SearchContracts(ctx context.Context, symbol string) ([]Contract, error)
}

// OrderClient defines order operations.
type OrderClient interface {
	PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*OrderResponse, error)
	ModifyOrder(ctx context.Context, orderID string, req *ModifyOrderRequest) (*OrderResponse, error)
	CancelOrder(ctx context.Context, orderID string) error
	GetLiveOrders(ctx context.Context) ([]Order, error)
}

// PortfolioClient defines portfolio operations.
type PortfolioClient interface {
	GetPortfolio(ctx context.Context) ([]Position, error)
	GetAccountSummary(ctx context.Context) (*AccountSummary, error)
}

// IBKRClient defines the interface for interacting with the IBKR Gateway API.
type IBKRClient interface {
	BasicClient
	MarketDataClient
	OrderClient
	PortfolioClient
}
