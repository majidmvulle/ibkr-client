package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/config"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/database"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/session"
)

// TestContext holds shared test resources.
type TestContext struct {
	DB             *database.DB
	IBKRClient     *ibkr.Client
	SessionService *session.Service
	Config         *config.Config
}

var testCtx *TestContext

// TestMain sets up and tears down test resources.
func TestMain(m *testing.M) {
	// Skip integration tests when running with -short flag
	if testing.Short() {
		fmt.Println("Skipping integration tests in short mode")
		os.Exit(0)
	}

	ctx := context.Background()

	// Load test configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load test config: %v\n", err)
		os.Exit(1)
	}

	// Initialize database
	db, err := database.New(ctx, cfg.DBWriteDSN, cfg.DBReadDSN)
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize IBKR client
	ibkrClient := ibkr.NewClient(cfg.IBKRGatewayURL, cfg.IBKRAccountID)

	// Initialize session service
	sessionService := session.NewService(db, cfg.EncryptionKey, 24*time.Hour)

	// Create test context
	testCtx = &TestContext{
		DB:             db,
		IBKRClient:     ibkrClient,
		SessionService: sessionService,
		Config:         cfg,
	}

	// Run tests
	code := m.Run()

	// Cleanup
	cleanupTestData(ctx, db)

	os.Exit(code)
}

// loadTestConfig loads configuration for integration tests.
func loadTestConfig() (*config.Config, error) {
	// Override with test values if needed
	if os.Getenv("APP_ENV") == "" {
		err := os.Setenv("APP_ENV", "test")
		if err != nil {
			return nil, err
		}
	}

	return config.Load()
}

// cleanupTestData removes test data from database.
func cleanupTestData(ctx context.Context, db *database.DB) {
	// Clean up test sessions
	queries := db.Queries
	_ = queries.DeleteExpiredSessions(ctx)
}

// CreateTestSession creates a test session and returns the token.
func CreateTestSession(t *testing.T, accountID string) string {
	t.Helper()

	ctx := context.Background()
	token, err := testCtx.SessionService.Create(ctx, accountID)
	if err != nil {
		t.Fatalf("Failed to create test session: %v", err)
	}

	return token
}

// ValidateTestSession validates a test session token.
func ValidateTestSession(t *testing.T, token string) string {
	t.Helper()

	ctx := context.Background()
	accountID, err := testCtx.SessionService.Validate(ctx, token)
	if err != nil {
		t.Fatalf("Failed to validate test session: %v", err)
	}

	return accountID
}

// DeleteTestSession deletes a test session.
func DeleteTestSession(t *testing.T, token string) {
	t.Helper()

	ctx := context.Background()
	err := testCtx.SessionService.Delete(ctx, token)
	if err != nil {
		t.Fatalf("Failed to delete test session: %v", err)
	}
}
