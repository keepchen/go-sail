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
	// Email validates the email address format
	Email(email string) bool
	// IdentityCard validates the (mainland China) identity card format
	//
	// Calculation rules reference the Standardization Administration of China
	//
	// Official documentation @see http://www.gb688.cn/bzgk/gb/newGbInfo?hcno=080D6FBF2BB468F9007657F26D60013E
	IdentityCard(idCard string) bool
}

var vi IValidator = &validatorImpl{}

// Validator instantiates validator utilities
func Validator() IValidator {
	return vi
}

// Email validates the email address format
func (validatorImpl) Email(email string) bool {
	reg, err := regexp.Compile("^(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)])$")
	if err != nil {
		return false
	}

	return reg.Match([]byte(email))
}

// IdentityCard validates the (mainland China) identity card format
//
// Calculation rules reference the Standardization Administration of China
//
// Official documentation @see http://www.gb688.cn/bzgk/gb/newGbInfo?hcno=080D6FBF2BB468F9007657F26D60013E
func (validatorImpl) IdentityCard(idCard string) bool {
	//Mapping from a1 to check digits, where key is a1, value is the check digit, and 10 represents check digit X
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
				//Calculate weighting factors
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
