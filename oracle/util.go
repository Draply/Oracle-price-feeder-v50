package oracle

import (
	"fmt"

	"price-feeder/oracle/provider"
	"price-feeder/oracle/types"

	"cosmossdk.io/math"
)

// ComputeVWAP computes the volume weighted average price for all tickers.
// If all tickers report a volume of 0, treat all volumes as 1 and
// effectively return the average price instead.
// Ref: https://en.wikipedia.org/wiki/Volume-weighted_average_price
func ComputeVWAP(tickers []types.TickerPrice) (math.LegacyDec, error) {
	if len(tickers) == 0 {
		return math.LegacyDec{}, fmt.Errorf("no tickers supplied")
	}

	volumeSum := math.LegacyZeroDec()

	for _, tp := range tickers {
		volumeSum = volumeSum.Add(tp.Volume)
	}

	weightedPrice := math.LegacyZeroDec()

	for _, tp := range tickers {
		volume := tp.Volume
		if volumeSum.Equal(math.LegacyZeroDec()) {
			volume = math.LegacyNewDec(1)
		}

		// weightedPrice = Σ {P * V} for all TickerPrice
		weightedPrice = weightedPrice.Add(tp.Price.Mul(volume))
	}

	if volumeSum.Equal(math.LegacyZeroDec()) {
		volumeSum = math.LegacyNewDec(int64(len(tickers)))
	}

	return weightedPrice.Quo(volumeSum), nil
}

// StandardDeviation returns standard deviation and mean of assets.
// Will skip calculating for an asset if there are less than 3 prices.
func StandardDeviation(prices []math.LegacyDec) (math.LegacyDec, math.LegacyDec, error) {
	// Skip if standard deviation would not be meaningful
	if len(prices) < 3 {
		err := fmt.Errorf("not enough values to calculate deviation")
		return math.LegacyDec{}, math.LegacyDec{}, err
	}

	sum := math.LegacyZeroDec()

	for _, price := range prices {
		sum = sum.Add(price)
	}

	numPrices := int64(len(prices))
	mean := sum.QuoInt64(numPrices)
	varianceSum := math.LegacyZeroDec()

	for _, price := range prices {
		deviation := price.Sub(mean)
		varianceSum = varianceSum.Add(deviation.Mul(deviation))
	}

	variance := varianceSum.QuoInt64(numPrices)

	deviation, err := variance.ApproxSqrt()
	if err != nil {
		return math.LegacyDec{}, math.LegacyDec{}, err
	}

	return deviation, mean, nil
}

func SetWeight(
	rates map[provider.Name]types.TickerPrice,
	weight ProviderWeight,
) (
	map[provider.Name]types.TickerPrice,
	error,
) {
	if len(weight.Weight) == 0 {
		return rates, nil
	}

	for name, volume := range weight.Weight {
		providerName := provider.Name(name)
		ticker, found := rates[providerName]
		if found {
			ticker.Volume = volume
			rates[providerName] = ticker
		}
	}

	return rates, nil
}
