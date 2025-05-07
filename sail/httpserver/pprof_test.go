package httpserver

import (
	"testing"

	"github.com/keepchen/go-sail/v3/sail/config"
)

func TestEnablePProfOnDebugMode(t *testing.T) {
	t.Run("EnablePProfOnDebugMode", func(t *testing.T) {
		conf := config.HttpServerConf{}
		r := InitGinEngine(conf)
		EnablePProfOnDebugMode(conf, r)
	})

	t.Run("EnablePProfOnDebugMode-Enable", func(t *testing.T) {
		conf := config.HttpServerConf{
			Debug: true,
		}
		r := InitGinEngine(conf)
		EnablePProfOnDebugMode(conf, r)
	})
}
