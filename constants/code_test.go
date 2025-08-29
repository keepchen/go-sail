package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterCode(t *testing.T) {
	const (
		c1000000 = CodeType(1000000)
		c1000001 = CodeType(1000001)
	)
	RegisterCode(LanguageEnglish, map[ICodeType]string{c1000000: "1000000", c1000001: "1000001"})
	assert.Equal(t, "1000000", c1000000.String())
	assert.Equal(t, "1000001", c1000001.String())
}

func TestRegisterCodeTable(t *testing.T) {
	const (
		c1000000 = CodeType(1000000)
		c1000001 = CodeType(1000001)
	)
	RegisterCodeTable(LanguageEnglish, map[ICodeType]string{c1000000: "1000000", c1000001: "1000001"})
	assert.Equal(t, "1000000", c1000000.String())
	assert.Equal(t, "1000001", c1000001.String())
}

func TestRegisterCodeSingle(t *testing.T) {
	const (
		c1000000 = CodeType(1000000)
		c1000001 = CodeType(1000001)
	)
	RegisterCodeSingle(LanguageEnglish, c1000000, "1000000")
	RegisterCodeSingle(LanguageEnglish, c1000001, "1000001")
	assert.Equal(t, "1000000", c1000000.String())
	assert.Equal(t, "1000001", c1000001.String())
}

func TestGetCodeMsg(t *testing.T) {
	t.Run("GetCodeMsg-OK", func(t *testing.T) {
		t.Log(GetCodeMsg(LanguageEnglish, ErrNone))
	})
	t.Run("GetCodeMsg-None-Language", func(t *testing.T) {
		t.Log(GetCodeMsg(LanguageCode("unknown"), ErrNone))
	})
	t.Run("GetCodeMsg-None-Code", func(t *testing.T) {
		t.Log(GetCodeMsg(LanguageEnglish, CodeType(1000000000000000)))
	})
}
