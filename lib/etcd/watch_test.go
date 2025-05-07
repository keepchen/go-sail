package etcd

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWatch(t *testing.T) {
	t.Run("Watch", func(t *testing.T) {
		assert.Panics(t, func() {
			fn := func(k, v []byte) {
				fmt.Println(k, v)
			}
			Watch(context.Background(), "go-sail", fn)
		})
	})
}

func TestWatchWithPrefix(t *testing.T) {
	t.Run("WatchWithPrefix", func(t *testing.T) {
		assert.Panics(t, func() {
			fn := func(k, v []byte) {
				fmt.Println(k, v)
			}
			WatchWithPrefix(context.Background(), "go-sail", fn)
		})
	})
}

func TestWatchWith(t *testing.T) {
	t.Run("WatchWith", func(t *testing.T) {
		assert.Panics(t, func() {
			fn := func(k, v []byte) {
				fmt.Println(k, v)
			}
			WatchWith(context.Background(), "go-sail", fn)
		})
	})
}
