package money

import (
	"math"
	"testing"

	moneyv1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/common/money/v1"
)

func TestFromFloat64(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		currency string
		wantErr  bool
	}{
		{
			name:     "positive USD amount",
			amount:   100.50,
			currency: "USD",
			wantErr:  false,
		},
		{
			name:     "zero amount",
			amount:   0.0,
			currency: "EUR",
			wantErr:  false,
		},
		{
			name:     "negative amount",
			amount:   -50.25,
			currency: "GBP",
			wantErr:  false,
		},
		{
			name:     "large amount",
			amount:   1000000.99,
			currency: "USD",
			wantErr:  false,
		},
		{
			name:     "small fractional amount",
			amount:   0.01,
			currency: "USD",
			wantErr:  false,
		},
		{
			name:     "very small amount",
			amount:   0.001,
			currency: "USD",
			wantErr:  false,
		},
		{
			name:     "amount with many decimal places",
			amount:   123.456789,
			currency: "JPY",
			wantErr:  false,
		},
		{
			name:     "NaN",
			amount:   math.NaN(),
			currency: "USD",
			wantErr:  true,
		},
		{
			name:     "positive infinity",
			amount:   math.Inf(1),
			currency: "USD",
			wantErr:  true,
		},
		{
			name:     "negative infinity",
			amount:   math.Inf(-1),
			currency: "USD",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromFloat64(tt.amount, tt.currency)

			if (err != nil) != tt.wantErr {
				t.Errorf("FromFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got == nil {
				t.Fatal("FromFloat64() returned nil")
			}

			if got.CurrencyCode != tt.currency {
				t.Errorf("FromFloat64() CurrencyCode = %v, want %v", got.CurrencyCode, tt.currency)
			}

			// Verify round-trip conversion
			roundTrip := ToFloat64(got)
			diff := math.Abs(roundTrip - tt.amount)
			if diff > 0.000000001 { // 1 nano precision
				t.Errorf("FromFloat64() round-trip failed: got %v, want %v, diff %v", roundTrip, tt.amount, diff)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name  string
		money *moneyv1.Money
		want  float64
	}{
		{
			name: "positive amount",
			money: &moneyv1.Money{
				Units:        100,
				Nanos:        500000000000000000,
				CurrencyCode: "USD",
			},
			want: 100.50,
		},
		{
			name: "zero amount",
			money: &moneyv1.Money{
				Units:        0,
				Nanos:        0,
				CurrencyCode: "USD",
			},
			want: 0.0,
		},
		{
			name: "negative amount",
			money: &moneyv1.Money{
				Units:        -50,
				Nanos:        -250000000000000000,
				CurrencyCode: "GBP",
			},
			want: -50.25,
		},
		{
			name: "large amount",
			money: &moneyv1.Money{
				Units:        1000000,
				Nanos:        990000000000000000,
				CurrencyCode: "USD",
			},
			want: 1000000.99,
		},
		{
			name: "small fractional amount",
			money: &moneyv1.Money{
				Units:        0,
				Nanos:        10000000000000000,
				CurrencyCode: "USD",
			},
			want: 0.01,
		},
		{
			name:  "nil money",
			money: nil,
			want:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToFloat64(tt.money)

			// Use approximate equality for floating point comparison
			if math.Abs(got-tt.want) > 0.000001 {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		currency string
	}{
		{
			name:     "positive amount",
			amount:   123.45,
			currency: "USD",
		},
		{
			name:     "negative amount",
			amount:   -67.89,
			currency: "EUR",
		},
		{
			name:     "zero",
			amount:   0.0,
			currency: "GBP",
		},
		{
			name:     "large amount",
			amount:   999999.99,
			currency: "JPY",
		},
		{
			name:     "small amount",
			amount:   0.01,
			currency: "USD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to Money
			money, err := FromFloat64(tt.amount, tt.currency)
			if err != nil {
				t.Fatalf("FromFloat64() error = %v", err)
			}

			// Convert back to float64
			got := ToFloat64(money)

			// Verify round-trip (with small tolerance for floating point)
			if math.Abs(got-tt.amount) > 0.000001 {
				t.Errorf("Round-trip failed: got %v, want %v", got, tt.amount)
			}
		})
	}
}

func TestPrecision(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
	}{
		{
			name:   "two decimal places",
			amount: 10.50,
		},
		{
			name:   "three decimal places",
			amount: 10.505,
		},
		{
			name:   "six decimal places",
			amount: 10.505050,
		},
		{
			name:   "nine decimal places",
			amount: 10.505050505,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := FromFloat64(tt.amount, "USD")
			if err != nil {
				t.Fatalf("FromFloat64() error = %v", err)
			}

			got := ToFloat64(money)

			// Verify precision is maintained (within nanos precision)
			diff := math.Abs(got - tt.amount)
			if diff > 0.000000001 { // 1 nano
				t.Errorf("Precision lost: got %v, want %v, diff %v", got, tt.amount, diff)
			}
		})
	}
}

func TestDifferentCurrencies(t *testing.T) {
	currencies := []string{"USD", "EUR", "GBP", "JPY", "CHF", "CAD", "AUD"}
	amount := 100.50

	for _, currency := range currencies {
		t.Run(currency, func(t *testing.T) {
			money, err := FromFloat64(amount, currency)
			if err != nil {
				t.Fatalf("FromFloat64() error = %v", err)
			}

			if money.CurrencyCode != currency {
				t.Errorf("Currency code = %v, want %v", money.CurrencyCode, currency)
			}

			got := ToFloat64(money)
			if math.Abs(got-amount) > 0.000001 {
				t.Errorf("Amount mismatch for %s: got %v, want %v", currency, got, amount)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
	}{
		{
			name:   "very large positive",
			amount: 999999999.99,
		},
		{
			name:   "very large negative",
			amount: -999999999.99,
		},
		{
			name:   "very small positive",
			amount: 0.000000001,
		},
		{
			name:   "very small negative",
			amount: -0.000000001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := FromFloat64(tt.amount, "USD")
			if err != nil {
				t.Fatalf("FromFloat64() error = %v", err)
			}

			got := ToFloat64(money)

			// Verify conversion works without panic
			if money == nil {
				t.Error("FromFloat64() returned nil")
			}

			// Verify reasonable precision
			diff := math.Abs(got - tt.amount)
			if diff > 0.01 { // Allow 1 cent tolerance for very large numbers
				t.Errorf("Large difference: got %v, want %v, diff %v", got, tt.amount, diff)
			}
		})
	}
}

func TestSignConsistency(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
	}{
		{name: "positive", amount: 100.50},
		{name: "negative", amount: -100.50},
		{name: "small positive", amount: 0.01},
		{name: "small negative", amount: -0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := FromFloat64(tt.amount, "USD")
			if err != nil {
				t.Fatalf("FromFloat64() error = %v", err)
			}

			// Verify sign consistency between units and nanos
			if (money.Units > 0 && money.Nanos < 0) || (money.Units < 0 && money.Nanos > 0) {
				t.Errorf("Sign inconsistency: units=%d, nanos=%d", money.Units, money.Nanos)
			}
		})
	}
}
