package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestI18NString(t *testing.T) {
	assert.Equal(t, "en", LanguageEnglish.String())
}

func TestI18NToLowerCase(t *testing.T) {
	assert.Equal(t, "en", LanguageEnglish.ToLowerCase())
}

func TestI18NToUpperCase(t *testing.T) {
	assert.Equal(t, "EN", LanguageEnglish.ToUpperCase())
}

func TestI18NExist(t *testing.T) {
	assert.Equal(t, true, LanguageEnglish.Exist())
}
