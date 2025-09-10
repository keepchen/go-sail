package sail

import (
	"testing"

	"github.com/keepchen/go-sail/v3/constants"
	"github.com/stretchr/testify/assert"
)

func TestCode(t *testing.T) {
	t.Run("Code", func(t *testing.T) {
		t.Log(Code())
	})

	t.Run("Code-Register", func(t *testing.T) {
		var (
			code1 = 1
			msg   = "test message"
		)
		Code().Register("en", code1, msg)

		assert.Equal(t, msg, constants.CodeType(code1).String())
	})
}
