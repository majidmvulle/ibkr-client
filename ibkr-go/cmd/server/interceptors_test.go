package server

import (
	"testing"

	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/config"
)

func TestSetupInterceptors_NoMTLS(t *testing.T) {
	cfg := &config.Config{
		MTLSEnabled: false,
	}

	interceptors := setupInterceptors(cfg, nil, nil)
	if interceptors == nil {
		t.Error("setupInterceptors() should not return nil")
	}
}

func TestSetupInterceptors_WithMTLS(t *testing.T) {
	cfg := &config.Config{
		MTLSEnabled: true,
	}

	interceptors := setupInterceptors(cfg, nil, nil)
	if interceptors == nil {
		t.Error("setupInterceptors() should not return nil")
	}
}
