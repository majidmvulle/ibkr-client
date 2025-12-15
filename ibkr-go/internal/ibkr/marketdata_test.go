package ibkr

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SearchContracts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"conid":265598,"symbol":"AAPL"}]`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	contracts, err := client.SearchContracts(context.Background(), "AAPL")
	if err != nil {
		t.Fatalf("SearchContracts() error = %v", err)
	}
	if len(contracts) != 1 {
		t.Errorf("Expected 1 contract, got %d", len(contracts))
	}
}

func TestClient_GetMarketData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"31":150.0,"84":149.5,"86":150.5}]`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	snapshots, err := client.GetMarketData(context.Background(), []int{265598}, nil)
	if err != nil {
		t.Fatalf("GetMarketData() error = %v", err)
	}
	if len(snapshots) != 1 {
		t.Errorf("Expected 1 snapshot, got %d", len(snapshots))
	}
}

func TestClient_GetHistoricalData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"t":1234567890,"o":100,"h":105,"l":99,"c":103,"v":50000}]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	data, err := client.GetHistoricalData(context.Background(), 265598, "1d", "1d")
	if err != nil {
		t.Fatalf("GetHistoricalData() error = %v", err)
	}
	if len(data.Data) != 1 {
		t.Errorf("Expected 1 bar, got %d", len(data.Data))
	}
}
