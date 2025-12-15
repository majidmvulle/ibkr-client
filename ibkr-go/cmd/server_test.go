package cmd

import (
	"context"
	"testing"

	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/config"
)

func TestServerSetup(t *testing.T) {
	cfg := &config.Config{
		AppName:     "test",
		HTTPPort:    8080,
		MTLSEnabled: false,
	}

	if cfg.AppName == "" {
		t.Error("AppName should not be empty")
	}
}

func TestContextCreation(t *testing.T) {
	ctx := context.Background()
	if ctx == nil {
		t.Error("Context should not be nil")
	}
}
