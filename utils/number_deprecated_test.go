package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomInt64(t *testing.T) {
	var ranges = [][]int64{
		{0, 0},
		{1, 1},
		{0, 1000},
		{100, 1000},
		{-10, 0},
		{-10, 100},
		{-100, 100},
	}
	t.Run("RandomInt64", func(t *testing.T) {
		for _, rn := range ranges {
			t.Log("range: from", rn[0], "to", rn[1])
			for i := 0; i < 10; i++ {
				t.Log(RandomInt64(rn[0], rn[1]))
			}

			for j := 0; j < 1000; j++ {
				r := RandomInt64(rn[0], rn[1])
				assert.GreaterOrEqual(t, r, rn[0])
				if rn[0] != rn[1] {
					assert.Less(t, r, rn[1])
				} else {
					assert.Equal(t, r, rn[1])
				}
			}
		}
	})

	t.Run("RandomInt64-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			t.Log(RandomInt64(1, 0))
		})
	})

	t.Run("RandomInt64-BothNegative", func(t *testing.T) {
		t.Log(RandomInt64(-2, -1))
	})
}

func TestRandomFloat64(t *testing.T) {
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
		{-1000, 100, 2},
		{-1000, 0, 3},
	}
	t.Run("RandomFloat64", func(t *testing.T) {
		for _, rn := range ranges {
			t.Log("range: from", rn.start, "to", rn.end, "precision", rn.precision)
			for i := 0; i < 10; i++ {
				t.Log(RandomFloat64(rn.start, rn.end, rn.precision))
			}

			for j := 0; j < 1000; j++ {
				r := RandomFloat64(rn.start, rn.end, rn.precision)
				assert.GreaterOrEqual(t, r, rn.start)
				if rn.start != rn.end {
					assert.Less(t, r, rn.end)
				} else {
					assert.Equal(t, r, rn.end)
				}
			}
		}
	})

	t.Run("RandomFloat64-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			t.Log(RandomFloat64(100, 10, 4))
		})
	})
}

func TestPow(t *testing.T) {
	var cases = [][]int64{
		{1, 0, 1},
		{-1, 0, 1},
		{2, 2, 4},
		{3, 2, 9},
	}
	t.Run("Pow", func(t *testing.T) {
		for _, ca := range cases {
			assert.Equal(t, Pow(ca[0], ca[1]), ca[2])
		}
	})

	t.Run("Pow-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Pow(1, -1)
		})
	})
}
