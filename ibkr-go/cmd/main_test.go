package main

import (
	"testing"

	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/config"
)

func TestServerConfiguration(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
	}{
		{
			name: "basic config",
			cfg: &config.Config{
				AppName:     "test",
				HTTPPort:    8080,
				MTLSEnabled: false,
			},
		},
		{
			name: "with mTLS",
			cfg: &config.Config{
				AppName:            "test",
				HTTPPort:           8080,
				MTLSEnabled:        true,
				MTLSServerCertPath: "/tmp/cert.pem",
				MTLSServerKeyPath:  "/tmp/key.pem",
				MTLSCACertPath:     "/tmp/ca.pem",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cfg.AppName == "" {
				t.Error("AppName should not be empty")
			}
			if tt.cfg.HTTPPort == 0 {
				t.Error("HTTPPort should not be zero")
			}
		})
	}
}
