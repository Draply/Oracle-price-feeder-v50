package oracle_test

import (
	"testing"

	"price-feeder/oracle"
	"price-feeder/oracle/types"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestComputeVWAP(t *testing.T) {
	prices := map[string][]types.TickerPrice{}

	prices["ATOM"] = []types.TickerPrice{{
		Price:  math.LegacyNewDecFromStr("28.21000000"),
		Volume: math.LegacyNewDecFromStr("2749102.78000000"),
	}, {
		Price:  math.LegacyNewDecFromStr("28.268700"),
		Volume: math.LegacyNewDecFromStr("178277.53314385"),
	}, {
		Price:  math.LegacyNewDecFromStr("28.168700"),
		Volume: math.LegacyNewDecFromStr("4749102.53314385"),
	}}

	prices["UMEE"] = []types.TickerPrice{{
		Price:  math.LegacyNewDecFromStr("1.13000000"),
		Volume: math.LegacyNewDecFromStr("249102.38000000"),
	}}

	prices["LUNA"] = []types.TickerPrice{{
		Price:  math.LegacyNewDecFromStr("64.87000000"),
		Volume: math.LegacyNewDecFromStr("7854934.69000000"),
	}, {
		Price:  math.LegacyNewDecFromStr("64.87853000"),
		Volume: math.LegacyNewDecFromStr("458917.46353577"),
	}}

	prices["ZERO1"] = []types.TickerPrice{{
		Price:  math.LegacyNewDecFromStr("12.34000000"),
		Volume: math.LegacyNewDecFromStr("0"),
	}}

	prices["ZERO2"] = []types.TickerPrice{{
		Price:  math.LegacyNewDecFromStr("10"),
		Volume: math.LegacyNewDecFromStr("0"),
	}, {
		Price:  math.LegacyNewDecFromStr("20"),
		Volume: math.LegacyNewDecFromStr("0"),
	}}

	expected := map[string]math.LegacyDec{
		"ATOM":  math.LegacyNewDecFromStr("28.185812745610043621"),
		"UMEE":  math.LegacyNewDecFromStr("1.13000000"),
		"LUNA":  math.LegacyNewDecFromStr("64.870470848638112395"),
		"ZERO1": math.LegacyNewDecFromStr("12.34000000"),
		"ZERO2": math.LegacyNewDecFromStr("15"),
	}

	for denom, tickers := range prices {
		t.Run(denom, func(t *testing.T) {
			vwap, err := oracle.ComputeVWAP(tickers)
			require.NoError(t, err)
			require.Equal(t, expected[denom], vwap)
		})
	}

	// Supply emty list of tickers
	t.Run("EMPTY", func(t *testing.T) {
		_, err := oracle.ComputeVWAP([]types.TickerPrice{})
		require.Error(t, err)
	})
}

func TestStandardDeviation(t *testing.T) {
	type result struct {
		mean      math.LegacyDec
		deviation math.LegacyDec
		err       bool
	}
	restCases := map[string]struct {
		prices   []math.LegacyDec
		expected result
	}{
		"empty prices": {
			prices:   []math.LegacyDec{},
			expected: result{},
		},
		"nil prices": {
			prices:   nil,
			expected: result{},
		},
		"not enough prices": {
			prices: []math.LegacyDec{
				math.LegacyNewDecFromStr("28.21000000"),
				math.LegacyNewDecFromStr("28.23000000"),
			},
			expected: result{},
		},
		"enough prices 1": {
			prices: []math.LegacyDec{
				math.LegacyNewDecFromStr("28.21000000"),
				math.LegacyNewDecFromStr("28.23000000"),
				math.LegacyNewDecFromStr("28.40000000"),
			},
			expected: result{
				mean:      math.LegacyNewDecFromStr("28.28"),
				deviation: math.LegacyNewDecFromStr("0.085244745683629475"),
				err:       false,
			},
		},
		"enough prices 2": {
			prices: []math.LegacyDec{
				math.LegacyNewDecFromStr("1.13000000"),
				math.LegacyNewDecFromStr("1.13050000"),
				math.LegacyNewDecFromStr("1.14000000"),
			},
			expected: result{
				mean:      math.LegacyNewDecFromStr("1.1335"),
				deviation: math.LegacyNewDecFromStr("0.004600724580614015"),
				err:       false,
			},
		},
	}

	for name, test := range restCases {
		test := test

		t.Run(name, func(t *testing.T) {
			deviation, mean, _ := oracle.StandardDeviation(test.prices)
			// if test.expected.err == false {
			// 	require.NoError(t, err)
			// }
			require.Equal(t, test.expected.deviation, deviation)
			require.Equal(t, test.expected.mean, mean)
		})
	}
}
