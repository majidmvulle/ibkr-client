package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/config"
)

// configureTLS sets up TLS configuration for the HTTP server with client certificate validation.
func configureTLS(server *http.Server, cfg *config.Config) error {
	// Load server certificate and key.
	cert, err := tls.LoadX509KeyPair(cfg.MTLSServerCertPath, cfg.MTLSServerKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load server certificate: %w", err)
	}

	// Load CA certificate for validating client certificates.
	caCert, err := os.ReadFile(cfg.MTLSCACertPath)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return fmt.Errorf("failed to parse CA certificate")
	}

	// Configure TLS with client certificate verification.
	server.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}

	slog.Info("mTLS configured",
		slog.String("server_cert", cfg.MTLSServerCertPath),
		slog.String("ca_cert", cfg.MTLSCACertPath),
	)

	return nil
}
