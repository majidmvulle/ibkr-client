package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	// DefaultGRPCPort is the default port for gRPC server.
	DefaultGRPCPort = 50051
	// DefaultHTTPPort is the default port for HTTP server.
	DefaultHTTPPort = 8080
	// DefaultDBLogLevel is the default database log level.
	DefaultDBLogLevel = 2
	// RequiredEncryptionKeyLength is the required length for AES-256 encryption key.
	RequiredEncryptionKeyLength = 32
)

// Config holds all application configuration.
type Config struct {
	// App.
	AppName    string
	AppEnv     string
	AppDebug   bool
	GRPCPort   int
	HTTPPort   int
	AppVersion string
	LogLevel   int

	// Database.
	DBWriteDSN string
	DBReadDSN  string
	DBLogLevel int

	// IBKR Gateway.
	IBKRGatewayURL     string
	IBKRGatewayKeyPath string
	IBKRAccountID      string

	// mTLS.
	MTLSEnabled        bool
	MTLSCACertPath     string
	MTLSServerCertPath string
	MTLSServerKeyPath  string

	// Encryption.
	EncryptionKey []byte

	// OpenTelemetry.
	OtelCollectorEndpoint string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		AppName:    getEnv("APP_NAME", "ibkr-client"),
		AppEnv:     getEnv("APP_ENV", "local"),
		AppDebug:   getEnvBool("APP_DEBUG", true),
		GRPCPort:   getEnvInt("APP_GRPC_PORT", DefaultGRPCPort),
		HTTPPort:   getEnvInt("APP_HTTP_PORT", DefaultHTTPPort),
		AppVersion: getEnv("APP_VERSION", "v0.0.1"),
		LogLevel:   getEnvInt("APP_LOG_LEVEL", 0),

		DBWriteDSN: getEnv("DB_WRITE_DSN", ""),
		DBReadDSN:  getEnv("DB_READ_DSN", ""),
		DBLogLevel: getEnvInt("DB_LOG_LEVEL", DefaultDBLogLevel),

		IBKRGatewayURL:     getEnv("IBKR_GATEWAY_URL", "http://localhost:5000"),
		IBKRGatewayKeyPath: getEnv("IBKR_GATEWAY_KEY_PATH", ""),
		IBKRAccountID:      getEnv("IBKR_ACCOUNT_ID", ""),

		MTLSEnabled:        getEnvBool("MTLS_ENABLED", false),
		MTLSCACertPath:     getEnv("MTLS_CA_CERT_PATH", ""),
		MTLSServerCertPath: getEnv("MTLS_SERVER_CERT_PATH", ""),
		MTLSServerKeyPath:  getEnv("MTLS_SERVER_KEY_PATH", ""),

		OtelCollectorEndpoint: getEnv("OTEL_COLLECTOR_ENDPOINT", ""),
	}

	// Load encryption key.
	encKeyStr := getEnv("ENCRYPTION_KEY", "")
	if encKeyStr != "" {
		cfg.EncryptionKey = []byte(encKeyStr)
		if len(cfg.EncryptionKey) != RequiredEncryptionKeyLength {
			return nil, fmt.Errorf("ENCRYPTION_KEY must be exactly 32 bytes, got %d", len(cfg.EncryptionKey))
		}
	}

	// Validate required fields.
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if required configuration is present.
func (c *Config) Validate() error {
	if c.DBWriteDSN == "" {
		return fmt.Errorf("DB_WRITE_DSN is required")
	}

	if c.DBReadDSN == "" {
		return fmt.Errorf("DB_READ_DSN is required")
	}

	if len(c.EncryptionKey) != RequiredEncryptionKeyLength {
		return fmt.Errorf("ENCRYPTION_KEY must be set and be 32 bytes")
	}

	return nil
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

// getEnvInt retrieves an environment variable as int or returns a default value.
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	return defaultValue
}

// getEnvBool retrieves an environment variable as bool or returns a default value.
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}

	return defaultValue
}
