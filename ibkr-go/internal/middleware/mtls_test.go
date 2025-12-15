package middleware

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"log/slog"
	"testing"

	"connectrpc.com/connect"
)

func TestExtractClientIdentityFromCert(t *testing.T) {
	tests := []struct {
		name     string
		cert     *x509.Certificate
		expected string
	}{
		{
			name:     "nil certificate",
			cert:     nil,
			expected: "",
		},
		{
			name: "certificate with common name",
			cert: &x509.Certificate{
				Subject: pkix.Name{
					CommonName: "test-service",
				},
			},
			expected: "test-service",
		},
		{
			name: "certificate with DNS names",
			cert: &x509.Certificate{
				Subject: pkix.Name{
					CommonName: "",
				},
				DNSNames: []string{"service1.example.com", "service2.example.com"},
			},
			expected: "service1.example.com",
		},
		{
			name: "certificate with both CN and DNS names prefers CN",
			cert: &x509.Certificate{
				Subject: pkix.Name{
					CommonName: "primary-service",
				},
				DNSNames: []string{"service1.example.com"},
			},
			expected: "primary-service",
		},
		{
			name: "certificate with no identity",
			cert: &x509.Certificate{
				Subject: pkix.Name{
					CommonName: "",
				},
				DNSNames: []string{},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractClientIdentityFromCert(tt.cert)
			if result != tt.expected {
				t.Errorf("ExtractClientIdentityFromCert() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewMTLSInterceptor(t *testing.T) {
	// Create a test logger
	logger := slog.Default()

	// Create the interceptor
	interceptor := NewMTLSInterceptor(logger)

	// Verify it returns a function
	if interceptor == nil {
		t.Fatal("NewMTLSInterceptor() returned nil")
	}

	// Create a mock next handler
	nextCalled := false
	next := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		nextCalled = true
		return connect.NewResponse(&struct{}{}), nil
	}

	// Wrap with interceptor
	handler := interceptor(next)

	// Create a test request with no peer info
	req := connect.NewRequest(&struct{}{})

	// Call should fail due to no peer information
	_, err := handler(context.Background(), req)
	if err == nil {
		t.Error("Expected error for request with no peer information, got nil")
	}

	// Verify next was not called
	if nextCalled {
		t.Error("Next handler should not be called when peer info is missing")
	}
}
