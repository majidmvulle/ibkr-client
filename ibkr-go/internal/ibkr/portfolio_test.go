package ibkr

import (
	"context"
	"encoding/json"
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
		// Return valid JSON matching AccountSummary struct
		json.NewEncoder(w).Encode(map[string]interface{}{
			"accountcode":                    "U12345",
			"accounttype":                    "INDIVIDUAL",
			"netliquidation":                 100000.50,
			"totalcashvalue":                 50000.25,
			"settledcash":                    45000.00,
			"accruedcash":                    100.50,
			"buyingpower":                    200000.00,
			"equitywithloanvalue":            95000.00,
			"previousdayequitywithloanvalue": 94000.00,
			"grosspositionvalue":             50000.00,
			"regtequity":                     100000.00,
			"regtmargin":                     25000.00,
			"sma":                            10000.00,
			"currency":                       "USD",
		})
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
