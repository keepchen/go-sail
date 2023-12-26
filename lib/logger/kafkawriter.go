package logger

import (
	"context"

	kafkaLib "github.com/segmentio/kafka-go"
)

type kafkaWriterStd struct {
	writer *kafkaLib.Writer
	topic  string
}

func (w *kafkaWriterStd) Sync() error {
	//TODO implement me

	return nil
}

func (w *kafkaWriterStd) Write(p []byte) (int, error) {
	err := w.writer.WriteMessages(context.Background(), kafkaLib.Message{
		Topic: w.topic,
		Value: p,
	})

	return len(p), err
}
