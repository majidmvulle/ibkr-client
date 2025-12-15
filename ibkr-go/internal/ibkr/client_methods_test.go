package ibkr

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SearchContracts_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Contract{
			{ConID: 265598, Symbol: "AAPL"},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	contracts, err := client.SearchContracts(context.Background(), "AAPL")
	if err != nil {
		t.Errorf("SearchContracts() error = %v", err)
	}
	if len(contracts) != 1 {
		t.Errorf("Expected 1 contract, got %d", len(contracts))
	}
	if contracts[0].Symbol != "AAPL" {
		t.Errorf("Symbol = %v, want AAPL", contracts[0].Symbol)
	}
}

func TestClient_SearchContracts_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	_, err := client.SearchContracts(context.Background(), "INVALID")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestClient_GetMarketData_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]MarketDataSnapshot{
			{LastPrice: 150.50, Bid: 150.25, Ask: 150.75},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	snapshots, err := client.GetMarketData(context.Background(), []int{265598}, nil)
	if err != nil {
		t.Errorf("GetMarketData() error = %v", err)
	}
	if len(snapshots) != 1 {
		t.Errorf("Expected 1 snapshot, got %d", len(snapshots))
	}
}

func TestClient_GetMarketData_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	_, err := client.GetMarketData(context.Background(), []int{999999}, nil)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestClient_GetHistoricalData_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&HistoricalDataResponse{
			Data: []HistoricalBar{
				{Time: 1234567890, Open: 150.0, High: 151.0, Low: 149.0, Close: 150.5, Volume: 1000000},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	resp, err := client.GetHistoricalData(context.Background(), 265598, "1d", "1h")
	if err != nil {
		t.Errorf("GetHistoricalData() error = %v", err)
	}
	if len(resp.Data) != 1 {
		t.Errorf("Expected 1 bar, got %d", len(resp.Data))
	}
}

func TestClient_PlaceOrder_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&OrderResponse{
			OrderID: "12345",
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	req := &PlaceOrderRequest{
		ConID:     265598,
		OrderType: "MKT",
		Side:      "BUY",
		Quantity:  100,
	}
	resp, err := client.PlaceOrder(context.Background(), req)
	if err != nil {
		t.Errorf("PlaceOrder() error = %v", err)
	}
	if resp.OrderID != "12345" {
		t.Errorf("OrderID = %v, want 12345", resp.OrderID)
	}
}

func TestClient_PlaceOrder_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	_, err := client.PlaceOrder(context.Background(), &PlaceOrderRequest{})
	if err == nil {
		t.Error("Expected error")
	}
}

func TestClient_PlaceOrder_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	_, err := client.PlaceOrder(context.Background(), &PlaceOrderRequest{})
	if err == nil {
		t.Error("Expected error")
	}
}

func TestClient_ModifyOrder_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&OrderResponse{
			OrderID: "12345",
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	req := &ModifyOrderRequest{
		Quantity: 200,
	}
	resp, err := client.ModifyOrder(context.Background(), "12345", req)
	if err != nil {
		t.Errorf("ModifyOrder() error = %v", err)
	}
	if resp.OrderID != "12345" {
		t.Errorf("OrderID = %v, want 12345", resp.OrderID)
	}
}

func TestClient_ModifyOrder_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	_, err := client.ModifyOrder(context.Background(), "123", &ModifyOrderRequest{})
	if err == nil {
		t.Error("Expected error")
	}
}

func TestClient_ModifyOrder_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	_, err := client.ModifyOrder(context.Background(), "123", &ModifyOrderRequest{})
	if err == nil {
		t.Error("Expected error")
	}
}

func TestClient_CancelOrder_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	err := client.CancelOrder(context.Background(), "12345")
	if err != nil {
		t.Errorf("CancelOrder() error = %v", err)
	}
}

func TestClient_CancelOrder_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	err := client.CancelOrder(context.Background(), "12345")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestClient_GetLiveOrders_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"orders": []Order{
				{OrderID: "12345", Status: "Filled"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	orders, err := client.GetLiveOrders(context.Background())
	if err != nil {
		t.Errorf("GetLiveOrders() error = %v", err)
	}
	if len(orders) != 1 {
		t.Errorf("Expected 1 order, got %d", len(orders))
	}
}

func TestClient_GetLiveOrders_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	_, err := client.GetLiveOrders(context.Background())
	if err == nil {
		t.Error("Expected error")
	}
}

func TestClient_GetLiveOrders_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	_, err := client.GetLiveOrders(context.Background())
	if err == nil {
		t.Error("Expected error")
	}
}
