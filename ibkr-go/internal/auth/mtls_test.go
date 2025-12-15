package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTLSCredentials_InvalidPaths(t *testing.T) {
	tests := []struct {
		name           string
		caCertPath     string
		serverCertPath string
		serverKeyPath  string
		wantErr        bool
	}{
		{
			name:           "invalid CA cert path",
			caCertPath:     "/nonexistent/ca.pem",
			serverCertPath: "/tmp/server.pem",
			serverKeyPath:  "/tmp/server-key.pem",
			wantErr:        true,
		},
		{
			name:           "invalid server cert path",
			caCertPath:     "/tmp/ca.pem",
			serverCertPath: "/nonexistent/server.pem",
			serverKeyPath:  "/tmp/server-key.pem",
			wantErr:        true,
		},
		{
			name:           "all invalid paths",
			caCertPath:     "/nonexistent/ca.pem",
			serverCertPath: "/nonexistent/server.pem",
			serverKeyPath:  "/nonexistent/server-key.pem",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadTLSCredentials(tt.caCertPath, tt.serverCertPath, tt.serverKeyPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTLSCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadTLSCredentials_InvalidPEM(t *testing.T) {
	// Create temp directory for test files
	tmpDir := t.TempDir()

	// Create invalid PEM file (not a valid certificate)
	invalidCA := filepath.Join(tmpDir, "invalid-ca.pem")
	if err := os.WriteFile(invalidCA, []byte("not a valid PEM"), 0600); err != nil {
		t.Fatal(err)
	}

	// Create dummy cert and key files (will fail at CA parsing)
	serverCert := filepath.Join(tmpDir, "server.pem")
	serverKey := filepath.Join(tmpDir, "server-key.pem")
	if err := os.WriteFile(serverCert, []byte("dummy"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(serverKey, []byte("dummy"), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := LoadTLSCredentials(invalidCA, serverCert, serverKey)
	if err == nil {
		t.Error("Expected error with invalid CA PEM, got nil")
	}
}

func TestMTLSInterceptor_ReturnsFunction(t *testing.T) {
	interceptor := MTLSInterceptor()
	if interceptor == nil {
		t.Error("MTLSInterceptor() returned nil")
	}
}
