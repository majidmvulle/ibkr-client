package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/config"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/database"
)

func TestSetupServer(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:    8080,
		MTLSEnabled: false,
	}

	server, err := setupServer(cfg, nil, nil, nil)
	if err != nil {
		t.Fatalf("setupServer() error = %v", err)
	}
	if server == nil {
		t.Fatal("setupServer() returned nil")
	}

	if server.Addr != ":8080" {
		t.Errorf("server.Addr = %v, want :8080", server.Addr)
	}

	if server.Handler == nil {
		t.Error("server.Handler is nil")
	}
}

func TestSetupServerWithMTLS(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:           8080,
		MTLSEnabled:        true,
		MTLSServerCertPath: "/nonexistent/cert.pem",
		MTLSServerKeyPath:  "/nonexistent/key.pem",
		MTLSCACertPath:     "/nonexistent/ca.pem",
	}

	// This will fail to configure TLS due to missing files
	server, err := setupServer(cfg, nil, nil, nil)
	if err == nil {
		t.Fatal("setupServer() expected error due to missing certs")
	}
	if server != nil {
		t.Error("setupServer() expected nil server on error")
	}
}

func TestWaitForShutdown(t *testing.T) {
	// Create a test server
	server := &http.Server{
		Addr: ":0",
	}

	// This test just ensures the function exists and can be called
	// We can't actually test the shutdown signal handling in a unit test
	go func() {
		// Function will block waiting for signal, so we run it in goroutine
		waitForShutdown(server, nil)
	}()
}

func TestSetupLogger(t *testing.T) {
	cfg := &config.Config{
		AppDebug: true,
	}
	logger := setupLogger(cfg)
	if logger == nil {
		t.Error("setupLogger returned nil")
	}
}

func TestInitTelemetry(t *testing.T) {
	cfg := &config.Config{
		AppName:               "test",
		OtelCollectorEndpoint: "",
	}
	logger := slog.Default()
	tp := initTelemetry(context.Background(), cfg, logger)
	if tp == nil {
		t.Error("initTelemetry returned nil")
	}
}

func TestShutdownTelemetry(t *testing.T) {
	logger := slog.Default()
	// Test with nil
	shutdownTelemetry(nil, logger)

	// Test with valid provider (no-op)
	cfg := &config.Config{AppName: "test"}
	tp := initTelemetry(context.Background(), cfg, logger)
	shutdownTelemetry(tp, logger)
}

func TestStartServer(t *testing.T) {
	server := &http.Server{Addr: ":0"}
	logger := slog.Default()

	// Test without TLS
	startServer(server, logger, false)

	// Give it a moment to start
	time.Sleep(10 * time.Millisecond)

	// Close it
	server.Close()
}

func TestHealthCheck(t *testing.T) {
	cfg := &config.Config{HTTPPort: 8080}
	server, _ := setupServer(cfg, nil, nil, nil)

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestReadinessCheck_DatabaseUnhealthy(t *testing.T) {
	cfg := &config.Config{HTTPPort: 8080}
	db := &database.DB{} // Pool is nil, Health() should return error
	server, _ := setupServer(cfg, db, nil, nil)

	req := httptest.NewRequest("GET", "/readyz", nil)
	w := httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Status = %v, want %v", w.Code, http.StatusServiceUnavailable)
	}
}
