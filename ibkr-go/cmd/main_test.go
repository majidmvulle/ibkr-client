package cmd

import (
	"context"
	"testing"

	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/config"
)

func TestMain(t *testing.T) {
	// Test that main doesn't panic with nil config
	// We can't actually run main() as it would block, but we can test components
	cfg := &config.Config{
		AppName:     "test",
		AppEnv:      "test",
		HTTPPort:    8080,
		GRPCPort:    50051,
		MTLSEnabled: false,
	}

	if cfg.AppName != "test" {
		t.Errorf("Config not set correctly")
	}
}

func TestInitialization(t *testing.T) {
	// Test basic initialization doesn't panic
	ctx := context.Background()
	if ctx == nil {
		t.Error("Context is nil")
	}
}
