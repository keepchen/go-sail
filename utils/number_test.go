package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberImplRandomInt64(t *testing.T) {
	var ranges = [][]int64{
		{0, 0},
		{1, 1},
		{0, 1000},
		{-10, 0},
		{-10, 100},
		{-100, 100},
	}
	for _, rn := range ranges {
		t.Log("range: from", rn[0], "to", rn[1])
		for i := 0; i < 10; i++ {
			t.Log(Number().RandomInt64(rn[0], rn[1]))
		}

		for j := 0; j < 1000; j++ {
			r := Number().RandomInt64(rn[0], rn[1])
			assert.GreaterOrEqual(t, r, rn[0])
			if rn[0] != rn[1] {
				assert.Less(t, r, rn[1])
			} else {
				assert.Equal(t, r, rn[1])
			}
		}
	}
}

func TestNumberImplRandomFloat64(t *testing.T) {
	type target struct {
		start     float64
		end       float64
		precision int
	}
	var ranges = []target{
		{0, 0, 2},
		{1, 1, 2},
		{0, 1000, 3},
		{-10.11, 0.11, 4},
		{-10.101, 100.10111111, 5},
		{-100.1011, 100.1011111111, 6},
		{-200, -100, 6},
	}
	for _, rn := range ranges {
		t.Log("range: from", rn.start, "to", rn.end, "precision", rn.precision)
		for i := 0; i < 10; i++ {
			t.Log(Number().RandomFloat64(rn.start, rn.end, rn.precision))
		}

		for j := 0; j < 1000; j++ {
			r := Number().RandomFloat64(rn.start, rn.end, rn.precision)
			assert.GreaterOrEqual(t, r, rn.start)
			if rn.start != rn.end {
				assert.Less(t, r, rn.end)
			} else {
				assert.Equal(t, r, rn.end)
			}
		}
	}
}

func TestNumberImplPow(t *testing.T) {
	var cases = [][]int64{
		{1, 0, 1},
		{-1, 0, 1},
		{2, 2, 4},
		{3, 2, 9},
	}
	for _, ca := range cases {
		assert.Equal(t, Number().Pow(ca[0], ca[1]), ca[2])
	}
}
