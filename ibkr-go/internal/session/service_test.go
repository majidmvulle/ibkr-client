package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/crypto"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/db"
	"github.com/stretchr/testify/assert"
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

func TestNewService(t *testing.T) {
	tests := []struct {
		name       string
		sessionTTL time.Duration
		wantTTL    time.Duration
	}{
		{
			name:       "with custom TTL",
			sessionTTL: 1 * time.Hour,
			wantTTL:    1 * time.Hour,
		},
		{
			name:       "with zero TTL uses default",
			sessionTTL: 0,
			wantTTL:    DefaultSessionTTL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQuerier := new(MockQuerier)
			service := NewService(mockQuerier, []byte("test-key"), tt.sessionTTL)
			if service == nil {
				t.Fatal("NewService() returned nil")
			}
			if service.sessionTTL != tt.wantTTL {
				t.Errorf("NewService() sessionTTL = %v, want %v", service.sessionTTL, tt.wantTTL)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	mockQuerier := new(MockQuerier)
	key := []byte("12345678901234567890123456789012") // 32 bytes key
	service := NewService(mockQuerier, key, time.Hour)
	ctx := context.Background()
	accountID := "acc-123"

	mockQuerier.On("CreateSession", ctx, mock.AnythingOfType("db.CreateSessionParams")).Return(db.Session{}, nil)

	token, err := service.Create(ctx, accountID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockQuerier.AssertExpectations(t)
}

func TestService_Create_Error(t *testing.T) {
	mockQuerier := new(MockQuerier)
	key := []byte("12345678901234567890123456789012")
	service := NewService(mockQuerier, key, time.Hour)
	ctx := context.Background()
	accountID := "acc-123"

	mockQuerier.On("CreateSession", ctx, mock.AnythingOfType("db.CreateSessionParams")).Return(db.Session{}, errors.New("db error"))

	token, err := service.Create(ctx, accountID)
	assert.Error(t, err)
	assert.Empty(t, token)

	mockQuerier.AssertExpectations(t)
}

func TestService_Validate_Success(t *testing.T) {
	mockQuerier := new(MockQuerier)
	key := []byte("12345678901234567890123456789012")
	service := NewService(mockQuerier, key, time.Hour)
	ctx := context.Background()

	// Pre-calculate token and hash
	token := "some-valid-token"
	tokenHash := crypto.HashToken(token)
	encryptedToken, _ := crypto.EncryptToken([]byte(token), key)

	session := db.Session{
		AccountID:             "acc-123",
		SessionTokenEncrypted: encryptedToken,
	}

	mockQuerier.On("GetSessionByHash", ctx, tokenHash).Return(session, nil)

	accountID, err := service.Validate(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, "acc-123", accountID)

	mockQuerier.AssertExpectations(t)
}

func TestService_Validate_NotFound(t *testing.T) {
	mockQuerier := new(MockQuerier)
	key := []byte("12345678901234567890123456789012")
	service := NewService(mockQuerier, key, time.Hour)
	ctx := context.Background()
	token := "invalid-token"
	tokenHash := crypto.HashToken(token)

	mockQuerier.On("GetSessionByHash", ctx, tokenHash).Return(db.Session{}, errors.New("not found"))

	accountID, err := service.Validate(ctx, token)
	assert.Error(t, err)
	assert.Empty(t, accountID)

	mockQuerier.AssertExpectations(t)
}

func TestService_Validate_TokenMismatch(t *testing.T) {
	mockQuerier := new(MockQuerier)
	key := []byte("12345678901234567890123456789012")
	service := NewService(mockQuerier, key, time.Hour)
	ctx := context.Background()

	token := "token-a"
	tokenHash := crypto.HashToken(token)

	// Encrypt a DIFFERENT token
	encryptedToken, _ := crypto.EncryptToken([]byte("token-b"), key)

	session := db.Session{
		AccountID:             "acc-123",
		SessionTokenEncrypted: encryptedToken,
	}

	mockQuerier.On("GetSessionByHash", ctx, tokenHash).Return(session, nil)

	accountID, err := service.Validate(ctx, token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token mismatch")
	assert.Empty(t, accountID)
}

func TestService_Delete(t *testing.T) {
	mockQuerier := new(MockQuerier)
	service := NewService(mockQuerier, nil, time.Hour)
	ctx := context.Background()
	token := "token"
	tokenHash := crypto.HashToken(token)

	mockQuerier.On("DeleteSessionByHash", ctx, tokenHash).Return(nil)

	err := service.Delete(ctx, token)
	assert.NoError(t, err)

	mockQuerier.AssertExpectations(t)
}

func TestService_Delete_Error(t *testing.T) {
	mockQuerier := new(MockQuerier)
	service := NewService(mockQuerier, nil, time.Hour)
	ctx := context.Background()
	token := "token"
	tokenHash := crypto.HashToken(token)

	mockQuerier.On("DeleteSessionByHash", ctx, tokenHash).Return(errors.New("db error"))

	err := service.Delete(ctx, token)
	assert.Error(t, err)
}

func TestService_CleanupExpired(t *testing.T) {
	mockQuerier := new(MockQuerier)
	service := NewService(mockQuerier, nil, time.Hour)
	ctx := context.Background()

	mockQuerier.On("DeleteExpiredSessions", ctx).Return(nil)

	err := service.CleanupExpired(ctx)
	assert.NoError(t, err)

	mockQuerier.AssertExpectations(t)
}

func TestService_CleanupExpired_Error(t *testing.T) {
	mockQuerier := new(MockQuerier)
	service := NewService(mockQuerier, nil, time.Hour)
	ctx := context.Background()

	mockQuerier.On("DeleteExpiredSessions", ctx).Return(errors.New("db error"))

	err := service.CleanupExpired(ctx)
	assert.Error(t, err)
}
