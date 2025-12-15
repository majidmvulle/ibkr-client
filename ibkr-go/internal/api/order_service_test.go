package api

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	orderv1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/order/v1"
)

func TestOrderServiceHandler_GetOrder(t *testing.T) {
	mock := &mockIBKRClient{
		getLiveOrdersFunc: func(ctx context.Context) ([]ibkr.Order, error) {
			return []ibkr.Order{{OrderID: "12345", Status: "Filled"}}, nil
		},
	}
	handler := NewOrderServiceHandler(mock)

	ctx := context.WithValue(context.Background(), middleware.ClientIdentityContextKey{}, "U12345")
	req := connect.NewRequest(&orderv1.GetOrderRequest{OrderId: "12345"})

	resp, err := handler.GetOrder(ctx, req)
	if err != nil {
		t.Fatalf("GetOrder() error = %v", err)
	}
	if resp.Msg.Order.OrderId != "12345" {
		t.Errorf("OrderID = %v, want 12345", resp.Msg.Order.OrderId)
	}
}

func TestOrderServiceHandler_NoAccountID(t *testing.T) {
	mock := &mockIBKRClient{}
	handler := NewOrderServiceHandler(mock)

	ctx := context.Background() // No account ID
	req := connect.NewRequest(&orderv1.PlaceOrderRequest{
		Symbol:    "AAPL",
		OrderType: orderv1.OrderType_ORDER_TYPE_MARKET,
		Side:      orderv1.OrderSide_ORDER_SIDE_BUY,
		Quantity:  100,
	})

	_, err := handler.PlaceOrder(ctx, req)
	if err == nil {
		t.Error("Expected error when account ID missing")
	}
}

func TestOrderServiceHandler_IBKRError(t *testing.T) {
	mock := &mockIBKRClient{
		placeOrderFunc: func(ctx context.Context, req *ibkr.PlaceOrderRequest) (*ibkr.OrderResponse, error) {
			return nil, errors.New("IBKR error")
		},
	}
	handler := NewOrderServiceHandler(mock)

	ctx := context.WithValue(context.Background(), middleware.ClientIdentityContextKey{}, "U12345")
	req := connect.NewRequest(&orderv1.PlaceOrderRequest{
		Symbol:    "AAPL",
		OrderType: orderv1.OrderType_ORDER_TYPE_MARKET,
		Side:      orderv1.OrderSide_ORDER_SIDE_BUY,
		Quantity:  100,
	})

	_, err := handler.PlaceOrder(ctx, req)
	if err == nil {
		t.Error("Expected error from IBKR")
	}
}

func TestOrderServiceHandler_ListOrders_WithFilter(t *testing.T) {
	mock := &mockIBKRClient{
		getLiveOrdersFunc: func(ctx context.Context) ([]ibkr.Order, error) {
			return []ibkr.Order{
				{OrderID: "1", Status: "Filled"},
				{OrderID: "2", Status: "Submitted"},
			}, nil
		},
	}
	handler := NewOrderServiceHandler(mock)

	ctx := context.WithValue(context.Background(), middleware.ClientIdentityContextKey{}, "U12345")
	status := orderv1.OrderStatus_ORDER_STATUS_FILLED
	req := connect.NewRequest(&orderv1.ListOrdersRequest{
		Status: &status,
	})

	resp, err := handler.ListOrders(ctx, req)
	if err != nil {
		t.Fatalf("ListOrders() error = %v", err)
	}
	if len(resp.Msg.Orders) == 0 {
		t.Error("Expected filtered orders")
	}
}
