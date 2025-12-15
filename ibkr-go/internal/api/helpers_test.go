package api

import (
	"testing"

	"github.com/majidmvulle/ibkr-client/ibkr-go/internal/ibkr"
	orderv1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/order/v1"
)

func TestCalculatePnLPercent(t *testing.T) {
	tests := []struct {
		name          string
		unrealizedPnl float64
		avgCost       float64
		quantity      float64
		want          float64
	}{
		{"positive pnl", 100, 50, 10, 20.0},
		{"negative pnl", -100, 50, 10, -20.0},
		{"zero avg cost", 100, 0, 10, 0},
		{"zero quantity", 100, 50, 0, 0},
		{"zero total cost", 100, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculatePnLPercent(tt.unrealizedPnl, tt.avgCost, tt.quantity)
			if got != tt.want {
				t.Errorf("calculatePnLPercent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapIBKRPositionToProto(t *testing.T) {
	pos := &ibkr.Position{
		ContractDesc:  "AAPL",
		Position:      100,
		MktValue:      15000,
		AvgCost:       140,
		UnrealizedPnl: 1000,
		Currency:      "USD",
	}

	result, err := mapIBKRPositionToProto(pos)
	if err != nil {
		t.Fatalf("mapIBKRPositionToProto() error = %v", err)
	}

	if result.Symbol != "AAPL" {
		t.Errorf("Symbol = %v, want AAPL", result.Symbol)
	}
	if result.Quantity != 100 {
		t.Errorf("Quantity = %v, want 100", result.Quantity)
	}
}

func TestMapSnapshotToQuote(t *testing.T) {
	snapshot := &ibkr.MarketDataSnapshot{
		LastPrice: 150.0,
		Bid:       149.5,
		Ask:       150.5,
		Volume:    1000000,
		High:      151.0,
		Low:       149.0,
	}

	quote := mapSnapshotToQuote(snapshot, "AAPL")

	if quote.Symbol != "AAPL" {
		t.Errorf("Symbol = %v, want AAPL", quote.Symbol)
	}
	if quote.Last != 150.0 {
		t.Errorf("Last = %v, want 150.0", quote.Last)
	}
	if quote.Bid != 149.5 {
		t.Errorf("Bid = %v, want 149.5", quote.Bid)
	}
}

func TestMapHistoricalBarToProto(t *testing.T) {
	bar := &ibkr.HistoricalBar{
		Time:   1234567890,
		Open:   100.0,
		High:   105.0,
		Low:    99.0,
		Close:  103.0,
		Volume: 50000,
	}

	result := mapHistoricalBarToProto(bar)

	if result.Open != 100.0 {
		t.Errorf("Open = %v, want 100.0", result.Open)
	}
	if result.High != 105.0 {
		t.Errorf("High = %v, want 105.0", result.High)
	}
	if result.Close != 103.0 {
		t.Errorf("Close = %v, want 103.0", result.Close)
	}
}

func TestFilterAndMapOrders(t *testing.T) {
	orders := []ibkr.Order{
		{OrderID: "1", Status: "Submitted"},
		{OrderID: "2", Status: "Filled"},
		{OrderID: "3", Status: "Cancelled"},
	}

	// Test with no filter
	result := filterAndMapOrders(orders, nil, nil)
	if len(result) != 3 {
		t.Errorf("filterAndMapOrders() returned %d orders, want 3", len(result))
	}

	// Test with limit
	limit := int32(2)
	result = filterAndMapOrders(orders, nil, &limit)
	if len(result) != 2 {
		t.Errorf("filterAndMapOrders() with limit returned %d orders, want 2", len(result))
	}

	// Test with status filter
	status := orderv1.OrderStatus_ORDER_STATUS_FILLED
	result = filterAndMapOrders(orders, &status, nil)
	if len(result) != 1 {
		t.Errorf("filterAndMapOrders() with status filter returned %d orders, want 1", len(result))
	}
}
