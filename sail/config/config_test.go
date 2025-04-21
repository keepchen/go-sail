package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintConfigFields(t *testing.T) {
	t.Log(&Config{})
}

func TestConfigSet(t *testing.T) {
	cfg := &Config{}
	cfg.HttpServer.Addr = ":8000"
	Set(cfg)
	assert.Equal(t, cfg.HttpServer.Addr, config.HttpServer.Addr)
}

func TestConfigGet(t *testing.T) {
	cfg := &Config{}
	cfg.HttpServer.Addr = ":8000"
	Set(cfg)
	assert.Equal(t, cfg.HttpServer.Addr, Get().HttpServer.Addr)
}
