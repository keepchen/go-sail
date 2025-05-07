package jwt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	contents = []byte(`hello
world
hello
go-sail`)
	dst = "./file_test_will_be_deleted_after_testcase.txt"

	rawStrArr = []string{
		"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",
		"abc",
		"",
	}
)

func TestFileGetContents(t *testing.T) {
	_ = os.WriteFile(dst, contents, 0644)
	defer func() {
		_ = os.Remove(dst)
	}()
	t.Run("fileGetContents", func(t *testing.T) {
		byt, err := fileGetContents(dst)
		assert.NoError(t, err)
		assert.Equal(t, true, len(byt) > 0)
	})
}

func TestFileExists(t *testing.T) {
	_ = os.WriteFile(dst, contents, 0644)
	defer func() {
		_ = os.Remove(dst)
	}()
	t.Run("fileExists", func(t *testing.T) {
		assert.Equal(t, true, fileExists(dst))
	})
}

func TestFileExistsWithError(t *testing.T) {
	_ = os.WriteFile(dst, contents, 0644)
	defer func() {
		_ = os.Remove(dst)
	}()
	t.Run("fileExistsWithError", func(t *testing.T) {
		ok, err := fileExistsWithError(dst)
		assert.Equal(t, true, ok)
		assert.NoError(t, err)
	})

	t.Run("fileExistsWithError-ErrNoExist", func(t *testing.T) {
		ok, err := fileExistsWithError(dst + "-unknown")
		assert.Equal(t, false, ok)
		assert.NoError(t, err)
	})
}

func TestWordwrap(t *testing.T) {
	t.Run("wordwrap", func(t *testing.T) {
		for _, rawStr := range rawStrArr {
			t.Log(wordwrap(rawStr, 20, "\n"))
		}
	})
}
