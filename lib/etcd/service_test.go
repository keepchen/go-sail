package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterService(t *testing.T) {
	t.Run("RegisterService", func(t *testing.T) {
		assert.Panics(t, func() {
			Init(conf)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			t.Log(RegisterService(ctx, "go-sail", "endpoint", 60))
		})
	})
}

func TestDiscoverService(t *testing.T) {
	t.Run("DiscoverService", func(t *testing.T) {
		assert.Panics(t, func() {
			Init(conf)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			t.Log(DiscoverService(ctx, "go-sail"))
		})
	})
}

func TestWatchService(t *testing.T) {
	t.Run("WatchService", func(t *testing.T) {
		assert.Panics(t, func() {
			Init(conf)
			fn := func(k, v []byte) {
				fmt.Println(k, v)
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			WatchService(ctx, "go-sail", fn)
		})
	})
}
