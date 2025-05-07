package etcd

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterService(t *testing.T) {
	t.Run("RegisterService", func(t *testing.T) {
		assert.Panics(t, func() {
			t.Log(RegisterService(context.Background(), "go-sail", "endpoint", 60))
		})
	})
}

func TestDiscoverService(t *testing.T) {
	t.Run("DiscoverService", func(t *testing.T) {
		assert.Panics(t, func() {
			t.Log(DiscoverService(context.Background(), "go-sail"))
		})
	})
}

func TestWatchService(t *testing.T) {
	t.Run("WatchService", func(t *testing.T) {
		assert.Panics(t, func() {
			fn := func(k, v []byte) {
				fmt.Println(k, v)
			}
			WatchService(context.Background(), "go-sail", fn)
		})
	})
}
