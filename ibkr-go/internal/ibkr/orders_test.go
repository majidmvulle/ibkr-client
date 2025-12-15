package ibkr

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_PlaceOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"order_id":"12345","order_status":"Submitted"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	req := &PlaceOrderRequest{
		ConID:     12345,
		OrderType: "MKT",
		Side:      "BUY",
		Quantity:  100,
		Tif:       "DAY",
	}

	resp, err := client.PlaceOrder(context.Background(), req)
	if err != nil {
		t.Fatalf("PlaceOrder() error = %v", err)
	}
	if resp.OrderID != "12345" {
		t.Errorf("OrderID = %v, want 12345", resp.OrderID)
	}
}

func TestClient_ModifyOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"order_id":"12345","order_status":"Modified"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	req := &ModifyOrderRequest{
		Quantity: 200,
	}

	resp, err := client.ModifyOrder(context.Background(), "12345", req)
	if err != nil {
		t.Fatalf("ModifyOrder() error = %v", err)
	}
	if resp.OrderID != "12345" {
		t.Errorf("OrderID = %v, want 12345", resp.OrderID)
	}
}

func TestClient_CancelOrder(t *testing.T) {
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

func TestClient_GetLiveOrders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"orders":[{"orderId":"123","status":"Filled"}]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "U12345")
	orders, err := client.GetLiveOrders(context.Background())
	if err != nil {
		t.Fatalf("GetLiveOrders() error = %v", err)
	}
	if len(orders) != 1 {
		t.Errorf("Expected 1 order, got %d", len(orders))
	}
}
