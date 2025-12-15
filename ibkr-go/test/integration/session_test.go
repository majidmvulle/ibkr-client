//go:build integration

package integration

import (
	"context"
	"testing"
	"time"
)

func TestIntegration_Session_CreateAndValidate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	accountID := testCtx.Config.IBKRAccountID

	// Create session
	token, err := testCtx.SessionService.Create(ctx, accountID)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if token == "" {
		t.Fatal("Expected token to be non-empty")
	}

	t.Logf("Session created: token length=%d", len(token))

	// Validate session
	validatedAccountID, err := testCtx.SessionService.Validate(ctx, token)
	if err != nil {
		t.Fatalf("Failed to validate session: %v", err)
	}

	if validatedAccountID != accountID {
		t.Errorf("Expected account ID %s, got %s", accountID, validatedAccountID)
	}

	// Cleanup
	err = testCtx.SessionService.Delete(ctx, token)
	if err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}
}

func TestIntegration_Session_InvalidToken(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Try to validate invalid token
	_, err := testCtx.SessionService.Validate(ctx, "invalid_token_12345")
	if err == nil {
		t.Error("Expected error for invalid token")
	}

	t.Logf("Invalid token correctly rejected: %v", err)
}

func TestIntegration_Session_DeleteAndValidate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	accountID := testCtx.Config.IBKRAccountID

	// Create session
	token, err := testCtx.SessionService.Create(ctx, accountID)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Delete session
	err = testCtx.SessionService.Delete(ctx, token)
	if err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}

	// Try to validate deleted session
	_, err = testCtx.SessionService.Validate(ctx, token)
	if err == nil {
		t.Error("Expected error for deleted session")
	}

	t.Logf("Deleted session correctly rejected: %v", err)
}

func TestIntegration_Session_ExpiredCleanup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	accountID := testCtx.Config.IBKRAccountID

	// Create session
	token, err := testCtx.SessionService.Create(ctx, accountID)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}
	defer testCtx.SessionService.Delete(ctx, token)

	// Cleanup expired sessions (this one shouldn't be deleted as it's not expired)
	err = testCtx.SessionService.CleanupExpired(ctx)
	if err != nil {
		t.Fatalf("Failed to cleanup expired sessions: %v", err)
	}

	// Validate session still exists
	_, err = testCtx.SessionService.Validate(ctx, token)
	if err != nil {
		t.Errorf("Session should still be valid after cleanup: %v", err)
	}
}

func TestIntegration_Session_MultipleTokens(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	accountID := testCtx.Config.IBKRAccountID

	// Create multiple sessions
	tokens := make([]string, 3)
	for i := 0; i < 3; i++ {
		token, err := testCtx.SessionService.Create(ctx, accountID)
		if err != nil {
			t.Fatalf("Failed to create session %d: %v", i, err)
		}
		tokens[i] = token
	}

	// Validate all sessions
	for i, token := range tokens {
		validatedAccountID, err := testCtx.SessionService.Validate(ctx, token)
		if err != nil {
			t.Errorf("Failed to validate session %d: %v", i, err)
		}
		if validatedAccountID != accountID {
			t.Errorf("Session %d: expected account ID %s, got %s", i, accountID, validatedAccountID)
		}
	}

	// Cleanup all sessions
	for i, token := range tokens {
		err := testCtx.SessionService.Delete(ctx, token)
		if err != nil {
			t.Errorf("Failed to delete session %d: %v", i, err)
		}
	}
}

func TestIntegration_Session_TokenUniqueness(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	accountID := testCtx.Config.IBKRAccountID

	// Create two sessions
	token1, err := testCtx.SessionService.Create(ctx, accountID)
	if err != nil {
		t.Fatalf("Failed to create session 1: %v", err)
	}
	defer testCtx.SessionService.Delete(ctx, token1)

	// Small delay to ensure different timestamps
	time.Sleep(10 * time.Millisecond)

	token2, err := testCtx.SessionService.Create(ctx, accountID)
	if err != nil {
		t.Fatalf("Failed to create session 2: %v", err)
	}
	defer testCtx.SessionService.Delete(ctx, token2)

	// Tokens should be different
	if token1 == token2 {
		t.Error("Expected tokens to be unique")
	}

	t.Logf("Tokens are unique: token1 length=%d, token2 length=%d", len(token1), len(token2))
}
