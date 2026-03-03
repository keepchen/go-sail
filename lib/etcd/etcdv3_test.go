package etcd

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var conf = Conf{
	Enable:    true,
	Endpoints: []string{"127.0.0.1:2379"},
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		instance, err := New(conf)
		//t.Log(instance, err)
		assert.NoError(t, err)
		assert.NotNil(t, instance)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		pr, pe := instance.Put(ctx, "go-sail-key", "go-sail")
		assert.NoError(t, pe)
		assert.NotNil(t, pr)
		gr, ge := instance.Get(ctx, "go-sail-key")
		assert.NoError(t, ge)
		assert.NotNil(t, gr)
	})

	t.Run("New-Auth", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		conf.Username = "username"
		conf.Password = "password"
		instance, err := New(conf)
		//t.Log(instance, err)
		assert.NoError(t, err)
		assert.NotNil(t, instance)
	})

	t.Run("New-TlsEnable", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		assert.Panics(t, func() {
			conf.Tls = &tls.Config{}
			instance, err := New(conf)
			assert.Error(t, err)
			assert.Nil(t, instance)
			//t.Log(instance, err)
			//t.Log(instance, err)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			//t.Log("----Put----")
			resultP, errP := instance.Put(ctx, "go-sail-key", "go-sail")
			//t.Log(resultP, errP)
			assert.Error(t, errP)
			assert.Nil(t, resultP)
			//t.Log("----Get----")
			resultG, errG := instance.Get(ctx, "go-sail-key")
			//t.Log()
			//t.Log(resultG, errG)
			assert.Error(t, errG)
			assert.Nil(t, resultG)
		})
	})
}

func TestInit(t *testing.T) {
	t.Run("Init", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		conf.Tls = nil
		Init(conf)
	})

	t.Run("Init-Auth", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		conf.Tls = nil
		conf.Username = "username"
		conf.Password = "password"
		Init(conf)
	})

	t.Run("Init-TlsEnable", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		assert.Panics(t, func() {
			conf.Tls = &tls.Config{}
			Init(conf)
		})
	})
}
