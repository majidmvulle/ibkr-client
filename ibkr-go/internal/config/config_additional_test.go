package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("APP_NAME", "test-app")
	os.Setenv("APP_ENV", "test")
	os.Setenv("HTTP_PORT", "8080")
	os.Setenv("GRPC_PORT", "50051")
	os.Setenv("DB_WRITE_DSN", "postgres://test")
	os.Setenv("IBKR_GATEWAY_URL", "http://localhost:5000")
	os.Setenv("IBKR_ACCOUNT_ID", "U12345")
	os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")

	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("APP_ENV")
		os.Unsetenv("HTTP_PORT")
		os.Unsetenv("GRPC_PORT")
		os.Unsetenv("DB_WRITE_DSN")
		os.Unsetenv("IBKR_GATEWAY_URL")
		os.Unsetenv("IBKR_ACCOUNT_ID")
		os.Unsetenv("ENCRYPTION_KEY")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AppName != "test-app" {
		t.Errorf("AppName = %v, want test-app", cfg.AppName)
	}
	if cfg.HTTPPort != 8080 {
		t.Errorf("HTTPPort = %v, want 8080", cfg.HTTPPort)
	}
	if cfg.GRPCPort != 50051 {
		t.Errorf("GRPCPort = %v, want 50051", cfg.GRPCPort)
	}
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "test-value")
	defer os.Unsetenv("TEST_VAR")

	val := getEnv("TEST_VAR", "default")
	if val != "test-value" {
		t.Errorf("getEnv() = %v, want test-value", val)
	}

	val = getEnv("NONEXISTENT", "default")
	if val != "default" {
		t.Errorf("getEnv() = %v, want default", val)
	}
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	val := getEnvInt("TEST_INT", 10)
	if val != 42 {
		t.Errorf("getEnvInt() = %v, want 42", val)
	}

	val = getEnvInt("NONEXISTENT", 10)
	if val != 10 {
		t.Errorf("getEnvInt() = %v, want 10", val)
	}

	os.Setenv("TEST_INVALID", "not-a-number")
	defer os.Unsetenv("TEST_INVALID")
	val = getEnvInt("TEST_INVALID", 10)
	if val != 10 {
		t.Errorf("getEnvInt() with invalid value = %v, want 10", val)
	}
}

func TestGetEnvBool(t *testing.T) {
	os.Setenv("TEST_BOOL", "true")
	defer os.Unsetenv("TEST_BOOL")

	val := getEnvBool("TEST_BOOL", false)
	if val != true {
		t.Errorf("getEnvBool() = %v, want true", val)
	}

	val = getEnvBool("NONEXISTENT", false)
	if val != false {
		t.Errorf("getEnvBool() = %v, want false", val)
	}
}
