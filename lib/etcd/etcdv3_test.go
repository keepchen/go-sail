package etcd

import (
	"crypto/tls"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

var conf = Conf{
	Enable:    true,
	Endpoints: []string{"127.0.0.1:2379"},
}

func TestNew(t *testing.T) {
	t.Run("New-NonValue", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		t.Log(New(conf))
	})

	t.Run("New", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		conf.Tls = &tls.Config{}
		t.Log(New(conf))
	})
}

func TestInit(t *testing.T) {
	t.Run("Init-NonValue", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s", conf.Endpoints[0]))
		if err != nil {
			return
		}
		_ = conn.Close()
		Init(conf)
	})

	t.Run("Init", func(t *testing.T) {
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
