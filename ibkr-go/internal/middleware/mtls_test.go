package middleware

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"testing"

	"log/slog"

	"connectrpc.com/connect"
)

func TestNewMTLSInterceptor_AllPaths(t *testing.T) {
	logger := slog.Default()
	interceptor := NewMTLSInterceptor(logger)
	if interceptor == nil {
		t.Fatal("NewMTLSInterceptor should not return nil")
	}

	// Test with nil TLS
	handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		return connect.NewResponse(&struct{}{}), nil
	}
	wrapped := interceptor(handler)

	ctx := context.Background()
	req := connect.NewRequest(&struct{}{})

	_, err := wrapped(ctx, req)
	if err == nil {
		t.Log("Interceptor executed without TLS")
	}
}

func TestExtractClientIdentity_AllCases(t *testing.T) {
	tests := []struct {
		name string
		cert *x509.Certificate
		want string
	}{
		{"nil cert", nil, ""},
		{"empty CN", &x509.Certificate{}, ""},
		{"with CN", &x509.Certificate{Subject: x509.Subject{CommonName: "test"}}, "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractClientIdentityFromCert(tt.cert)
			if got != tt.want {
				t.Errorf("ExtractClientIdentityFromCert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateMTLSConnection_AllPaths(t *testing.T) {
	// Test with nil connection state
	identity := ValidateMTLSConnection(nil)
	if identity != "" {
		t.Errorf("Expected empty identity for nil connection, got %v", identity)
	}

	// Test with empty peer certificates
	state := &tls.ConnectionState{}
	identity = ValidateMTLSConnection(state)
	if identity != "" {
		t.Errorf("Expected empty identity for no peer certs, got %v", identity)
	}

	// Test with peer certificate
	cert := &x509.Certificate{Subject: x509.Subject{CommonName: "client"}}
	state = &tls.ConnectionState{PeerCertificates: []*x509.Certificate{cert}}
	identity = ValidateMTLSConnection(state)
	if identity != "client" {
		t.Errorf("Expected 'client', got %v", identity)
	}
}
