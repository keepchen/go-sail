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

	t.Log(ErrNone.String("unknown"))
	assert.Contains(t, ErrNone.String("unknown"), "not defined")

	var notRegisteredCode = CodeType(-1)
	t.Log(notRegisteredCode.String("unknown"))
	assert.Contains(t, notRegisteredCode.String("unknown"), "not defined")
}

func TestInt(t *testing.T) {
	for lang, msgMap := range initErrorCodeMsgMap {
		RegisterCodeTable(lang, msgMap)
	}
	assert.Equal(t, 0, ErrNone.Int())
}
