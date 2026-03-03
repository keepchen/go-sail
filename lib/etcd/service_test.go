package etcd

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterService(t *testing.T) {
	t.Run("RegisterService", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		conf.Tls = nil
		Init(conf)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		rr, re := RegisterService(ctx, "go-sail", "testcase", "endpoint-local-tester", 60)
		assert.NoError(t, re)
		assert.NotNil(t, rr)

		//clear
		_ = GetInstance().Close()
		client = nil
	})

	t.Run("RegisterService-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			rr, re := RegisterService(ctx, "go-sail", "testcase", "endpoint-local-tester", 60)
			assert.Error(t, re)
			assert.Nil(t, rr)
		})
	})
}

func TestDiscoverService(t *testing.T) {
	t.Run("DiscoverService", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		conf.Tls = nil
		Init(conf)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		dr, de := DiscoverService(ctx, "go-sail", "testcase")
		if len(dr) == 0 {
			assert.Error(t, de)
		} else {
			assert.NoError(t, de)
		}

		//clear
		_ = GetInstance().Close()
		client = nil
	})

	t.Run("DiscoverService-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			dr, de := DiscoverService(ctx, "go-sail", "testcase")
			assert.Error(t, de)
			assert.Nil(t, dr)
		})
	})
}

func TestGetAllServices(t *testing.T) {
	t.Run("GetAllServices", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		conf.Tls = nil
		Init(conf)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		gr, ge := GetAllServices(ctx, "go-sail", "testcase")
		assert.NoError(t, ge)
		assert.NotNil(t, gr)

		//clear
		_ = GetInstance().Close()
		client = nil
	})

	t.Run("GetAllServices-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			gr, ge := GetAllServices(ctx, "go-sail", "testcase")
			assert.Error(t, ge)
			assert.Nil(t, gr)
		})
	})
}

func TestWatchService(t *testing.T) {
	t.Run("WatchService", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		conf.Tls = nil
		Init(conf)
		fn := func(k, v []byte) {
			fmt.Println(string(k), string(v))
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		go WatchService(ctx, "go-sail", "testcase", fn)

		time.Sleep(time.Second)
		rr, re := RegisterService(ctx, "go-sail", "testcase",
			fmt.Sprintf("endpoint-local-tester-%s", time.Now().String()), 60)
		assert.NoError(t, re)
		assert.NotNil(t, rr)

		time.Sleep(5 * time.Second)

		//clear
		_ = GetInstance().Close()
		client = nil
	})

	t.Run("WatchService-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			fn := func(k, v []byte) {
				fmt.Println(string(k), string(v))
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			WatchService(ctx, "go-sail", "testcase", fn)

			go func() {
				time.Sleep(time.Second)
				rr, re := RegisterService(ctx, "go-sail", "testcase",
					fmt.Sprintf("endpoint-local-tester-%s", time.Now().String()), 60)
				assert.Error(t, re)
				assert.Nil(t, rr)
			}()

			time.Sleep(5 * time.Second)
		})
	})
}

func TestGenerateInstanceID(t *testing.T) {
	t.Run("GenerateInstanceID", func(t *testing.T) {
		var endpoints = []string{
			"127.0.0.1:5000",
			"127.0.0.1:6000",
			"127.0.0.1:7000",
			"127.0.0.1:8000",
			"127.0.0.1:9000",
		}
		for _, endpoint := range endpoints {
			t.Log(generateInstanceID(endpoint))
		}
	})
}
