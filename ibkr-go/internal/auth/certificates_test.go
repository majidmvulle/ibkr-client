package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadTLSCredentials_ValidCertificates(t *testing.T) {
	// Create temporary directory for test certificates
	tmpDir := t.TempDir()

	// Generate test CA certificate
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate CA key: %v", err)
	}

	caTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test CA"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caCertBytes, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("Failed to create CA certificate: %v", err)
	}

	// Write CA certificate
	caCertPath := filepath.Join(tmpDir, "ca.pem")
	caCertFile, err := os.Create(caCertPath)
	if err != nil {
		t.Fatalf("Failed to create CA cert file: %v", err)
	}
	pem.Encode(caCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: caCertBytes})
	caCertFile.Close()

	// Generate server certificate
	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate server key: %v", err)
	}

	serverTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{"Test Server"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(time.Hour),
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	serverCertBytes, err := x509.CreateCertificate(rand.Reader, serverTemplate, caTemplate, &serverKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("Failed to create server certificate: %v", err)
	}

	// Write server certificate
	serverCertPath := filepath.Join(tmpDir, "server.pem")
	serverCertFile, err := os.Create(serverCertPath)
	if err != nil {
		t.Fatalf("Failed to create server cert file: %v", err)
	}
	pem.Encode(serverCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: serverCertBytes})
	serverCertFile.Close()

	// Write server key
	serverKeyPath := filepath.Join(tmpDir, "server-key.pem")
	serverKeyFile, err := os.Create(serverKeyPath)
	if err != nil {
		t.Fatalf("Failed to create server key file: %v", err)
	}
	pem.Encode(serverKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(serverKey)})
	serverKeyFile.Close()

	// Test LoadTLSCredentials with valid certificates
	creds, err := LoadTLSCredentials(caCertPath, serverCertPath, serverKeyPath)
	if err != nil {
		t.Errorf("LoadTLSCredentials() with valid certs error = %v", err)
	}
	if creds == nil {
		t.Error("Expected non-nil credentials")
	}
}

func TestMTLSInterceptor_NotNil(t *testing.T) {
	interceptor := MTLSInterceptor()
	if interceptor == nil {
		t.Error("MTLSInterceptor() should not return nil")
	}
}
