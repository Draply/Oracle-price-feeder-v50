package types

import "cosmossdk.io/math"

type Denom struct {
	Amount math.LegacyDec
	Symbol string
}
