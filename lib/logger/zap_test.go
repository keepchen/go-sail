package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	t.Run("GetLogger-NonModule", func(t *testing.T) {
		assert.Panics(t, func() {
			GetLogger()
		})
	})

	t.Run("GetLogger", func(t *testing.T) {
		assert.Panics(t, func() {
			GetLogger("schedule")
		})
	})
}

func TestInit(t *testing.T) {
	t.Run("Init", func(t *testing.T) {
		conf := Conf{
			Filename: "../../examples/logs/logger_tester.log",
		}
		Init(conf, "go-sail")
	})

	t.Run("Init-WithSyncers", func(t *testing.T) {
		conf := Conf{
			Filename: "../../examples/logs/logger_tester.log",
		}
		redisWriter := &redisWriterStd{
			cli:     nil,
			listKey: "",
		}
		Init(conf, "go-sail", redisWriter)
	})

	t.Run("Init-Levels", func(t *testing.T) {
		conf := Conf{
			Filename: "../../examples/logs/logger_tester.log",
		}
		redisWriter := &redisWriterStd{
			cli:     nil,
			listKey: "",
		}

		levels := []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}

		for _, level := range levels {
			conf.Level = level
			Init(conf, "go-sail", redisWriter)
		}
	})

	t.Run("Init-Modules", func(t *testing.T) {
		conf := Conf{
			Filename: "../../examples/logs/logger_tester.log",
		}
		redisWriter := &redisWriterStd{
			cli:     nil,
			listKey: "",
		}

		modules := []string{"api", "schedule", "system"}
		levels := []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}

		for _, level := range levels {
			conf.Level = level
			conf.Modules = modules
			Init(conf, "go-sail", redisWriter)
		}
	})
}

func TestExporterProvider(t *testing.T) {
	t.Run("ExporterProvider", func(t *testing.T) {
		providers := []string{"redis", "redis-cluster", "nats", "kafka"}
		for _, provider := range providers {
			conf := Conf{
				Filename: "../../examples/logs/logger_tester.log",
			}
			conf.Exporter.Provider = provider
			t.Log(exporterProvider(conf))
		}
	})
}
