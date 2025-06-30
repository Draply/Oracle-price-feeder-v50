package types

import (
	"fmt"

	"cosmossdk.io/math"
)

// CandlePrice defines price, volume, and time information for an exchange rate.
type CandlePrice struct {
	Price     math.LegacyDec // last trade price
	Volume    math.LegacyDec // volume
	TimeStamp int64          // timestamp
}

// NewCandlePrice parses the lastPrice and volume to a decimal and returns a CandlePrice
func NewCandlePrice(provider, symbol, lastPrice, volume string, timeStamp int64) (CandlePrice, error) {
	price, err := math.LegacyNewDecFromStr(lastPrice)
	if err != nil {
		return CandlePrice{}, fmt.Errorf("failed to parse %s price (%s) for %s: %w", provider, lastPrice, symbol, err)
	}

	volumeDec, err := math.LegacyNewDecFromStr(volume)
	if err != nil {
		return CandlePrice{}, fmt.Errorf("failed to parse %s volume (%s) for %s: %w", provider, volume, symbol, err)
	}

	return CandlePrice{Price: price, Volume: volumeDec, TimeStamp: timeStamp}, nil
}
