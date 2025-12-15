package money

import (
	"testing"

	moneyv1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/common/money/v1"
)

func TestFromFloat64_AllCases(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		currency string
		wantErr  bool
	}{
		{"positive", 100.50, "USD", false},
		{"negative", -50.25, "USD", false},
		{"zero", 0.0, "USD", false},
		{"large", 999999.99, "USD", false},
		{"small", 0.01, "USD", false},
		{"empty currency", 100.0, "", false},
		{"EUR", 50.0, "EUR", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := FromFloat64(tt.value, tt.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if money == nil {
					t.Error("Expected non-nil money")
				}
				if money.CurrencyCode != tt.currency {
					t.Errorf("CurrencyCode = %v, want %v", money.CurrencyCode, tt.currency)
				}
			}
		})
	}
}

func TestToFloat64_AllCases(t *testing.T) {
	tests := []struct {
		name     string
		units    int64
		nanos    int64
		currency string
		want     float64
	}{
		{"positive", 100, 500000000000000000, "USD", 100.5},
		{"negative", -50, -250000000000000000, "USD", -50.25},
		{"zero", 0, 0, "USD", 0.0},
		{"large", 999999, 990000000000000000, "USD", 999999.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money := &moneyv1.Money{Units: tt.units, Nanos: tt.nanos, CurrencyCode: tt.currency}
			got := ToFloat64(money)
			if got != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}
