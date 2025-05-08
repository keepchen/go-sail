package logger

import (
	"github.com/keepchen/go-sail/v3/lib/kafka"
	"testing"
)

var (
	kafkaConf = kafka.Conf{
		Enable: true,
		Endpoints: []string{
			"127.0.0.1:9092",
			"127.0.0.1:9093",
		},
		Username: "username",
		Password: "password",
	}
)

func TestKafkaSync(t *testing.T) {
	t.Run("Sync", func(t *testing.T) {
		writer := &kafkaWriterStd{}
		t.Log(writer.Sync())
	})
}

func TestKafkaWrite(t *testing.T) {
	t.Run("Write", func(t *testing.T) {
		rd, err := kafka.NewWriter(kafkaConf, "go-sail-tester-logger-topic")
		t.Log(err)
		if rd == nil || err != nil {
			return
		}
		writer := &kafkaWriterStd{
			writer: rd,
			topic:  "go-sail-tester-logger-topic",
		}
		t.Log(writer.Sync())
	})
}
