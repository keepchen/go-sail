package kafka

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConnections(t *testing.T) {
	t.Run("GetConnections", func(t *testing.T) {
		assert.Panics(t, func() {
			GetConnections()
		})
	})
}

func TestGetWriter(t *testing.T) {
	t.Run("GetWriter", func(t *testing.T) {
		assert.Panics(t, func() {
			GetWriter()
		})
	})
}

func TestGetInstance(t *testing.T) {
	t.Run("GetInstance", func(t *testing.T) {
		assert.Nil(t, GetInstance())
	})
}

func TestInit(t *testing.T) {
	t.Run("Init", func(t *testing.T) {
		assert.Panics(t, func() {
			conf := Conf{}
			conf.Tls = &tls.Config{}
			Init(conf, "go-sail", "tester")
		})
	})
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		conf := Conf{}
		conf.Tls = &tls.Config{}
		assert.Panics(t, func() {
			t.Log(New(conf, "go-sail", "tester"))
		})
	})
}

func TestInitConnections(t *testing.T) {
	t.Run("InitConnections", func(t *testing.T) {
		conf := Conf{
			Endpoints: []string{"localhost:9092", "localhost:9093"},
		}
		conf.Tls = &tls.Config{}
		assert.Panics(t, func() {
			InitConnections(conf)
		})
	})

	t.Run("InitConnections-NonValue", func(t *testing.T) {
		conf := Conf{
			Endpoints: []string{"localhost:9092", "localhost:9093"},
			Username:  "username",
			Password:  "password",
		}
		conf.Tls = &tls.Config{}
		assert.Panics(t, func() {
			InitConnections(conf)
		})
	})
}

func TestNewConnections(t *testing.T) {
	t.Run("NewConnections", func(t *testing.T) {
		conf := Conf{
			Endpoints: []string{"localhost:9092", "localhost:9093"},
		}
		conf.Tls = &tls.Config{}
		assert.Panics(t, func() {
			NewConnections(conf)
		})
	})

	t.Run("NewConnections-NonValue", func(t *testing.T) {
		conf := Conf{
			Endpoints: []string{"localhost:9092", "localhost:9093"},
			Username:  "username",
			Password:  "password",
		}
		conf.Tls = &tls.Config{}
		assert.Panics(t, func() {
			NewConnections(conf)
		})
	})
}

func TestInitWriter(t *testing.T) {
	t.Run("InitWriter", func(t *testing.T) {
		conf := Conf{
			Endpoints: []string{"localhost:9092", "localhost:9093"},
		}
		conf.Tls = &tls.Config{}
		InitWriter(conf, "go-sail")
	})

	t.Run("InitWriter-NonValue", func(t *testing.T) {
		conf := Conf{
			Endpoints: []string{"localhost:9092", "localhost:9093"},
			Username:  "username",
			Password:  "password",
		}
		conf.Tls = &tls.Config{}
		InitWriter(conf, "go-sail")
	})
}

func TestNewWriter(t *testing.T) {
	t.Run("NewWriter", func(t *testing.T) {
		conf := Conf{
			Endpoints: []string{"localhost:9092", "localhost:9093"},
		}
		conf.Tls = &tls.Config{}
		t.Log(NewWriter(conf, "go-sail"))
	})

	t.Run("NewWriter-NonValue", func(t *testing.T) {
		conf := Conf{
			Endpoints: []string{"localhost:9092", "localhost:9093"},
			Username:  "username",
			Password:  "password",
		}
		conf.Tls = &tls.Config{}
		t.Log(NewWriter(conf, "go-sail"))
	})
}

func TestInitReader(t *testing.T) {
	t.Run("InitReader", func(t *testing.T) {
		conf := Conf{
			Endpoints: []string{"localhost:9092", "localhost:9093"},
		}
		conf.Tls = &tls.Config{}
		InitReader(conf, "go-sail", "tester")
	})

	t.Run("InitReader-NonValue", func(t *testing.T) {
		conf := Conf{
			Endpoints: []string{"localhost:9092", "localhost:9093"},
			Username:  "username",
			Password:  "password",
		}
		conf.Tls = &tls.Config{}
		InitReader(conf, "go-sail", "tester")
	})
}

func TestNewReader(t *testing.T) {
	t.Run("NewReader", func(t *testing.T) {
		conf := Conf{
			Endpoints: []string{"localhost:9092", "localhost:9093"},
		}
		conf.Tls = &tls.Config{}
		t.Log(NewReader(conf, "go-sail", "tester"))
	})

	t.Run("NewReader-NonValue", func(t *testing.T) {
		conf := Conf{
			Endpoints: []string{"localhost:9092", "localhost:9093"},
			Username:  "username",
			Password:  "password",
		}
		conf.Tls = &tls.Config{}
		t.Log(NewReader(conf, "go-sail", "tester"))
	})
}

func TestGetMechanism(t *testing.T) {
	t.Run("getMechanism", func(t *testing.T) {
		conf := Conf{
			Endpoints: []string{"localhost:9092", "localhost:9093"},
			Username:  "username",
			Password:  "password",
		}
		conf.Tls = &tls.Config{}

		authTypes := []string{"sha256", "sha512", "plain", "unknown"}
		for _, at := range authTypes {
			conf.SASLAuthType = at
			t.Log(getMechanism(conf))
		}
	})
}
