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
		ok2 := Put(key, v, time.Millisecond*100)
		assert.Equal(t, true, ok2)
	}
}

func TestGet(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		for _, v := range values {
			ok := Put(key, v, time.Millisecond*100)
			assert.Equal(t, true, ok)
			getV, ret := Get(key)
			assert.Equal(t, v, getV.(string))
			assert.Equal(t, 0, ret)
		}
	})
	t.Run("Get-Expired", func(t *testing.T) {
		ok := Put(key+"-expired", "go-sail", time.Second)
		assert.Equal(t, true, ok)
		time.Sleep(time.Second * 2)
		getV, ret := Get(key + "-expired")
		assert.Nil(t, getV)
		assert.Equal(t, -1, ret)
	})
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

	assert.Equal(t, false, Forget("not-exist-key"))
}

func TestExpire(t *testing.T) {
	for _, v := range values {
		ok := Put(key, v)
		assert.Equal(t, true, ok)

		ok = Expire(key, time.Minute)
		assert.Equal(t, true, ok)
	}
}
