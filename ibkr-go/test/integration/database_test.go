//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/db"
)

func TestIntegration_Database_Connection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Test pool health
	err := testCtx.DB.Pool.Ping(ctx)
	if err != nil {
		t.Fatalf("Pool ping failed: %v", err)
	}

	t.Log("Database connection healthy")
}

func TestIntegration_Database_SessionQueries(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	queries := testCtx.DB.Queries

	// Create a test session directly via queries
	accountID := testCtx.Config.IBKRAccountID
	tokenHash := "test_hash_12345"
	encryptedToken := []byte("encrypted_test_token")
	expiresAt := pgtype.Timestamp{
		Time:  time.Now().Add(24 * time.Hour),
		Valid: true,
	}

	session, err := queries.CreateSession(ctx, db.CreateSessionParams{
		AccountID:             accountID,
		SessionTokenHash:      tokenHash,
		SessionTokenEncrypted: encryptedToken,
		ExpiresAt:             expiresAt,
	})
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session.AccountID != accountID {
		t.Errorf("Expected account ID %s, got %s", accountID, session.AccountID)
	}

	// Get session
	retrievedSession, err := queries.GetSessionByHash(ctx, tokenHash)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if retrievedSession.AccountID != accountID {
		t.Errorf("Expected account ID %s, got %s", accountID, retrievedSession.AccountID)
	}

	if string(retrievedSession.SessionTokenEncrypted) != string(encryptedToken) {
		t.Error("Encrypted token mismatch")
	}

	// Delete session
	err = queries.DeleteSessionByHash(ctx, tokenHash)
	if err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}

	// Verify deletion
	_, err = queries.GetSessionByHash(ctx, tokenHash)
	if err == nil {
		t.Error("Expected error when getting deleted session")
	}
}

func TestIntegration_Database_HealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Test health check
	err := testCtx.DB.Health(ctx)
	if err != nil {
		t.Errorf("Database health check failed: %v", err)
	}

	t.Log("Database health check passed")
}

func TestIntegration_Database_TransactionIsolation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	queries := testCtx.DB.Queries

	// Create two sessions with different tokens
	accountID := testCtx.Config.IBKRAccountID

	token1Hash := "test_hash_1"
	encryptedToken1 := []byte("encrypted_test_token_1")
	expiresAt := pgtype.Timestamp{
		Time:  time.Now().Add(24 * time.Hour),
		Valid: true,
	}

	token2Hash := "test_hash_2"
	encryptedToken2 := []byte("encrypted_test_token_2")

	// Create first session
	_, err := queries.CreateSession(ctx, db.CreateSessionParams{
		AccountID:             accountID,
		SessionTokenHash:      token1Hash,
		SessionTokenEncrypted: encryptedToken1,
		ExpiresAt:             expiresAt,
	})
	if err != nil {
		t.Fatalf("Failed to create session 1: %v", err)
	}
	defer queries.DeleteSessionByHash(ctx, token1Hash)

	// Create second session
	_, err = queries.CreateSession(ctx, db.CreateSessionParams{
		AccountID:             accountID,
		SessionTokenHash:      token2Hash,
		SessionTokenEncrypted: encryptedToken2,
		ExpiresAt:             expiresAt,
	})
	if err != nil {
		t.Fatalf("Failed to create session 2: %v", err)
	}
	defer queries.DeleteSessionByHash(ctx, token2Hash)

	// Verify both sessions exist independently
	session1, err := queries.GetSessionByHash(ctx, token1Hash)
	if err != nil {
		t.Fatalf("Failed to get session 1: %v", err)
	}

	session2, err := queries.GetSessionByHash(ctx, token2Hash)
	if err != nil {
		t.Fatalf("Failed to get session 2: %v", err)
	}

	if string(session1.SessionTokenEncrypted) == string(session2.SessionTokenEncrypted) {
		t.Error("Sessions should have different encrypted tokens")
	}
}

func TestIntegration_Database_ConcurrentAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Test concurrent pings
	done := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func() {
			err := testCtx.DB.Pool.Ping(ctx)
			done <- err
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		err := <-done
		if err != nil {
			t.Errorf("Concurrent ping %d failed: %v", i, err)
		}
	}

	t.Log("Concurrent database access successful")
}
