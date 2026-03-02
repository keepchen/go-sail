package utils

import (
	"fmt"
	"math/rand"

	"github.com/shopspring/decimal"
)

type numberImpl struct {
}

type INumber interface {
	// RandomInt64 在指定范围内取随机整数
	//
	// start和end同时支持正负数
	//
	// 结果值区间 ∈ [start, end)
	//
	// # Note
	//
	// 若start大于end将panic
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
	// RandomFloat64 在指定范围内取随机浮点数
	//
	// start和end同时支持正负数
	//
	// precision为精度，此参数将限定返回值的最大小数位数
	//
	// 结果值区间 ∈ [start, end)
	//
	// # Note
	//
	// 若start大于end将panic
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
	// Pow 计算x的y次幂
	//
	// # Note
	//
	// 若y小于0,将panic
	Pow(x, y int64) int64
}

var nbi INumber = &numberImpl{}

// Number 实例化number工具类
func Number() INumber {
	return nbi
}

// RandomInt64 在指定范围内取随机整数
//
// start和end同时支持正负数
//
// 结果值区间 ∈ [start, end)
//
// # Note
//
// 若start大于end将panic
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
	//如果范围都是负值区间
	if start < 0 && end < 0 {
		fixedStart, fixedEnd := 0-start, 0-end
		return 0 - (fixedEnd + rand.Int63n(fixedStart))
	}
	//如果是一正一负
	if start < 0 {
		fixed := 0 - start
		return rand.Int63n(fixed+end) - fixed
	}
	//如果都为正
	if start > 0 && end > 0 {
		return rand.Int63n(end-start) + start
	}
	//起始为0
	if start == 0 {
		return rand.Int63n(end)
	}

	return start
}

// RandomFloat64 在指定范围内取随机浮点数
//
// start和end同时支持正负数
//
// precision为精度，此参数将限定返回值的最大小数位数
//
// 结果值区间 ∈ [start, end)
//
// # Note
//
// 若start大于end将panic
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

// Pow 计算x的y次幂
//
// # Note
//
// 若y小于0,将panic
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
