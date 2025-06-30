package volume

import (
	"cosmossdk.io/math"
)

type Total struct {
	Total  math.LegacyDec
	Values int
	First  uint64
}

func NewTotal() *Total {
	return &Total{
		Total:  math.LegacyZeroDec(),
		Values: 0,
	}
}

func (t *Total) Clear() {
	t.Total = math.LegacyZeroDec()
	t.Values = 0
	t.First = 0
}

func (t *Total) Sub(value math.LegacyDec) {
	if value.IsNil() || value.IsNegative() {
		return
	}
	t.Total = t.Total.Sub(value)
	t.Values -= 1
}

func (t *Total) Add(value math.LegacyDec, height uint64) {
	if value.IsNil() || value.IsNegative() {
		return
	}
	t.Total = t.Total.Add(value)
	t.Values += 1
	if height < t.First || t.First == 0 {
		t.First = height
	}
}
