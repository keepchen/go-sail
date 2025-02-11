package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type validatorImpl struct {
}

type IValidator interface {
	// Email 验证邮箱格式
	Email(email string) bool
	// IdentityCard （中国大陆）身份证格式校验
	//
	// 计算规则参考“中国国家标准化管理委员会”
	//
	// 官方文档@see http://www.gb688.cn/bzgk/gb/newGbInfo?hcno=080D6FBF2BB468F9007657F26D60013E
	IdentityCard(idCard string) bool
}

var _ IValidator = &validatorImpl{}

// Validator 实例化validator工具类
func Validator() IValidator {
	return &validatorImpl{}
}

// Email 验证邮箱格式
func (validatorImpl) Email(email string) bool {
	reg, err := regexp.Compile("^(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)])$")
	if err != nil {
		return false
	}

	return reg.Match([]byte(email))
}

// IdentityCard （中国大陆）身份证格式校验
//
// 计算规则参考“中国国家标准化管理委员会”
//
// 官方文档@see http://www.gb688.cn/bzgk/gb/newGbInfo?hcno=080D6FBF2BB468F9007657F26D60013E
func (validatorImpl) IdentityCard(idCard string) bool {
	//a1与对应的校验码对照表，其中key表示a1，value表示校验码，value中的10表示校验码X
	var a1Map = map[int]int{
		0:  1,
		1:  0,
		2:  10,
		3:  9,
		4:  8,
		5:  7,
		6:  6,
		7:  5,
		8:  4,
		9:  3,
		10: 2,
	}

	var idStr = strings.ToUpper(idCard)
	var reg, err = regexp.Compile(`^[0-9]{17}[0-9X]$`)
	if err != nil {
		return false
	}
	if !reg.Match([]byte(idStr)) {
		return false
	}
	var sum int
	var signChar = ""
	for index, c := range idStr {
		var i = 18 - index
		if i != 1 {
			if v, pErr := strconv.Atoi(string(c)); pErr == nil {
				//计算加权因子
				var weight = int(math.Pow(2, float64(i-1))) % 11
				sum += v * weight
			} else {
				return false
			}
		} else {
			signChar = string(c)
		}
	}
	var a1 = a1Map[sum%11]
	var a1Str = fmt.Sprintf("%d", a1)
	if a1 == 10 {
		a1Str = "X"
	}

	return a1Str == signChar
}
