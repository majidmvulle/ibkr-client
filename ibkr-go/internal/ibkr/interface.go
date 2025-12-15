package ibkr

import "context"

// ClientInterface defines the interface for IBKR client operations
// This allows for easy mocking in tests
type ClientInterface interface {
	Ping(ctx context.Context) error
	AuthStatus(ctx context.Context) (*AuthStatusResponse, error)
	Reauthenticate(ctx context.Context) error
	GetAccounts(ctx context.Context) ([]Account, error)
}

// Ensure Client implements ClientInterface
var _ ClientInterface = (*Client)(nil)
