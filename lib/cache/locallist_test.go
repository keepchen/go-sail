package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitList(t *testing.T) {
	t.Run("InitList", func(t *testing.T) {
		InitList()
	})
}

func TestGetListInstance(t *testing.T) {
	t.Run("GetListInstance", func(t *testing.T) {
		InitList()
		assert.NotNil(t, GetListInstance())
	})
}

func TestNewList(t *testing.T) {
	t.Run("NewList", func(t *testing.T) {
		assert.NotNil(t, NewList("new-list"))
	})
}
