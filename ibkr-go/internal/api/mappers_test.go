package api

import (
	"testing"

	orderv1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/ibkr/order/v1"
)

func TestMapOrderType(t *testing.T) {
	tests := []struct {
		name      string
		protoType orderv1.OrderType
		want      string
	}{
		{"market order", orderv1.OrderType_ORDER_TYPE_MARKET, "MKT"},
		{"limit order", orderv1.OrderType_ORDER_TYPE_LIMIT, "LMT"},
		{"stop order", orderv1.OrderType_ORDER_TYPE_STOP, "STP"},
		{"stop limit", orderv1.OrderType_ORDER_TYPE_STOP_LIMIT, "STP LMT"},
		{"unspecified", orderv1.OrderType_ORDER_TYPE_UNSPECIFIED, "MKT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapOrderType(tt.protoType)
			if got != tt.want {
				t.Errorf("mapOrderType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapOrderSide(t *testing.T) {
	tests := []struct {
		name      string
		protoSide orderv1.OrderSide
		want      string
	}{
		{"buy", orderv1.OrderSide_ORDER_SIDE_BUY, "BUY"},
		{"sell", orderv1.OrderSide_ORDER_SIDE_SELL, "SELL"},
		{"unspecified", orderv1.OrderSide_ORDER_SIDE_UNSPECIFIED, "BUY"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapOrderSide(tt.protoSide)
			if got != tt.want {
				t.Errorf("mapOrderSide() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapTimeInForce(t *testing.T) {
	tests := []struct {
		name     string
		protoTif orderv1.TimeInForce
		want     string
	}{
		{"day", orderv1.TimeInForce_TIME_IN_FORCE_DAY, "DAY"},
		{"gtc", orderv1.TimeInForce_TIME_IN_FORCE_GTC, "GTC"},
		{"ioc", orderv1.TimeInForce_TIME_IN_FORCE_IOC, "IOC"},
		{"fok", orderv1.TimeInForce_TIME_IN_FORCE_FOK, "FOK"},
		{"unspecified", orderv1.TimeInForce_TIME_IN_FORCE_UNSPECIFIED, "DAY"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapTimeInForce(tt.protoTif)
			if got != tt.want {
				t.Errorf("mapTimeInForce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapOrderStatus(t *testing.T) {
	tests := []struct {
		name       string
		ibkrStatus string
		want       orderv1.OrderStatus
	}{
		{"submitted", "Submitted", orderv1.OrderStatus_ORDER_STATUS_PENDING},
		{"filled", "Filled", orderv1.OrderStatus_ORDER_STATUS_FILLED},
		{"cancelled", "Cancelled", orderv1.OrderStatus_ORDER_STATUS_CANCELLED},
		{"inactive", "Inactive", orderv1.OrderStatus_ORDER_STATUS_CANCELLED},
		{"unknown", "Unknown", orderv1.OrderStatus_ORDER_STATUS_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapOrderStatus(tt.ibkrStatus)
			if got != tt.want {
				t.Errorf("mapOrderStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapOrderSideFromString(t *testing.T) {
	tests := []struct {
		name string
		side string
		want orderv1.OrderSide
	}{
		{"buy", "BUY", orderv1.OrderSide_ORDER_SIDE_BUY},
		{"sell", "SELL", orderv1.OrderSide_ORDER_SIDE_SELL},
		{"unknown", "UNKNOWN", orderv1.OrderSide_ORDER_SIDE_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapOrderSideFromString(tt.side)
			if got != tt.want {
				t.Errorf("mapOrderSideFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapOrderTypeFromString(t *testing.T) {
	tests := []struct {
		name      string
		orderType string
		want      orderv1.OrderType
	}{
		{"market", "MKT", orderv1.OrderType_ORDER_TYPE_MARKET},
		{"limit", "LMT", orderv1.OrderType_ORDER_TYPE_LIMIT},
		{"stop", "STP", orderv1.OrderType_ORDER_TYPE_STOP},
		{"stop limit", "STP LMT", orderv1.OrderType_ORDER_TYPE_STOP_LIMIT},
		{"unknown", "UNKNOWN", orderv1.OrderType_ORDER_TYPE_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapOrderTypeFromString(tt.orderType)
			if got != tt.want {
				t.Errorf("mapOrderTypeFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatMessages(t *testing.T) {
	tests := []struct {
		name     string
		messages []string
		want     string
	}{
		{"empty", []string{}, ""},
		{"single", []string{"error1"}, "error1"},
		{"multiple", []string{"error1", "error2"}, "error1; error2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatMessages(tt.messages)
			if got != tt.want {
				t.Errorf("formatMessages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseConID(t *testing.T) {
	tests := []struct {
		name    string
		conID   string
		want    int
		wantErr bool
	}{
		{"valid", "12345", 12345, false},
		{"invalid", "abc", 0, true},
		{"empty", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConID(tt.conID)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseConID() = %v, want %v", got, tt.want)
			}
		})
	}
}
