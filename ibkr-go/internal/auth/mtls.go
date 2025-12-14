package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"connectrpc.com/connect"
	"google.golang.org/grpc/credentials"
)

// LoadTLSCredentials loads mTLS credentials from certificate files
func LoadTLSCredentials(caCertPath, serverCertPath, serverKeyPath string) (credentials.TransportCredentials, error) {
	// Load CA certificate
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add CA certificate to pool")
	}

	// Load server certificate and key
	serverCert, err := tls.LoadX509KeyPair(serverCertPath, serverKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %w", err)
	}

	// Create TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}

	return credentials.NewTLS(tlsConfig), nil
}

// MTLSInterceptor validates client certificates
func MTLSInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// TODO: Extract and validate client certificate from context
			// For now, this is a placeholder
			return next(ctx, req)
		}
	}
}
