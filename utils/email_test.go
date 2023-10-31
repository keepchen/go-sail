package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type emailChallenge struct {
	email string
	valid bool
}

var ecs = []emailChallenge{
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
	for _, ec := range ecs {
		valid := ValidateEmail(ec.email)
		t.Log(ec.email)
		assert.Equal(t, valid, ec.valid)
	}
}
