package email

import (
	"testing"
	"time"
)

var (
	awsSesHost = "email-smtp.us-west-2.amazonaws.com"
	awsSesPort = 465
)

func TestNewPool(t *testing.T) {
	t.Run("NewPool", func(t *testing.T) {
		conf := Conf{
			Host: awsSesHost,
			Port: awsSesPort,
		}
		t.Log(NewPool(conf))
	})
}

func TestMount(t *testing.T) {
	t.Run("Mount", func(t *testing.T) {
		conf := Conf{
			Host: awsSesHost,
			Port: awsSesPort,
		}
		pool := NewPool(conf)
		t.Log(pool)
		envelope := &Envelope{}
		t.Log(envelope)
		pool.Mount(0, envelope)
	})
}

func TestEmit(t *testing.T) {
	t.Run("Emit", func(t *testing.T) {
		conf := Conf{
			Host: awsSesHost,
			Port: awsSesPort,
		}
		pool := NewPool(conf)
		t.Log(pool)
		envelope := &Envelope{}
		t.Log(envelope)
		pool.Mount(0, envelope)

		pool.Emit()
	})
}

func TestDone(t *testing.T) {
	t.Run("Done", func(t *testing.T) {
		conf := Conf{
			Host: awsSesHost,
			Port: awsSesPort,
		}
		pool := NewPool(conf)
		t.Log(pool)
		envelope := &Envelope{}
		t.Log(envelope)
		pool.Mount(0, envelope)

		pool.Emit()

		time.Sleep(time.Second)
		pool.Done()
	})
}
