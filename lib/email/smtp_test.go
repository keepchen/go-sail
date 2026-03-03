package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	conf := Conf{
		Host: awsSesHost,
		Port: awsSesPort,
	}
	t.Run("New", func(t *testing.T) {
		//t.Log(New(conf))
		assert.NotNil(t, New(conf))
	})
}

func TestSendMailWithTLS(t *testing.T) {
	conf := Conf{
		Host: awsSesHost,
		Port: awsSesPort,
	}
	client := New(conf)
	t.Run("New", func(t *testing.T) {
		//t.Log(client.SendMailWithTLS([]string{"tester@go-sail.dev"}, []byte(`test send`)))
		assert.Error(t, client.SendMailWithTLS([]string{"tester@go-sail.dev"}, []byte(`test send`)))
	})
}
