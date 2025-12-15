package ibkr

import "context"

// ClientInterface defines the interface for IBKR client operations
// This allows for easy mocking in tests
type ClientInterface interface {
	Ping(ctx context.Context) error
	AuthStatus(ctx context.Context) (*AuthStatusResponse, error)
	Reauthenticate(ctx context.Context) error
	GetAccounts(ctx context.Context) ([]Account, error)
	GetPortfolio(ctx context.Context) ([]Position, error)
	GetAccountSummary(ctx context.Context) (*AccountSummary, error)
	PlaceOrder(ctx context.Context, order *OrderRequest) (*OrderResponse, error)
	ModifyOrder(ctx context.Context, orderID string, order *OrderRequest) (*OrderResponse, error)
	CancelOrder(ctx context.Context, orderID string) error
	GetOrder(ctx context.Context, orderID string) (*Order, error)
	ListOrders(ctx context.Context) ([]Order, error)
	SearchContracts(ctx context.Context, symbol string) ([]Contract, error)
	GetMarketData(ctx context.Context, conIDs []int, fields []string) ([]MarketDataSnapshot, error)
	GetHistoricalData(ctx context.Context, conID int, period, barSize string) (*HistoricalData, error)
}

// Ensure Client implements ClientInterface
var _ ClientInterface = (*Client)(nil)
