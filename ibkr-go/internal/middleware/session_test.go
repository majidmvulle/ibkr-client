package middleware

import (
	"context"
	"errors"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/crypto"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/db"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/session"
	"github.com/stretchr/testify/mock"
)

// MockQuerier is a mock implementation of db.Querier
type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Session), args.Error(1)
}

func (m *MockQuerier) DeleteExpiredSessions(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockQuerier) DeleteSessionByHash(ctx context.Context, sessionTokenHash string) error {
	args := m.Called(ctx, sessionTokenHash)
	return args.Error(0)
}

func (m *MockQuerier) GetSessionByHash(ctx context.Context, sessionTokenHash string) (db.Session, error) {
	args := m.Called(ctx, sessionTokenHash)
	return args.Get(0).(db.Session), args.Error(1)
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		want       string
	}{
		{
			name:       "valid bearer token",
			authHeader: "Bearer abc123xyz",
			want:       "abc123xyz",
		},
		{
			name:       "valid bearer token with uppercase",
			authHeader: "BEARER token123",
			want:       "token123",
		},
		{
			name:       "valid bearer token with mixed case",
			authHeader: "BeArEr mixedCase",
			want:       "mixedCase",
		},
		{
			name:       "empty header",
			authHeader: "",
			want:       "",
		},
		{
			name:       "missing token",
			authHeader: "Bearer",
			want:       "",
		},
		{
			name:       "missing bearer prefix",
			authHeader: "token123",
			want:       "",
		},
		{
			name:       "wrong prefix",
			authHeader: "Basic token123",
			want:       "",
		},
		{
			name:       "token with spaces",
			authHeader: "Bearer token with spaces",
			want:       "token with spaces",
		},
		{
			name:       "only spaces",
			authHeader: "   ",
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractToken(tt.authHeader)
			if got != tt.want {
				t.Errorf("extractToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAccountIDFromContext(t *testing.T) {
	tests := []struct {
		name      string
		ctx       context.Context
		wantID    string
		wantFound bool
	}{
		{
			name:      "account ID present",
			ctx:       context.WithValue(context.Background(), SessionContextKey{}, "DU123456"),
			wantID:    "DU123456",
			wantFound: true,
		},
		{
			name:      "account ID not present",
			ctx:       context.Background(),
			wantID:    "",
			wantFound: false,
		},
		{
			name:      "wrong type in context",
			ctx:       context.WithValue(context.Background(), SessionContextKey{}, 12345),
			wantID:    "",
			wantFound: false,
		},
		{
			name:      "empty account ID",
			ctx:       context.WithValue(context.Background(), SessionContextKey{}, ""),
			wantID:    "",
			wantFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotFound := GetAccountIDFromContext(tt.ctx)
			if gotID != tt.wantID {
				t.Errorf("GetAccountIDFromContext() ID = %v, want %v", gotID, tt.wantID)
			}
			if gotFound != tt.wantFound {
				t.Errorf("GetAccountIDFromContext() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestSetAccountIDInContext(t *testing.T) {
	tests := []struct {
		name      string
		accountID string
	}{
		{
			name:      "set valid account ID",
			accountID: "DU123456",
		},
		{
			name:      "set empty account ID",
			accountID: "",
		},
		{
			name:      "set long account ID",
			accountID: "VERY_LONG_ACCOUNT_ID_12345678901234567890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := SetAccountIDInContext(context.Background(), tt.accountID)

			gotID, found := GetAccountIDFromContext(ctx)
			if !found {
				t.Error("SetAccountIDInContext() did not set account ID in context")
			}
			if gotID != tt.accountID {
				t.Errorf("SetAccountIDInContext() = %v, want %v", gotID, tt.accountID)
			}
		})
	}
}

func TestSessionContextKey(t *testing.T) {
	// Verify that SessionContextKey is a unique type
	key1 := SessionContextKey{}
	key2 := SessionContextKey{}

	ctx1 := context.WithValue(context.Background(), key1, "value1")
	ctx2 := context.WithValue(ctx1, key2, "value2")

	// Both keys should refer to the same value since they're the same type
	val1, ok1 := ctx1.Value(key1).(string)
	val2, ok2 := ctx2.Value(key2).(string)

	if !ok1 || !ok2 {
		t.Error("Failed to retrieve values from context")
	}

	// The second WithValue should overwrite the first
	if val2 != "value2" {
		t.Errorf("Expected value2, got %v", val2)
	}

	// Verify the first context still has its original value
	if val1 != "value1" {
		t.Errorf("Expected value1, got %v", val1)
	}
}

func TestContextChaining(t *testing.T) {
	// Test that context values can be chained properly
	ctx := context.Background()
	ctx = SetAccountIDInContext(ctx, "ACCOUNT1")

	// Add another value to the context
	type otherKey struct{}
	ctx = context.WithValue(ctx, otherKey{}, "other value")

	// Verify both values are accessible
	accountID, found := GetAccountIDFromContext(ctx)
	if !found || accountID != "ACCOUNT1" {
		t.Errorf("Account ID not found or incorrect: %v", accountID)
	}

	otherVal, ok := ctx.Value(otherKey{}).(string)
	if !ok || otherVal != "other value" {
		t.Errorf("Other value not found or incorrect: %v", otherVal)
	}
}

func TestSessionInterceptor_NoAuthHeader(t *testing.T) {
	mockQuerier := new(MockQuerier)
	svc := session.NewService(mockQuerier, []byte("key"), time.Hour)
	interceptor := NewSessionInterceptor(svc)

	handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		return connect.NewResponse(&struct{}{}), nil
	}
	wrapped := interceptor(handler)

	req := connect.NewRequest(&struct{}{})
	// No header
	_, err := wrapped(context.Background(), req)
	if err == nil {
		t.Error("Expected error for missing auth header")
	}
	if connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Errorf("Expected Unauthenticated, got %v", connect.CodeOf(err))
	}
}

func TestSessionInterceptor_InvalidToken(t *testing.T) {
	mockQuerier := new(MockQuerier)
	svc := session.NewService(mockQuerier, []byte("12345678901234567890123456789012"), time.Hour)
	interceptor := NewSessionInterceptor(svc)

	handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		return connect.NewResponse(&struct{}{}), nil
	}
	wrapped := interceptor(handler)

	req := connect.NewRequest(&struct{}{})
	req.Header().Set("Authorization", "Bearer invalid-token")

	mockQuerier.On("GetSessionByHash", mock.Anything, mock.Anything).Return(db.Session{}, errors.New("not found"))

	_, err := wrapped(context.Background(), req)
	if err == nil {
		t.Error("Expected error for invalid token")
	}
	if connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Errorf("Expected Unauthenticated, got %v", connect.CodeOf(err))
	}
}

func TestSessionInterceptor_Success(t *testing.T) {
	mockQuerier := new(MockQuerier)
	key := []byte("12345678901234567890123456789012")
	svc := session.NewService(mockQuerier, key, time.Hour)
	interceptor := NewSessionInterceptor(svc)

	accountID := "acc-123"
	token := "valid-token"
	tokenHash := crypto.HashToken(token)
	encryptedToken, _ := crypto.EncryptToken([]byte(token), key)

	sessionData := db.Session{
		AccountID:             accountID,
		SessionTokenEncrypted: encryptedToken,
	}

	mockQuerier.On("GetSessionByHash", mock.Anything, tokenHash).Return(sessionData, nil)

	handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		// Verify account ID is in context
		id, ok := GetAccountIDFromContext(ctx)
		if !ok || id != accountID {
			return nil, errors.New("account ID not found in context")
		}
		return connect.NewResponse(&struct{}{}), nil
	}
	wrapped := interceptor(handler)

	req := connect.NewRequest(&struct{}{})
	req.Header().Set("Authorization", "Bearer "+token)

	_, err := wrapped(context.Background(), req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
