package logger

import (
	"testing"

	"github.com/keepchen/go-sail/v3/lib/nats"
)

var (
	natsConf = nats.Conf{
		Enable: true,
		Endpoints: []string{
			"127.0.0.1:4222",
			"127.0.0.1:4223",
			"127.0.0.1:4224",
		},
		Username: "username",
		Password: "password",
	}
)

func TestNatsSync(t *testing.T) {
	t.Run("Sync", func(t *testing.T) {
		writer := &natsWriterStd{}
		t.Log(writer.Sync())
	})
}

func TestNatsWrite(t *testing.T) {
	t.Run("Write", func(t *testing.T) {
		rd, err := nats.New(natsConf)
		t.Log(err)
		if rd == nil {
			return
		}
		writer := &natsWriterStd{
			cli:        rd,
			subjectKey: "go-sail-tester-logger-subject",
		}
		t.Log(writer.Sync())
	})
}
