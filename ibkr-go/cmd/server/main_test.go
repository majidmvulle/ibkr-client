package main

import (
	"net/http"
	"testing"

	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/config"
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
