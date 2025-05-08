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
		t.Log(instance, err)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t.Log("----Put----")
		t.Log(instance.Put(ctx, "go-sail-key", "go-sail"))
		t.Log("----Get----")
		t.Log(instance.Get(ctx, "go-sail-key"))
	})

	t.Run("New-TlsEnable", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		conf.Tls = &tls.Config{}
		instance, err := New(conf)
		t.Log(instance, err)
		t.Log(instance, err)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t.Log("----Put----")
		resultP, errP := instance.Put(ctx, "go-sail-key", "go-sail")
		t.Log(resultP, errP)
		assert.Error(t, errP)
		t.Log("----Get----")
		resultG, errG := instance.Get(ctx, "go-sail-key")
		t.Log()
		t.Log(resultG, errG)
		assert.Error(t, errG)
	})
}

func TestInit(t *testing.T) {
	t.Run("Init", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		Init(conf)
	})

	t.Run("Init-TlsEnable", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		conf.Tls = &tls.Config{}
		Init(conf)
	})
}
