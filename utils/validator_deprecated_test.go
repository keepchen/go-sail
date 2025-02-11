package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type emailChallengeDeprecated struct {
	email string
	valid bool
}

var ecsDeprecated = []emailChallengeDeprecated{
	{"example@gmail.com", true},
	{"example@gmail", false},
	{"example@gmail.com.us", true},
	{"example-@gmail.com.us", true},
	{"example-x@gmail.com.us", true},
	{"example_x@gmail.com.us", true},
	{"example.@gmail.com.us", false},
	{"example.x@gmail.com.us", true},
	{"example.x.@gmail.com.us", false},
	{"example.x.x@gmail.com.us", true},
}

func TestValidateEmail(t *testing.T) {
	for _, ec := range ecsDeprecated {
		valid := ValidateEmail(ec.email)
		t.Log(ec.email)
		assert.Equal(t, valid, ec.valid)
	}
}

func TestValidateIdentityCard(t *testing.T) {
	testCases := []struct {
		Cid    string
		Result bool
	}{
		{
			Cid:    "123",
			Result: false,
		},
		{
			Cid:    "123",
			Result: false,
		},
		{
			Cid:    "512659199602034568",
			Result: false,
		},
		{
			Cid:    "110101199003073191",
			Result: true,
		},
	}

	for _, v := range testCases {
		result := ValidateIdentityCard(v.Cid)
		assert.Equal(t, v.Result, result)
	}
}
