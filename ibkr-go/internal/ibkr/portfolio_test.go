package ibkr

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetPortfolio(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"contractDesc":"AAPL","position":100}]`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	positions, err := client.GetPortfolio(context.Background())
	if err != nil {
		t.Fatalf("GetPortfolio() error = %v", err)
	}
	if len(positions) != 1 {
		t.Errorf("Expected 1 position, got %d", len(positions))
	}
}

func TestClient_GetAccountSummary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"netliquidation":{"amount":"10000"},"currency":"USD"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	summary, err := client.GetAccountSummary(context.Background())
	if err != nil {
		t.Fatalf("GetAccountSummary() error = %v", err)
	}
	if summary == nil {
		t.Error("Expected non-nil summary")
	}
}
