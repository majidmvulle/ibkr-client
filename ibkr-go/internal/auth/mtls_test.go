package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTLSCredentials(t *testing.T) {
	// Create temporary test certificates
	tmpDir := t.TempDir()

	caCertPath := filepath.Join(tmpDir, "ca.pem")
	serverCertPath := filepath.Join(tmpDir, "server.pem")
	serverKeyPath := filepath.Join(tmpDir, "server-key.pem")

	// Write dummy certificate data
	caCert := []byte(`-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJAKHHCgVZU6krMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnRl
c3RjYTAeFw0yMTAxMDEwMDAwMDBaFw0zMTAxMDEwMDAwMDBaMBExDzANBgNVBAMM
BnRlc3RjYTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEA0Z91qmxEplnSMFeY
JYIWuJ8p7czIWW==
-----END CERTIFICATE-----`)

	serverCert := []byte(`-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJAKHHCgVZU6krMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnRl
c3RjYTAeFw0yMTAxMDEwMDAwMDBaFw0zMTAxMDEwMDAwMDBaMBExDzANBgNVBAMM
BnRlc3RjYTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEA0Z91qmxEplnSMFeY
JYIWuJ8p7czIWW==
-----END CERTIFICATE-----`)

	serverKey := []byte(`-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBANGfdapsRKZZ0jBX
mCWCFrifKe3MyFlv
-----END PRIVATE KEY-----`)

	if err := os.WriteFile(caCertPath, caCert, 0600); err != nil {
		t.Fatalf("Failed to write CA cert: %v", err)
	}
	if err := os.WriteFile(serverCertPath, serverCert, 0600); err != nil {
		t.Fatalf("Failed to write server cert: %v", err)
	}
	if err := os.WriteFile(serverKeyPath, serverKey, 0600); err != nil {
		t.Fatalf("Failed to write server key: %v", err)
	}

	// Test with valid paths (will fail due to invalid cert data, but tests the code path)
	_, err := LoadTLSCredentials(caCertPath, serverCertPath, serverKeyPath)
	// We expect an error because the certs are dummy data
	if err == nil {
		t.Error("Expected error with dummy certificates")
	}

	// Test with missing CA cert
	_, err = LoadTLSCredentials("/nonexistent/ca.pem", serverCertPath, serverKeyPath)
	if err == nil {
		t.Error("Expected error with missing CA cert")
	}

	// Test with missing server cert
	_, err = LoadTLSCredentials(caCertPath, "/nonexistent/server.pem", serverKeyPath)
	if err == nil {
		t.Error("Expected error with missing server cert")
	}
}
