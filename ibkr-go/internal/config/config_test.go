package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original env vars
	originalEnv := map[string]string{
		"APP_NAME":                os.Getenv("APP_NAME"),
		"APP_ENV":                 os.Getenv("APP_ENV"),
		"DB_WRITE_DSN":            os.Getenv("DB_WRITE_DSN"),
		"DB_READ_DSN":             os.Getenv("DB_READ_DSN"),
		"IBKR_GATEWAY_URL":        os.Getenv("IBKR_GATEWAY_URL"),
		"IBKR_ACCOUNT_ID":         os.Getenv("IBKR_ACCOUNT_ID"),
		"ENCRYPTION_KEY":          os.Getenv("ENCRYPTION_KEY"),
		"APP_HTTP_PORT":           os.Getenv("APP_HTTP_PORT"),
		"APP_GRPC_PORT":           os.Getenv("APP_GRPC_PORT"),
		"MTLS_ENABLED":            os.Getenv("MTLS_ENABLED"),
		"OTEL_COLLECTOR_ENDPOINT": os.Getenv("OTEL_COLLECTOR_ENDPOINT"),
	}

	// Restore env vars after test
	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"APP_NAME":         "test-app",
				"APP_ENV":          "test",
				"DB_WRITE_DSN":     "postgres://user:pass@localhost:5432/db",
				"DB_READ_DSN":      "postgres://user:pass@localhost:5432/db",
				"IBKR_GATEWAY_URL": "http://localhost:5000",
				"IBKR_ACCOUNT_ID":  "DU123456",
				"ENCRYPTION_KEY":   "12345678901234567890123456789012",
				"APP_HTTP_PORT":    "8080",
				"APP_GRPC_PORT":    "50051",
				"MTLS_ENABLED":     "false",
			},
			wantErr: false,
		},
		{
			name: "missing DB_WRITE_DSN",
			envVars: map[string]string{
				"APP_NAME":         "test-app",
				"APP_ENV":          "test",
				"DB_READ_DSN":      "postgres://user:pass@localhost:5432/db",
				"IBKR_GATEWAY_URL": "http://localhost:5000",
				"IBKR_ACCOUNT_ID":  "DU123456",
				"ENCRYPTION_KEY":   "12345678901234567890123456789012",
			},
			wantErr: true,
		},
		{
			name: "missing ENCRYPTION_KEY",
			envVars: map[string]string{
				"APP_NAME":         "test-app",
				"APP_ENV":          "test",
				"DB_WRITE_DSN":     "postgres://user:pass@localhost:5432/db",
				"DB_READ_DSN":      "postgres://user:pass@localhost:5432/db",
				"IBKR_GATEWAY_URL": "http://localhost:5000",
				"IBKR_ACCOUNT_ID":  "DU123456",
			},
			wantErr: true,
		},
		{
			name: "with mTLS enabled",
			envVars: map[string]string{
				"APP_NAME":              "test-app",
				"APP_ENV":               "test",
				"DB_WRITE_DSN":          "postgres://user:pass@localhost:5432/db",
				"DB_READ_DSN":           "postgres://user:pass@localhost:5432/db",
				"IBKR_GATEWAY_URL":      "http://localhost:5000",
				"IBKR_ACCOUNT_ID":       "DU123456",
				"ENCRYPTION_KEY":        "12345678901234567890123456789012",
				"MTLS_ENABLED":          "true",
				"MTLS_CA_CERT_PATH":     "/path/to/ca.crt",
				"MTLS_SERVER_CERT_PATH": "/path/to/server.crt",
				"MTLS_SERVER_KEY_PATH":  "/path/to/server.key",
			},
			wantErr: false,
		},
		{
			name: "with OpenTelemetry",
			envVars: map[string]string{
				"APP_NAME":                "test-app",
				"APP_ENV":                 "test",
				"DB_WRITE_DSN":            "postgres://user:pass@localhost:5432/db",
				"DB_READ_DSN":             "postgres://user:pass@localhost:5432/db",
				"IBKR_GATEWAY_URL":        "http://localhost:5000",
				"IBKR_ACCOUNT_ID":         "DU123456",
				"ENCRYPTION_KEY":          "12345678901234567890123456789012",
				"OTEL_COLLECTOR_ENDPOINT": "http://localhost:4317",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all env vars first
			for key := range originalEnv {
				os.Unsetenv(key)
			}

			// Set test env vars
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg, err := Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Verify config values
			if cfg == nil {
				t.Fatal("Load() returned nil config")
			}

			if cfg.AppName != tt.envVars["APP_NAME"] {
				t.Errorf("AppName = %v, want %v", cfg.AppName, tt.envVars["APP_NAME"])
			}

			if cfg.DBWriteDSN != tt.envVars["DB_WRITE_DSN"] {
				t.Errorf("DBWriteDSN = %v, want %v", cfg.DBWriteDSN, tt.envVars["DB_WRITE_DSN"])
			}

			if cfg.IBKRGatewayURL != tt.envVars["IBKR_GATEWAY_URL"] {
				t.Errorf("IBKRGatewayURL = %v, want %v", cfg.IBKRGatewayURL, tt.envVars["IBKR_GATEWAY_URL"])
			}

			if string(cfg.EncryptionKey) != tt.envVars["ENCRYPTION_KEY"] {
				t.Errorf("EncryptionKey = %v, want %v", string(cfg.EncryptionKey), tt.envVars["ENCRYPTION_KEY"])
			}
		})
	}
}

func TestDefaultValues(t *testing.T) {
	// Save and clear env vars
	originalEnv := map[string]string{
		"APP_NAME":         os.Getenv("APP_NAME"),
		"APP_ENV":          os.Getenv("APP_ENV"),
		"DB_WRITE_DSN":     os.Getenv("DB_WRITE_DSN"),
		"DB_READ_DSN":      os.Getenv("DB_READ_DSN"),
		"IBKR_GATEWAY_URL": os.Getenv("IBKR_GATEWAY_URL"),
		"IBKR_ACCOUNT_ID":  os.Getenv("IBKR_ACCOUNT_ID"),
		"ENCRYPTION_KEY":   os.Getenv("ENCRYPTION_KEY"),
		"APP_HTTP_PORT":    os.Getenv("APP_HTTP_PORT"),
		"APP_GRPC_PORT":    os.Getenv("APP_GRPC_PORT"),
		"MTLS_ENABLED":     os.Getenv("MTLS_ENABLED"),
	}

	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// Set only required vars
	os.Setenv("DB_WRITE_DSN", "postgres://user:pass@localhost:5432/db")
	os.Setenv("DB_READ_DSN", "postgres://user:pass@localhost:5432/db")
	os.Setenv("IBKR_GATEWAY_URL", "http://localhost:5000")
	os.Setenv("IBKR_ACCOUNT_ID", "DU123456")
	os.Setenv("ENCRYPTION_KEY", "12345678901234567890123456789012")

	// Unset optional vars to test defaults
	os.Unsetenv("APP_NAME")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("APP_HTTP_PORT")
	os.Unsetenv("APP_GRPC_PORT")
	os.Unsetenv("MTLS_ENABLED")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Check defaults are set
	if cfg.AppName == "" {
		t.Error("AppName should have a default value")
	}

	if cfg.AppEnv == "" {
		t.Error("AppEnv should have a default value")
	}

	if cfg.HTTPPort == 0 {
		t.Error("HTTPPort should have a default value")
	}

	if cfg.GRPCPort == 0 {
		t.Error("GRPCPort should have a default value")
	}
}

func TestEncryptionKeyLength(t *testing.T) {
	// Save env vars
	originalKey := os.Getenv("ENCRYPTION_KEY")
	defer func() {
		if originalKey == "" {
			os.Unsetenv("ENCRYPTION_KEY")
		} else {
			os.Setenv("ENCRYPTION_KEY", originalKey)
		}
	}()

	// Set required vars
	os.Setenv("DB_WRITE_DSN", "postgres://user:pass@localhost:5432/db")
	os.Setenv("DB_READ_DSN", "postgres://user:pass@localhost:5432/db")
	os.Setenv("IBKR_GATEWAY_URL", "http://localhost:5000")
	os.Setenv("IBKR_ACCOUNT_ID", "DU123456")

	tests := []struct {
		name          string
		encryptionKey string
		wantErr       bool
	}{
		{
			name:          "valid 32-byte key",
			encryptionKey: "12345678901234567890123456789012",
			wantErr:       false,
		},
		{
			name:          "too short key",
			encryptionKey: "short",
			wantErr:       true,
		},
		{
			name:          "too long key",
			encryptionKey: "123456789012345678901234567890123", // 33 bytes
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ENCRYPTION_KEY", tt.encryptionKey)

			_, err := Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
