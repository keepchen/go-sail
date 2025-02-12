package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	for lang, msgMap := range initErrorCodeMsgMap {
		RegisterCodeTable(lang, msgMap)
	}
	assert.Equal(t, "SUCCESS", ErrNone.String(LanguageEnglish.String()))
	assert.Equal(t, "成功", ErrNone.String(LanguageChinesePRC.String()))
}

func TestInt(t *testing.T) {
	assert.Equal(t, 0, ErrNone.Int())
}
