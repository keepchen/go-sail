package email

import "testing"

func TestNew(t *testing.T) {
	conf := Conf{
		Host: awsSesHost,
		Port: awsSesPort,
	}
	t.Run("New", func(t *testing.T) {
		t.Log(New(conf))
	})
}

func TestSendMailWithTLS(t *testing.T) {
	conf := Conf{
		Host: awsSesHost,
		Port: awsSesPort,
	}
	client := New(conf)
	t.Run("New", func(t *testing.T) {
		t.Log(client.SendMailWithTLS([]string{"tester@go-sail.dev"}, []byte(`test send`)))
	})
}
