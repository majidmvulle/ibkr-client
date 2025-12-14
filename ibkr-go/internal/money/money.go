package money

import (
	"fmt"
	"math"
	"math/big"

	moneyv1 "github.com/majidmvulle/ibkr-client/proto/gen/go/api/common/money/v1"
)

const (
	quadPrecision = 128
	half          = 0.5
	nanoFactorBig = 1_000_000_000_000_000_000 // 10^18
)

var (
	bigNanoFactorBig      = big.NewInt(nanoFactorBig)
	bigFloatNanoFactorBig = new(big.Float).SetPrec(quadPrecision).SetInt64(nanoFactorBig)
	bigFloatHalf          = new(big.Float).SetPrec(quadPrecision).SetFloat64(half)
)

// FromFloat64 converts a float64 amount to a Money proto message using big.Float for precision.
// Inspired by Float64ToBigDecimal pattern for robust decimal conversion.
func FromFloat64(floatVal float64, currency string) (*moneyv1.Money, error) {
	if math.IsNaN(floatVal) || math.IsInf(floatVal, 0) {
		return nil, fmt.Errorf("float64 value %v is NaN or Inf", floatVal)
	}

	// Use big.Float for precise conversion.
	bf := new(big.Float).SetPrec(quadPrecision).SetFloat64(floatVal)
	bfScaled := new(big.Float).Mul(bf, bigFloatNanoFactorBig)
	bfRounded := new(big.Float).SetPrec(quadPrecision)

	// Round towards zero or away from zero depending on sign.
	if bfScaled.Signbit() {
		bfRounded.Sub(bfScaled, bigFloatHalf)
	} else {
		bfRounded.Add(bfScaled, bigFloatHalf)
	}

	// Convert to big.Int.
	biRounded, _ := bfRounded.Int(nil)

	// Split into units and nanos.
	bq := new(big.Int)
	br := new(big.Int)
	bq.QuoRem(biRounded, bigNanoFactorBig, br)

	// Check for overflow.
	if !bq.IsInt64() {
		return nil, fmt.Errorf("float64 value %v results in units overflow for int64", floatVal)
	}

	if !br.IsInt64() {
		return nil, fmt.Errorf("float64 value %v results in nanos overflow for int64", floatVal)
	}

	units := bq.Int64()
	nanos := br.Int64()

	// Verify sign consistency.
	if (units > 0 && nanos < 0) || (units < 0 && nanos > 0) {
		return nil, fmt.Errorf("sign inconsistency for %v -> units=%d, nanos=%d", floatVal, units, nanos)
	}

	return &moneyv1.Money{
		CurrencyCode: currency,
		Units:        units,
		Nanos:        nanos,
	}, nil
}

// ToFloat64 converts a Money proto message to a float64.
func ToFloat64(m *moneyv1.Money) float64 {
	if m == nil {
		return 0
	}

	bfUnits := new(big.Float).SetInt64(m.Units)
	bfNanos := new(big.Float).SetInt64(m.Nanos)
	bfScaledNanos := new(big.Float).Quo(bfNanos, bigFloatNanoFactorBig)
	resultFloat := new(big.Float).Add(bfUnits, bfScaledNanos)
	f64, _ := resultFloat.Float64()

	return f64
}
