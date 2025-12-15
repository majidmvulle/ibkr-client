package integration

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/api"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	orderv1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/order/v1"
)

func TestIntegration_OrderService_PlaceOrder(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create test session
	token := CreateTestSession(t, testCtx.Config.IBKRAccountID)
	defer DeleteTestSession(t, token)

	// Add account ID to context (simulating session middleware)
	ctx = middleware.SetAccountIDInContext(ctx, testCtx.Config.IBKRAccountID)

	// Create order service handler
	handler := api.NewOrderServiceHandler(testCtx.IBKRClient)

	// Create place order request
	req := connect.NewRequest(&orderv1.PlaceOrderRequest{
		Symbol:      "AAPL",
		Quantity:    100,
		Side:        orderv1.OrderSide_ORDER_SIDE_BUY,
		Type:        orderv1.OrderType_ORDER_TYPE_MARKET,
		TimeInForce: orderv1.TimeInForce_TIME_IN_FORCE_DAY,
	})

	// Place order
	resp, err := handler.PlaceOrder(ctx, req)
	if err != nil {
		t.Fatalf("PlaceOrder failed: %v", err)
	}

	// Verify response
	if resp.Msg.OrderId == "" {
		t.Error("Expected order ID to be set")
	}

	if resp.Msg.Status == orderv1.OrderStatus_ORDER_STATUS_UNSPECIFIED {
		t.Error("Expected order status to be set")
	}

	t.Logf("Order placed successfully: ID=%s, Status=%s", resp.Msg.OrderId, resp.Msg.Status)
}

func TestIntegration_OrderService_ListOrders(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create test session
	token := CreateTestSession(t, testCtx.Config.IBKRAccountID)
	defer DeleteTestSession(t, token)

	// Add account ID to context
	ctx = middleware.SetAccountIDInContext(ctx, testCtx.Config.IBKRAccountID)

	// Create order service handler
	handler := api.NewOrderServiceHandler(testCtx.IBKRClient)

	// Create list orders request
	req := connect.NewRequest(&orderv1.ListOrdersRequest{})

	// List orders
	resp, err := handler.ListOrders(ctx, req)
	if err != nil {
		t.Fatalf("ListOrders failed: %v", err)
	}

	// Verify response
	if resp.Msg.Orders == nil {
		t.Error("Expected orders list to be initialized")
	}

	t.Logf("Listed %d orders", len(resp.Msg.Orders))
}

func TestIntegration_OrderService_CancelOrder(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create test session
	token := CreateTestSession(t, testCtx.Config.IBKRAccountID)
	defer DeleteTestSession(t, token)

	// Add account ID to context
	ctx = middleware.SetAccountIDInContext(ctx, testCtx.Config.IBKRAccountID)

	// Create order service handler
	handler := api.NewOrderServiceHandler(testCtx.IBKRClient)

	// First, place an order
	placeReq := connect.NewRequest(&orderv1.PlaceOrderRequest{
		Symbol:      "AAPL",
		Quantity:    100,
		Side:        orderv1.OrderSide_ORDER_SIDE_BUY,
		Type:        orderv1.OrderType_ORDER_TYPE_LIMIT,
		LimitPrice:  &[]float64{150.00}[0],
		TimeInForce: orderv1.TimeInForce_TIME_IN_FORCE_DAY,
	})

	placeResp, err := handler.PlaceOrder(ctx, placeReq)
	if err != nil {
		t.Fatalf("PlaceOrder failed: %v", err)
	}

	orderID := placeResp.Msg.OrderId

	// Now cancel the order
	cancelReq := connect.NewRequest(&orderv1.CancelOrderRequest{
		OrderId: orderID,
	})

	cancelResp, err := handler.CancelOrder(ctx, cancelReq)
	if err != nil {
		t.Fatalf("CancelOrder failed: %v", err)
	}

	// Verify response
	if cancelResp.Msg.Status == orderv1.OrderStatus_ORDER_STATUS_UNSPECIFIED {
		t.Error("Expected cancel status to be set")
	}

	t.Logf("Order cancelled successfully: ID=%s", orderID)
}

func TestIntegration_OrderService_InvalidSymbol(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create test session
	token := CreateTestSession(t, testCtx.Config.IBKRAccountID)
	defer DeleteTestSession(t, token)

	// Add account ID to context
	ctx = middleware.SetAccountIDInContext(ctx, testCtx.Config.IBKRAccountID)

	// Create order service handler
	handler := api.NewOrderServiceHandler(testCtx.IBKRClient)

	// Create place order request with invalid symbol
	req := connect.NewRequest(&orderv1.PlaceOrderRequest{
		Symbol:      "INVALID_SYMBOL_12345",
		Quantity:    100,
		Side:        orderv1.OrderSide_ORDER_SIDE_BUY,
		Type:        orderv1.OrderType_ORDER_TYPE_MARKET,
		TimeInForce: orderv1.TimeInForce_TIME_IN_FORCE_DAY,
	})

	// Place order (should fail or handle gracefully)
	_, err := handler.PlaceOrder(ctx, req)

	// We expect this might fail with the mock, but it shouldn't panic
	if err != nil {
		t.Logf("Expected error for invalid symbol: %v", err)
	}
}
