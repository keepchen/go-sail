package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	key    = "test-local-cache"
	values = []string{"a", "b", "c", "d"}
)

func TestPut(t *testing.T) {
	for _, v := range values {
		ok := Put(key, v)
		assert.Equal(t, true, ok)
	}
}

func TestGet(t *testing.T) {
	for _, v := range values {
		ok := Put(key, v)
		assert.Equal(t, true, ok)
		getV, ret := Get(key)
		assert.Equal(t, v, getV.(string))
		assert.Equal(t, 0, ret)
	}
}

func TestForget(t *testing.T) {
	for _, v := range values {
		ok := Put(key, v)
		assert.Equal(t, true, ok)

		ok = Forget(key)
		assert.Equal(t, true, ok)

		getV, ret := Get(key)
		assert.Equal(t, nil, getV)
		assert.Equal(t, -2, ret)
	}
}

func TestExpire(t *testing.T) {
	for _, v := range values {
		ok := Put(key, v)
		assert.Equal(t, true, ok)

		ok = Expire(key, time.Minute)
		assert.Equal(t, true, ok)
	}
}
