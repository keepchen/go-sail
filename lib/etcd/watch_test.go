package etcd

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	t.Run("Watch", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		fn := func(k, v []byte) {
			fmt.Println(k, v)
		}
		Init(conf)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		Watch(ctx, "go-sail", fn)
	})
}

func TestWatchWithPrefix(t *testing.T) {
	t.Run("WatchWithPrefix", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		fn := func(k, v []byte) {
			fmt.Println(k, v)
		}
		Init(conf)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		WatchWithPrefix(ctx, "go-sail", fn)
	})
}

func TestWatchWith(t *testing.T) {
	t.Run("WatchWith", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		fn := func(k, v []byte) {
			fmt.Println(k, v)
		}
		Init(conf)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		WatchWith(ctx, "go-sail", fn)
	})
}
