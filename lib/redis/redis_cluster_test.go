package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var cConf = ClusterConf{
	Enable:    true,
	Endpoints: []Endpoint{},
}

func TestInitRedisCluster(t *testing.T) {
	t.Run("InitRedisCluster-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			InitRedisCluster(cConf)
		})
	})

	t.Run("InitRedisCluster-SSLEnable", func(t *testing.T) {
		assert.Panics(t, func() {
			cConf.SSLEnable = true
			InitRedisCluster(cConf)
		})
	})
}

func TestNewRedisCluster(t *testing.T) {
	t.Run("NewRedisCluster", func(t *testing.T) {
		t.Log(NewCluster(cConf))

	})

	t.Run("NewRedisCluster-SSLEnable", func(t *testing.T) {
		cConf.SSLEnable = true
		t.Log(NewCluster(cConf))
	})
}
