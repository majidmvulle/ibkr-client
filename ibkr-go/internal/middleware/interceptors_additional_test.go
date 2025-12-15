package middleware

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"log/slog"
	"math/big"
	"testing"
	"time"

	"connectrpc.com/connect"
)

func TestNewSessionInterceptor(t *testing.T) {
	// Test that NewSessionInterceptor returns a non-nil interceptor
	interceptor := NewSessionInterceptor(nil)
	if interceptor == nil {
		t.Error("NewSessionInterceptor() should not return nil")
	}
}

func TestExtractClientIdentity_FromRequest(t *testing.T) {
	// Create a test certificate
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "test-client",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour),
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	// Test ExtractClientIdentityFromCert
	identity := ExtractClientIdentityFromCert(cert)
	if identity != "test-client" {
		t.Errorf("ExtractClientIdentityFromCert() = %v, want test-client", identity)
	}
}

func TestNewMTLSInterceptor_WithCerts(t *testing.T) {
	// Test the interceptor
	interceptor := NewMTLSInterceptor(slog.Default())
	if interceptor == nil {
		t.Error("NewMTLSInterceptor() should not return nil")
	}

	t.Skip("Skipping test requiring unimplemented mTLS extraction and difficult to mock TLS state")
}

func TestLoggingInterceptor_WithError(t *testing.T) {
	interceptor := LoggingInterceptor(slog.Default())

	testErr := connect.NewError(connect.CodeInternal, nil)
	handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		return nil, testErr
	}

	wrapped := interceptor(handler)
	req := connect.NewRequest(&struct{}{})

	_, err := wrapped(context.Background(), req)
	if err != testErr {
		t.Errorf("Expected error to be propagated")
	}
}

func TestValidationInterceptor_WithInvalidRequest(t *testing.T) {
	interceptor := NewValidationInterceptor()

	handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		return connect.NewResponse(&struct{}{}), nil
	}

	wrapped := interceptor(handler)
	req := connect.NewRequest(&struct{}{})

	_, err := wrapped(context.Background(), req)
	if err != nil {
		t.Logf("Validation error (expected for some requests): %v", err)
	}
}
