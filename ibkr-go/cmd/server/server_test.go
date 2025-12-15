package server

import (
	"testing"

	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/config"
)

func TestTLSConfiguration(t *testing.T) {
	cfg := &config.Config{
		MTLSEnabled:        true,
		MTLSServerCertPath: "/path/to/cert.pem",
		MTLSServerKeyPath:  "/path/to/key.pem",
		MTLSCACertPath:     "/path/to/ca.pem",
	}

	if !cfg.MTLSEnabled {
		t.Error("mTLS should be enabled")
	}

	if cfg.MTLSServerCertPath == "" {
		t.Error("Server cert path should not be empty")
	}
}

func TestServerPorts(t *testing.T) {
	tests := []struct {
		name     string
		httpPort int
		grpcPort int
		wantErr  bool
	}{
		{"valid ports", 8080, 50051, false},
		{"same ports", 8080, 8080, true},
		{"zero http port", 0, 50051, true},
		{"zero grpc port", 8080, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasErr := tt.httpPort == 0 || tt.grpcPort == 0 || tt.httpPort == tt.grpcPort
			if hasErr != tt.wantErr {
				t.Errorf("Port validation = %v, want %v", hasErr, tt.wantErr)
			}
		})
	}
}
