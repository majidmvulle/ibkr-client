package ibkr

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Ping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/api/tickle" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	err := client.Ping(context.Background())
	if err != nil {
		t.Errorf("Ping() error = %v", err)
	}
}

func TestClient_Ping_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	err := client.Ping(context.Background())
	if err == nil {
		t.Error("Expected error from Ping()")
	}
}

func TestClient_AuthStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"authenticated":true,"connected":true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	status, err := client.AuthStatus(context.Background())
	if err != nil {
		t.Fatalf("AuthStatus() error = %v", err)
	}
	if !status.Authenticated {
		t.Error("Expected authenticated = true")
	}
}

func TestClient_AuthStatus_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	_, err := client.AuthStatus(context.Background())
	if err == nil {
		t.Error("Expected error")
	}
}

func TestClient_AuthStatus_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	_, err := client.AuthStatus(context.Background())
	if err == nil {
		t.Error("Expected error")
	}
}

func TestClient_Reauthenticate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	err := client.Reauthenticate(context.Background())
	if err != nil {
		t.Errorf("Reauthenticate() error = %v", err)
	}
}

func TestClient_Reauthenticate_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	err := client.Reauthenticate(context.Background())
	if err == nil {
		t.Error("Expected error")
	}
}

func TestClient_GetAccounts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":"1","accountId":"U12345"}]`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	accounts, err := client.GetAccounts(context.Background())
	if err != nil {
		t.Fatalf("GetAccounts() error = %v", err)
	}
	if len(accounts) != 1 {
		t.Errorf("Expected 1 account, got %d", len(accounts))
	}
}

func TestClient_GetAccounts_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	_, err := client.GetAccounts(context.Background())
	if err == nil {
		t.Error("Expected error")
	}
}

func TestClient_GetAccounts_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	_, err := client.GetAccounts(context.Background())
	if err == nil {
		t.Error("Expected error")
	}
}
