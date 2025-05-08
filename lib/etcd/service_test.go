package etcd

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestRegisterService(t *testing.T) {
	t.Run("RegisterService", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		Init(conf)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		t.Log(RegisterService(ctx, "go-sail", "endpoint", 60))
	})
}

func TestDiscoverService(t *testing.T) {
	t.Run("DiscoverService", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		Init(conf)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		t.Log(DiscoverService(ctx, "go-sail"))
	})
}

func TestWatchService(t *testing.T) {
	t.Run("WatchService", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		Init(conf)
		fn := func(k, v []byte) {
			fmt.Println(k, v)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		WatchService(ctx, "go-sail", fn)
	})
}
