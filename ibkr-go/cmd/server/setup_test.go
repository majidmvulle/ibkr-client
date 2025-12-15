package cmd

import (
	"testing"

	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/config"
)

func TestSetupInterceptors_AllCases(t *testing.T) {
	tests := []struct {
		name        string
		mtlsEnabled bool
	}{
		{"without mTLS", false},
		{"with mTLS", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{MTLSEnabled: tt.mtlsEnabled}
			interceptors := setupInterceptors(cfg, nil, nil)
			if interceptors == nil {
				t.Error("setupInterceptors should not return nil")
			}
		})
	}
}
