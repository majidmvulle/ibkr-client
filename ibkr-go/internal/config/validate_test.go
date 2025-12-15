package config

import (
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				AppName:        "test",
				DBWriteDSN:     "postgres://test",
				IBKRGatewayURL: "http://localhost:5000",
				IBKRAccountID:  "U12345",
				EncryptionKey:  []byte("0123456789abcdef0123456789abcdef"),
			},
			wantErr: false,
		},
		{
			name: "missing app name",
			cfg: &Config{
				DBWriteDSN:     "postgres://test",
				IBKRGatewayURL: "http://localhost:5000",
				IBKRAccountID:  "U12345",
				EncryptionKey:  []byte("0123456789abcdef0123456789abcdef"),
			},
			wantErr: true,
		},
		{
			name: "missing DB DSN",
			cfg: &Config{
				AppName:        "test",
				IBKRGatewayURL: "http://localhost:5000",
				IBKRAccountID:  "U12345",
				EncryptionKey:  []byte("0123456789abcdef0123456789abcdef"),
			},
			wantErr: true,
		},
		{
			name: "missing IBKR gateway URL",
			cfg: &Config{
				AppName:       "test",
				DBWriteDSN:    "postgres://test",
				IBKRAccountID: "U12345",
				EncryptionKey: []byte("0123456789abcdef0123456789abcdef"),
			},
			wantErr: true,
		},
		{
			name: "missing IBKR account ID",
			cfg: &Config{
				AppName:        "test",
				DBWriteDSN:     "postgres://test",
				IBKRGatewayURL: "http://localhost:5000",
				EncryptionKey:  []byte("0123456789abcdef0123456789abcdef"),
			},
			wantErr: true,
		},
		{
			name: "missing encryption key",
			cfg: &Config{
				AppName:        "test",
				DBWriteDSN:     "postgres://test",
				IBKRGatewayURL: "http://localhost:5000",
				IBKRAccountID:  "U12345",
			},
			wantErr: true,
		},
		{
			name: "short encryption key",
			cfg: &Config{
				AppName:        "test",
				DBWriteDSN:     "postgres://test",
				IBKRGatewayURL: "http://localhost:5000",
				IBKRAccountID:  "U12345",
				EncryptionKey:  []byte("short"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadWithDefaults(t *testing.T) {
	// Clear all env vars
	os.Clearenv()

	// Set only required vars
	os.Setenv("APP_NAME", "test")
	os.Setenv("DB_WRITE_DSN", "postgres://test")
	os.Setenv("IBKR_GATEWAY_URL", "http://localhost:5000")
	os.Setenv("IBKR_ACCOUNT_ID", "U12345")
	os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Check defaults
	if cfg.AppEnv != "development" {
		t.Errorf("Default AppEnv = %v, want development", cfg.AppEnv)
	}
	if cfg.HTTPPort != 8080 {
		t.Errorf("Default HTTPPort = %v, want 8080", cfg.HTTPPort)
	}
	if cfg.GRPCPort != 50051 {
		t.Errorf("Default GRPCPort = %v, want 50051", cfg.GRPCPort)
	}
	if cfg.LogLevel != 0 {
		t.Errorf("Default LogLevel = %v, want 0", cfg.LogLevel)
	}
	if cfg.MTLSEnabled != false {
		t.Errorf("Default MTLSEnabled = %v, want false", cfg.MTLSEnabled)
	}
}
