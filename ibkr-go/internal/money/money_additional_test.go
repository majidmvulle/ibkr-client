package money

import (
	"testing"
)

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		currency string
		want     float64
	}{
		{"positive amount", 12345, "USD", 123.45},
		{"negative amount", -12345, "USD", -123.45},
		{"zero amount", 0, "USD", 0.0},
		{"large amount", 1000000, "USD", 10000.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money := &Money{
				Amount:   tt.amount,
				Currency: tt.currency,
			}
			got := ToFloat64(money)
			if got != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromFloat64_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		currency string
		wantErr  bool
	}{
		{"very large positive", 999999999.99, "USD", false},
		{"very large negative", -999999999.99, "USD", false},
		{"very small positive", 0.01, "USD", false},
		{"very small negative", -0.01, "USD", false},
		{"zero", 0.0, "USD", false},
		{"empty currency", 100.0, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FromFloat64(tt.value, tt.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromFloat64() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
