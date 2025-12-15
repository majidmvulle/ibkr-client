package money

import (
	"testing"
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
				if money.Currency != tt.currency {
					t.Errorf("Currency = %v, want %v", money.Currency, tt.currency)
				}
			}
		})
	}
}

func TestToFloat64_AllCases(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		currency string
		want     float64
	}{
		{"positive", 10050, "USD", 100.50},
		{"negative", -5025, "USD", -50.25},
		{"zero", 0, "USD", 0.0},
		{"large", 99999999, "USD", 999999.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money := &Money{Amount: tt.amount, Currency: tt.currency}
			got := ToFloat64(money)
			if got != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}
