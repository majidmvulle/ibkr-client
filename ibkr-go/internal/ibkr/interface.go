package ibkr

import "context"

// IBKRClient defines the interface for interacting with the IBKR Gateway API.
type IBKRClient interface {
	// Basic operations.
	Ping(ctx context.Context) error
	AuthStatus(ctx context.Context) (*AuthStatusResponse, error)
	Reauthenticate(ctx context.Context) error
	GetAccounts(ctx context.Context) ([]Account, error)

	// Market data operations.
	GetMarketData(ctx context.Context, conIDs []int, fields []string) ([]MarketDataSnapshot, error)
	GetHistoricalData(ctx context.Context, conID int, period, barSize string) (*HistoricalDataResponse, error)
	SearchContracts(ctx context.Context, symbol string) ([]Contract, error)

	// Order operations.
	PlaceOrder(ctx context.Context, req *PlaceOrderRequest) (*OrderResponse, error)
	ModifyOrder(ctx context.Context, orderID string, req *ModifyOrderRequest) (*OrderResponse, error)
	CancelOrder(ctx context.Context, orderID string) error
	GetLiveOrders(ctx context.Context) ([]Order, error)

	// Portfolio operations.
	GetPortfolio(ctx context.Context) ([]Position, error)
	GetAccountSummary(ctx context.Context) (*AccountSummary, error)
}
