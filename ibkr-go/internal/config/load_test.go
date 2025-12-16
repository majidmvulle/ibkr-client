package config

import (
	"os"
	"testing"
)

func TestLoad_AllEnvVars(t *testing.T) {
	// Set all environment variables
	os.Setenv("APP_NAME", "test-app")
	os.Setenv("APP_ENV", "production")
	os.Setenv("APP_HTTP_PORT", "9090")
	os.Setenv("APP_GRPC_PORT", "50052")
	os.Setenv("APP_LOG_LEVEL", "1")
	os.Setenv("DB_WRITE_DSN", "postgres://write")
	os.Setenv("DB_READ_DSN", "postgres://read")
	os.Setenv("IBKR_GATEWAY_URL", "http://gateway:5000")
	os.Setenv("IBKR_ACCOUNT_ID", "U99999")
	os.Setenv("ENCRYPTION_KEY", "12345678901234567890123456789012")
	os.Setenv("OTEL_COLLECTOR_ENDPOINT", "localhost:4317")
	os.Setenv("MTLS_ENABLED", "true")
	os.Setenv("MTLS_SERVER_CERT_PATH", "/certs/server.pem")
	os.Setenv("MTLS_SERVER_KEY_PATH", "/certs/server-key.pem")
	os.Setenv("MTLS_CA_CERT_PATH", "/certs/ca.pem")

	defer func() {
		os.Clearenv()
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AppName != "test-app" {
		t.Errorf("AppName = %v, want test-app", cfg.AppName)
	}
	if cfg.AppEnv != "production" {
		t.Errorf("AppEnv = %v, want production", cfg.AppEnv)
	}
	if cfg.HTTPPort != 9090 {
		t.Errorf("HTTPPort = %v, want 9090", cfg.HTTPPort)
	}
	if cfg.GRPCPort != 50052 {
		t.Errorf("GRPCPort = %v, want 50052", cfg.GRPCPort)
	}
	if !cfg.MTLSEnabled {
		t.Error("MTLSEnabled should be true")
	}
}
