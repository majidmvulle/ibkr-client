package api

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	orderv1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/order/v1"
)

type mockIBKRClient struct {
	placeOrderFunc    func(ctx context.Context, req *ibkr.PlaceOrderRequest) (*ibkr.OrderResponse, error)
	modifyOrderFunc   func(ctx context.Context, orderID string, req *ibkr.ModifyOrderRequest) (*ibkr.OrderResponse, error)
	cancelOrderFunc   func(ctx context.Context, orderID string) error
	getLiveOrdersFunc func(ctx context.Context) ([]ibkr.Order, error)
}

func (m *mockIBKRClient) Ping(ctx context.Context) error { return nil }
func (m *mockIBKRClient) AuthStatus(ctx context.Context) (*ibkr.AuthStatusResponse, error) {
	return nil, nil
}
func (m *mockIBKRClient) Reauthenticate(ctx context.Context) error                { return nil }
func (m *mockIBKRClient) GetAccounts(ctx context.Context) ([]ibkr.Account, error) { return nil, nil }

func (m *mockIBKRClient) PlaceOrder(ctx context.Context, req *ibkr.PlaceOrderRequest) (*ibkr.OrderResponse, error) {
	if m.placeOrderFunc != nil {
		return m.placeOrderFunc(ctx, req)
	}
	return &ibkr.OrderResponse{OrderID: "12345"}, nil
}

func (m *mockIBKRClient) ModifyOrder(ctx context.Context, orderID string, req *ibkr.ModifyOrderRequest) (*ibkr.OrderResponse, error) {
	if m.modifyOrderFunc != nil {
		return m.modifyOrderFunc(ctx, orderID, req)
	}
	return &ibkr.OrderResponse{OrderID: orderID}, nil
}

func (m *mockIBKRClient) CancelOrder(ctx context.Context, orderID string) error {
	if m.cancelOrderFunc != nil {
		return m.cancelOrderFunc(ctx, orderID)
	}
	return nil
}

func (m *mockIBKRClient) GetLiveOrders(ctx context.Context) ([]ibkr.Order, error) {
	if m.getLiveOrdersFunc != nil {
		return m.getLiveOrdersFunc(ctx)
	}
	return []ibkr.Order{{OrderID: "123", Status: "Filled"}}, nil
}

func TestOrderServiceHandler_PlaceOrder(t *testing.T) {
	mock := &mockIBKRClient{}
	handler := NewOrderServiceHandler(mock)

	ctx := context.WithValue(context.Background(), middleware.ClientIdentityContextKey{}, "U12345")
	req := connect.NewRequest(&orderv1.PlaceOrderRequest{
		Symbol:    "AAPL",
		OrderType: orderv1.OrderType_ORDER_TYPE_MARKET,
		Side:      orderv1.OrderSide_ORDER_SIDE_BUY,
		Quantity:  100,
	})

	resp, err := handler.PlaceOrder(ctx, req)
	if err != nil {
		t.Fatalf("PlaceOrder() error = %v", err)
	}
	if resp.Msg.OrderId == "" {
		t.Error("Expected non-empty order ID")
	}
}

func TestOrderServiceHandler_ModifyOrder(t *testing.T) {
	mock := &mockIBKRClient{}
	handler := NewOrderServiceHandler(mock)

	ctx := context.WithValue(context.Background(), middleware.ClientIdentityContextKey{}, "U12345")
	req := connect.NewRequest(&orderv1.ModifyOrderRequest{
		OrderId:  "12345",
		Quantity: 200,
	})

	resp, err := handler.ModifyOrder(ctx, req)
	if err != nil {
		t.Fatalf("ModifyOrder() error = %v", err)
	}
	if resp.Msg.OrderId != "12345" {
		t.Errorf("OrderID = %v, want 12345", resp.Msg.OrderId)
	}
}

func TestOrderServiceHandler_CancelOrder(t *testing.T) {
	mock := &mockIBKRClient{}
	handler := NewOrderServiceHandler(mock)

	ctx := context.WithValue(context.Background(), middleware.ClientIdentityContextKey{}, "U12345")
	req := connect.NewRequest(&orderv1.CancelOrderRequest{
		OrderId: "12345",
	})

	_, err := handler.CancelOrder(ctx, req)
	if err != nil {
		t.Errorf("CancelOrder() error = %v", err)
	}
}

func TestOrderServiceHandler_ListOrders(t *testing.T) {
	mock := &mockIBKRClient{}
	handler := NewOrderServiceHandler(mock)

	ctx := context.WithValue(context.Background(), middleware.ClientIdentityContextKey{}, "U12345")
	req := connect.NewRequest(&orderv1.ListOrdersRequest{})

	resp, err := handler.ListOrders(ctx, req)
	if err != nil {
		t.Fatalf("ListOrders() error = %v", err)
	}
	if len(resp.Msg.Orders) == 0 {
		t.Error("Expected at least one order")
	}
}
