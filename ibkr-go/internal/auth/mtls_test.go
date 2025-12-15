package auth

import (
	"testing"
)

func TestLoadTLSCredentials_AllPaths(t *testing.T) {
	// Test with non-existent files to cover error paths
	_, err := LoadTLSCredentials("/nonexistent/ca.pem", "/nonexistent/cert.pem", "/nonexistent/key.pem")
	if err == nil {
		t.Error("Expected error with non-existent files")
	}

	// Test with empty paths
	_, err = LoadTLSCredentials("", "", "")
	if err == nil {
		t.Error("Expected error with empty paths")
	}
}

func TestMTLSInterceptor_Coverage(t *testing.T) {
	interceptor := MTLSInterceptor()
	if interceptor == nil {
		t.Error("MTLSInterceptor should not return nil")
	}
}
