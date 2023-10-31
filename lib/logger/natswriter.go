package logger

import (
	"github.com/nats-io/nats.go"
)

type natsWriterStd struct {
	cli        *nats.Conn
	subjectKey string
}

func (w *natsWriterStd) Sync() error {
	//TODO implement me

	return nil
}

func (w *natsWriterStd) Write(p []byte) (int, error) {
	err := w.cli.Publish(w.subjectKey, p)

	return len(p), err
}
