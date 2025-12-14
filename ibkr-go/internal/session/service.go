package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/crypto"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/database"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/db"
)

const (
	// SessionTokenLength is the length of the session token in bytes.
	SessionTokenLength = 32
	// DefaultSessionTTL is the default session time-to-live.
	DefaultSessionTTL = 24 * time.Hour
)

// Service handles session management operations.
type Service struct {
	db            *database.DB
	encryptionKey []byte
	sessionTTL    time.Duration
}

// NewService creates a new session service.
func NewService(db *database.DB, encryptionKey []byte, sessionTTL time.Duration) *Service {
	if sessionTTL == 0 {
		sessionTTL = DefaultSessionTTL
	}

	return &Service{
		db:            db,
		encryptionKey: encryptionKey,
		sessionTTL:    sessionTTL,
	}
}

// Create creates a new session for the given account ID.
func (s *Service) Create(ctx context.Context, accountID string) (string, error) {
	// Generate random session token.
	tokenBytes := make([]byte, SessionTokenLength)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// Encrypt the token for storage.
	encryptedToken, err := crypto.EncryptToken([]byte(token), s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt token: %w", err)
	}

	// Hash the token for lookup.
	tokenHash := crypto.HashToken(token)

	// Calculate expiration time.
	expiresAt := time.Now().Add(s.sessionTTL)

	// Store in database.
	params := db.CreateSessionParams{
		AccountID:             accountID,
		SessionTokenEncrypted: encryptedToken,
		SessionTokenHash:      tokenHash,
		ExpiresAt: pgtype.Timestamp{
			Time:  expiresAt,
			Valid: true,
		},
	}

	session, err := s.db.Queries.CreateSession(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	_ = session // Session ID is in the database, but we return the token to the client

	return token, nil
}

// Validate validates a session token and returns the account ID if valid.
func (s *Service) Validate(ctx context.Context, token string) (string, error) {
	// Hash the token for lookup.
	tokenHash := crypto.HashToken(token)

	// Get session from database (query already filters expired sessions).
	session, err := s.db.Queries.GetSessionByHash(ctx, tokenHash)
	if err != nil {
		return "", fmt.Errorf("session not found: %w", err)
	}

	// Decrypt and verify the token.
	decryptedToken, err := crypto.DecryptToken(session.SessionTokenEncrypted, s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt token: %w", err)
	}

	// Verify the decrypted token matches the provided token.
	if string(decryptedToken) != token {
		return "", fmt.Errorf("token mismatch")
	}

	return session.AccountID, nil
}

// Delete deletes a session by token.
func (s *Service) Delete(ctx context.Context, token string) error {
	// Hash the token for lookup.
	tokenHash := crypto.HashToken(token)

	// Delete the session.
	if err := s.db.Queries.DeleteSessionByHash(ctx, tokenHash); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// CleanupExpired deletes all expired sessions.
func (s *Service) CleanupExpired(ctx context.Context) error {
	if err := s.db.Queries.DeleteExpiredSessions(ctx); err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	return nil
}
