package utils

import (
	"fmt"
	"math/rand/v2"

	"github.com/shopspring/decimal"
)

type numberImpl struct {
}

type INumber interface {
	// RandomInt64 generates a random integer within the specified range
	//
	// start and end support both positive and negative values
	//
	// Result range ∈ [start, end)
	//
	// # Note
	//
	// Panics if start is greater than end
	//
	// # Example:
	//
	// result := RandomInt64(10, 20)
	// //-> 13
	//
	// result := RandomInt64(-10, 20)
	// //-> 3
	//
	// result := RandomInt64(-20, -10)
	// //-> -7
	RandomInt64(start, end int64) int64
	// RandomFloat64 generates a random floating-point number within the specified range
	//
	// start and end support both positive and negative values
	//
	// precision limits the maximum number of decimal places in the return value
	//
	// Result range ∈ [start, end)
	//
	// # Note
	//
	// Panics if start is greater than end
	//
	// # Example:
	//
	// result := RandomFloat64(10.10, 20.20, 2)
	// //-> 16.22
	//
	// result := RandomFloat64(-10.10, 20.20, 3)
	// //-> -7.222
	//
	// result := RandomFloat64(-20.20, -10.10101010101, 4)
	// //-> -8.1234
	RandomFloat64(start, end float64, precision int) float64
	// Pow calculates x raised to the power of y
	//
	// # Note
	//
	// Panics if y is less than 0
	Pow(x, y int64) int64
}

var nbi INumber = &numberImpl{}

// Number instantiates number utilities
func Number() INumber {
	return nbi
}

// RandomInt64 generates a random integer within the specified range
//
// start and end support both positive and negative values
//
// Result range ∈ [start, end)
//
// # Note
//
// Panics if start is greater than end
//
// # Example:
//
// result := RandomInt64(10, 20)
// //-> 13
//
// result := RandomInt64(-10, 20)
// //-> 3
//
// result := RandomInt64(-20, -10)
// //-> -7
func (numberImpl) RandomInt64(start, end int64) int64 {
	if start > end {
		panic(fmt.Errorf("range invalid: start great than end"))
	}
	if start == end {
		return start
	}
	//If the entire range is negative
	if start < 0 && end < 0 {
		fixedStart, fixedEnd := 0-start, 0-end
		return 0 - (fixedEnd + rand.Int64N(fixedStart))
	}
	//If one bound is positive and the other negative
	if start < 0 {
		fixed := 0 - start
		return rand.Int64N(fixed+end) - fixed
	}
	//If both bounds are positive
	if start > 0 && end > 0 {
		return rand.Int64N(end-start) + start
	}
	//Start is 0
	if start == 0 {
		return rand.Int64N(end)
	}

	return start
}

// RandomFloat64 generates a random floating-point number within the specified range
//
// start and end support both positive and negative values
//
// precision limits the maximum number of decimal places in the return value
//
// Result range ∈ [start, end)
//
// # Note
//
// Panics if start is greater than end
//
// # Example:
//
// result := RandomFloat64(10.10, 20.20, 2)
// //-> 16.22
//
// result := RandomFloat64(-10.10, 20.20, 3)
// //-> -7.222
//
// result := RandomFloat64(-20.20, -10.10101010101, 4)
// //-> -8.1234
func (numberImpl) RandomFloat64(start, end float64, precision int) float64 {
	if start > end {
		panic(fmt.Errorf("range invalid: start great than end"))
	}
	if start == end {
		return start
	}

	delta := end - start
	result := start + rand.Float64()*delta

	return decimal.NewFromFloat(result).
		Truncate(int32(precision)).InexactFloat64()
}

// Pow calculates x raised to the power of y
//
// # Note
//
// Panics if y is less than 0
func (numberImpl) Pow(x, y int64) int64 {
	if y < 0 {
		panic(fmt.Errorf("y less than zero"))
	}
	if y == 0 {
		return 1
	}
	if y == 1 {
		return x
	}
	var times int64
	for {
		times++
		if times == y {
			break
		}
		x *= x
	}

	return x
}
