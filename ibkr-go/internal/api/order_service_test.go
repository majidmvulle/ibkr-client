package api

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	orderv1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/order/v1"
	"github.com/stretchr/testify/mock"
)

func TestPlaceOrder(t *testing.T) {
	mockClient := new(MockOrderClient)
	handler := NewOrderServiceHandler(mockClient)

	// Setup context with account ID
	ctx := middleware.SetAccountIDInContext(context.Background(), "U12345")

	// Setup request
	req := connect.NewRequest(&orderv1.PlaceOrderRequest{
		Symbol:      "AAPL",
		Side:        orderv1.OrderSide_ORDER_SIDE_BUY,
		Type:        orderv1.OrderType_ORDER_TYPE_MARKET,
		Quantity:    10,
		TimeInForce: orderv1.TimeInForce_TIME_IN_FORCE_DAY,
	})

	// Mock behavior
	mockClient.On("PlaceOrder", ctx, mock.Anything).Return(&ibkr.OrderResponse{
		OrderID:     "1001",
		OrderStatus: "Submitted",
	}, nil)

	resp, err := handler.PlaceOrder(ctx, req)
	if err != nil {
		t.Fatalf("PlaceOrder() error = %v", err)
	}

	if resp.Msg.OrderId != "1001" {
		t.Errorf("OrderID = %v, want 1001", resp.Msg.OrderId)
	}
}

func TestPlaceOrder_NoAccount(t *testing.T) {
	mockClient := new(MockOrderClient)
	handler := NewOrderServiceHandler(mockClient)

	ctx := context.Background() // No account ID
	req := connect.NewRequest(&orderv1.PlaceOrderRequest{})

	_, err := handler.PlaceOrder(ctx, req)
	if err == nil {
		t.Error("Expected error when no account ID in context")
	}
	if connect.CodeOf(err) != connect.CodeUnauthenticated {
		t.Errorf("Code = %v, want Unauthenticated", connect.CodeOf(err))
	}
}

func TestGetOrder(t *testing.T) {
	mockClient := new(MockOrderClient)
	handler := NewOrderServiceHandler(mockClient)

	ctx := middleware.SetAccountIDInContext(context.Background(), "U12345")
	req := connect.NewRequest(&orderv1.GetOrderRequest{OrderId: "1001"})

	orders := []ibkr.Order{
		{OrderID: "1001", Status: "Submitted", Ticker: "AAPL"},
	}

	mockClient.On("GetLiveOrders", ctx).Return(orders, nil)

	resp, err := handler.GetOrder(ctx, req)
	if err != nil {
		t.Fatalf("GetOrder() error = %v", err)
	}

	if resp.Msg.Order.OrderId != "1001" {
		t.Errorf("OrderId = %v, want 1001", resp.Msg.Order.OrderId)
	}
}

func TestCancelOrder(t *testing.T) {
	mockClient := new(MockOrderClient)
	handler := NewOrderServiceHandler(mockClient)

	ctx := middleware.SetAccountIDInContext(context.Background(), "U12345")
	req := connect.NewRequest(&orderv1.CancelOrderRequest{OrderId: "1001"})

	mockClient.On("CancelOrder", ctx, "1001").Return(nil)

	resp, err := handler.CancelOrder(ctx, req)
	if err != nil {
		t.Fatalf("CancelOrder() error = %v", err)
	}

	if resp.Msg.Status != orderv1.OrderStatus_ORDER_STATUS_CANCELLED {
		t.Errorf("Status = %v, want CANCELLED", resp.Msg.Status)
	}
}
