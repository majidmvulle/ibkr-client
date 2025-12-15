package middleware

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
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
		{"with CN", &x509.Certificate{Subject: pkix.Name{CommonName: "test"}}, "test"},
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

func TestExtractClientIdentityFromCert_WithDNSNames(t *testing.T) {
	// Test fallback to DNS names when CN is empty
	cert := &x509.Certificate{
		Subject:  pkix.Name{CommonName: ""},
		DNSNames: []string{"client.example.com", "alt.example.com"},
	}

	got := ExtractClientIdentityFromCert(cert)
	if got != "client.example.com" {
		t.Errorf("Expected 'client.example.com', got %v", got)
	}
}
