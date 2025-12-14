package middleware

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/session"
)

const (
	// BearerTokenParts is the expected number of parts in a bearer token header.
	BearerTokenParts = 2
)

// SessionContextKey is the context key for storing the account ID from the session.
type SessionContextKey struct{}

// NewSessionInterceptor creates a ConnectRPC interceptor that validates session tokens.
func NewSessionInterceptor(sessionService *session.Service) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Extract token from Authorization header.
			token := extractToken(req.Header().Get("Authorization"))
			if token == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, nil)
			}

			// Validate the session.
			accountID, err := sessionService.Validate(ctx, token)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}

			// Add account ID to context.
			ctx = context.WithValue(ctx, SessionContextKey{}, accountID)

			return next(ctx, req)
		}
	}
}

// extractToken extracts the bearer token from the Authorization header.
func extractToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	// Expected format: "Bearer <token>".
	parts := strings.SplitN(authHeader, " ", BearerTokenParts)
	if len(parts) != BearerTokenParts || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// GetAccountIDFromContext retrieves the account ID from the context.
func GetAccountIDFromContext(ctx context.Context) (string, bool) {
	accountID, ok := ctx.Value(SessionContextKey{}).(string)

	return accountID, ok
}
