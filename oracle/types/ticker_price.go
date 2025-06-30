package types

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
)

// TickerPrice defines price and volume information for a symbol or ticker exchange rate.
type TickerPrice struct {
	Price  math.LegacyDec `json:"price"`  // last trade price
	Volume math.LegacyDec `json:"volume"` // 24h volume
	Time   time.Time      `json:"time"`
}

func NewTickerPrice(price string, volume string, timestamp time.Time) (TickerPrice, error) {
	priceDec, err := math.LegacyNewDecFromStr(price)
	if err != nil {
		return TickerPrice{}, fmt.Errorf("failed to convert ticker price: %v", err)
	}
	volumeDec, err := math.LegacyNewDecFromStr(volume)
	if err != nil {
		return TickerPrice{}, fmt.Errorf("failed to convert ticker volume: %v", err)
	}
	ticker := TickerPrice{
		Price:  priceDec,
		Volume: volumeDec,
		Time:   timestamp,
	}
	return ticker, nil
}
