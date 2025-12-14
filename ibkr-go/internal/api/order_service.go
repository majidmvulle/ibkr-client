package api

import (
	"context"
	"fmt"
	"strconv"

	"connectrpc.com/connect"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/middleware"
	orderv1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/order/v1"
	"github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/order/v1/orderv1connect"
)

const (
	// IBKR order type constants.
	ibkrOrderTypeMarket    = "MKT"
	ibkrOrderTypeLimit     = "LMT"
	ibkrOrderTypeStop      = "STP"
	ibkrOrderTypeStopLimit = "STP LMT"

	// IBKR order side constants.
	ibkrOrderSideBuy  = "BUY"
	ibkrOrderSideSell = "SELL"

	// IBKR time in force constants.
	ibkrTifDay = "DAY"
	ibkrTifGTC = "GTC"
	ibkrTifIOC = "IOC"
	ibkrTifFOK = "FOK"
)

// OrderServiceHandler implements the OrderService ConnectRPC service.
type OrderServiceHandler struct {
	ibkrClient *ibkr.Client
}

// NewOrderServiceHandler creates a new OrderService handler.
func NewOrderServiceHandler(ibkrClient *ibkr.Client) orderv1connect.OrderServiceHandler {
	return &OrderServiceHandler{
		ibkrClient: ibkrClient,
	}
}

// PlaceOrder places a new order.
func (h *OrderServiceHandler) PlaceOrder(
	ctx context.Context,
	req *connect.Request[orderv1.PlaceOrderRequest],
) (*connect.Response[orderv1.PlaceOrderResponse], error) {
	// Get account ID from context (set by session middleware).
	accountID, ok := middleware.GetAccountIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("account ID not found in context"))
	}

	// Map proto request to IBKR request.
	ibkrReq := &ibkr.PlaceOrderRequest{
		SecType:   "STK", // Default to stock, could be enhanced to support other types.
		OrderType: mapOrderType(req.Msg.Type),
		Side:      mapOrderSide(req.Msg.Side),
		Quantity:  req.Msg.Quantity,
		Tif:       mapTimeInForce(req.Msg.TimeInForce),
		Ticker:    req.Msg.Symbol,
	}

	// Set price fields based on order type.
	if req.Msg.LimitPrice != nil {
		ibkrReq.Price = *req.Msg.LimitPrice
	}

	// Place order via IBKR Gateway.
	resp, err := h.ibkrClient.PlaceOrder(ctx, ibkrReq)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to place order: %w", err))
	}

	// Map IBKR response to proto response.
	protoResp := &orderv1.PlaceOrderResponse{
		OrderId: resp.OrderID,
		Status:  mapOrderStatus(resp.OrderStatus),
		Message: formatMessages(resp.Message),
	}

	_ = accountID // Will be used for logging/auditing.

	return connect.NewResponse(protoResp), nil
}

// ModifyOrder modifies an existing order.
func (h *OrderServiceHandler) ModifyOrder(
	ctx context.Context,
	req *connect.Request[orderv1.ModifyOrderRequest],
) (*connect.Response[orderv1.ModifyOrderResponse], error) {
	// Get account ID from context.
	accountID, ok := middleware.GetAccountIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("account ID not found in context"))
	}

	// Map proto request to IBKR request.
	ibkrReq := &ibkr.ModifyOrderRequest{}

	if req.Msg.Quantity != nil {
		ibkrReq.Quantity = *req.Msg.Quantity
	}

	if req.Msg.LimitPrice != nil {
		ibkrReq.Price = *req.Msg.LimitPrice
	}

	// Modify order via IBKR Gateway.
	resp, err := h.ibkrClient.ModifyOrder(ctx, req.Msg.OrderId, ibkrReq)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to modify order: %w", err))
	}

	// Map IBKR response to proto response.
	protoResp := &orderv1.ModifyOrderResponse{
		OrderId: resp.OrderID,
		Status:  mapOrderStatus(resp.OrderStatus),
		Message: formatMessages(resp.Message),
	}

	_ = accountID

	return connect.NewResponse(protoResp), nil
}

// CancelOrder cancels an existing order.
func (h *OrderServiceHandler) CancelOrder(
	ctx context.Context,
	req *connect.Request[orderv1.CancelOrderRequest],
) (*connect.Response[orderv1.CancelOrderResponse], error) {
	// Get account ID from context.
	accountID, ok := middleware.GetAccountIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("account ID not found in context"))
	}

	// Cancel order via IBKR Gateway.
	if err := h.ibkrClient.CancelOrder(ctx, req.Msg.OrderId); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to cancel order: %w", err))
	}

	// Return success response.
	protoResp := &orderv1.CancelOrderResponse{
		OrderId: req.Msg.OrderId,
		Status:  orderv1.OrderStatus_ORDER_STATUS_CANCELLED,
		Message: "Order cancelled successfully",
	}

	_ = accountID

	return connect.NewResponse(protoResp), nil
}

// GetOrder retrieves order details.
func (h *OrderServiceHandler) GetOrder(
	ctx context.Context,
	req *connect.Request[orderv1.GetOrderRequest],
) (*connect.Response[orderv1.GetOrderResponse], error) {
	// Get account ID from context.
	accountID, ok := middleware.GetAccountIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("account ID not found in context"))
	}

	// Get live orders from IBKR Gateway.
	orders, err := h.ibkrClient.GetLiveOrders(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get orders: %w", err))
	}

	// Find the requested order.
	for _, order := range orders {
		if order.OrderID == req.Msg.OrderId {
			protoOrder := mapIBKROrderToProto(&order)

			return connect.NewResponse(&orderv1.GetOrderResponse{
				Order: protoOrder,
			}), nil
		}
	}

	_ = accountID

	return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("order not found"))
}

// ListOrders lists orders for an account.
func (h *OrderServiceHandler) ListOrders(
	ctx context.Context,
	req *connect.Request[orderv1.ListOrdersRequest],
) (*connect.Response[orderv1.ListOrdersResponse], error) {
	// Get account ID from context.
	accountID, ok := middleware.GetAccountIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("account ID not found in context"))
	}

	// Get live orders from IBKR Gateway.
	orders, err := h.ibkrClient.GetLiveOrders(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get orders: %w", err))
	}

	// Filter and map orders.
	protoOrders := filterAndMapOrders(orders, req.Msg.StatusFilter, req.Msg.Limit)

	_ = accountID

	return connect.NewResponse(&orderv1.ListOrdersResponse{
		Orders: protoOrders,
	}), nil
}

// filterAndMapOrders filters orders by status and applies limit.
func filterAndMapOrders(orders []ibkr.Order, statusFilter *orderv1.OrderStatus, limit *int32) []*orderv1.Order {
	protoOrders := make([]*orderv1.Order, 0, len(orders))

	for i := range orders {
		// Filter by status if requested.
		if statusFilter != nil {
			orderStatus := mapOrderStatusFromString(orders[i].Status)
			if orderStatus != *statusFilter {
				continue
			}
		}

		protoOrders = append(protoOrders, mapIBKROrderToProto(&orders[i]))

		// Apply limit if specified.
		if limit != nil && len(protoOrders) >= int(*limit) {
			break
		}
	}

	return protoOrders
}

// Helper functions for mapping between proto and IBKR types.

func mapOrderType(protoType orderv1.OrderType) string {
	switch protoType {
	case orderv1.OrderType_ORDER_TYPE_UNSPECIFIED:
		return ibkrOrderTypeMarket
	case orderv1.OrderType_ORDER_TYPE_MARKET:
		return ibkrOrderTypeMarket
	case orderv1.OrderType_ORDER_TYPE_LIMIT:
		return ibkrOrderTypeLimit
	case orderv1.OrderType_ORDER_TYPE_STOP:
		return ibkrOrderTypeStop
	case orderv1.OrderType_ORDER_TYPE_STOP_LIMIT:
		return ibkrOrderTypeStopLimit
	default:
		return ibkrOrderTypeMarket
	}
}

func mapOrderSide(protoSide orderv1.OrderSide) string {
	switch protoSide {
	case orderv1.OrderSide_ORDER_SIDE_UNSPECIFIED:
		return ibkrOrderSideBuy
	case orderv1.OrderSide_ORDER_SIDE_BUY:
		return ibkrOrderSideBuy
	case orderv1.OrderSide_ORDER_SIDE_SELL:
		return ibkrOrderSideSell
	default:
		return ibkrOrderSideBuy
	}
}

func mapTimeInForce(protoTif orderv1.TimeInForce) string {
	switch protoTif {
	case orderv1.TimeInForce_TIME_IN_FORCE_UNSPECIFIED:
		return ibkrTifDay
	case orderv1.TimeInForce_TIME_IN_FORCE_DAY:
		return ibkrTifDay
	case orderv1.TimeInForce_TIME_IN_FORCE_GTC:
		return ibkrTifGTC
	case orderv1.TimeInForce_TIME_IN_FORCE_IOC:
		return ibkrTifIOC
	case orderv1.TimeInForce_TIME_IN_FORCE_FOK:
		return ibkrTifFOK
	default:
		return ibkrTifDay
	}
}

func mapOrderStatus(ibkrStatus string) orderv1.OrderStatus {
	switch ibkrStatus {
	case "Submitted", "PreSubmitted":
		return orderv1.OrderStatus_ORDER_STATUS_SUBMITTED
	case "Filled":
		return orderv1.OrderStatus_ORDER_STATUS_FILLED
	case "PartiallyFilled":
		return orderv1.OrderStatus_ORDER_STATUS_PARTIALLY_FILLED
	case "Cancelled":
		return orderv1.OrderStatus_ORDER_STATUS_CANCELLED
	case "Inactive", "PendingSubmit":
		return orderv1.OrderStatus_ORDER_STATUS_PENDING
	default:
		return orderv1.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

func mapOrderStatusFromString(status string) orderv1.OrderStatus {
	return mapOrderStatus(status)
}

func formatMessages(messages []string) string {
	if len(messages) == 0 {
		return ""
	}

	result := messages[0]
	for i := 1; i < len(messages); i++ {
		result += "; " + messages[i]
	}

	return result
}

func mapIBKROrderToProto(ibkrOrder *ibkr.Order) *orderv1.Order {
	order := &orderv1.Order{
		OrderId:        ibkrOrder.OrderID,
		AccountId:      ibkrOrder.AcctID,
		Symbol:         ibkrOrder.Ticker,
		Side:           mapOrderSideFromString(ibkrOrder.Side),
		Type:           mapOrderTypeFromString(ibkrOrder.OrigOrderType),
		Quantity:       ibkrOrder.TotalSize,
		FilledQuantity: ibkrOrder.FilledQuantity,
		Status:         mapOrderStatus(ibkrOrder.Status),
		TimeInForce:    orderv1.TimeInForce_TIME_IN_FORCE_DAY, // Default, IBKR doesn't return this.
	}

	if ibkrOrder.Price > 0 {
		order.LimitPrice = &ibkrOrder.Price
	}

	return order
}

func mapOrderSideFromString(side string) orderv1.OrderSide {
	switch side {
	case ibkrOrderSideBuy, "B":
		return orderv1.OrderSide_ORDER_SIDE_BUY
	case ibkrOrderSideSell, "S":
		return orderv1.OrderSide_ORDER_SIDE_SELL
	default:
		return orderv1.OrderSide_ORDER_SIDE_UNSPECIFIED
	}
}

func mapOrderTypeFromString(orderType string) orderv1.OrderType {
	switch orderType {
	case ibkrOrderTypeMarket:
		return orderv1.OrderType_ORDER_TYPE_MARKET
	case ibkrOrderTypeLimit:
		return orderv1.OrderType_ORDER_TYPE_LIMIT
	case ibkrOrderTypeStop:
		return orderv1.OrderType_ORDER_TYPE_STOP
	case ibkrOrderTypeStopLimit:
		return orderv1.OrderType_ORDER_TYPE_STOP_LIMIT
	default:
		return orderv1.OrderType_ORDER_TYPE_UNSPECIFIED
	}
}

// parseConID parses a contract ID from a string (helper for future use).
func parseConID(conIDStr string) (int, error) {
	conID, err := strconv.Atoi(conIDStr)
	if err != nil {
		return 0, fmt.Errorf("invalid contract ID: %w", err)
	}

	return conID, nil
}
