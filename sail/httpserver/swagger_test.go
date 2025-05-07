package httpserver

import (
	"testing"

	"github.com/keepchen/go-sail/v3/sail/config"
)

func TestRunSwaggerServerWhenEnable(t *testing.T) {
	t.Run("RunSwaggerServerWhenEnable", func(t *testing.T) {
		conf := config.HttpServerConf{}
		sConf := config.SwaggerConf{}
		r := InitGinEngine(conf)
		RunSwaggerServerWhenEnable(sConf, r)
	})

	t.Run("RunSwaggerServerWhenEnable-Path", func(t *testing.T) {
		conf := config.HttpServerConf{}
		sConf := config.SwaggerConf{
			Enable:      true,
			JsonPath:    "/swagger.json",
			RedocUIPath: "/api.html",
			FaviconPath: "/favicon.ico",
		}
		r := InitGinEngine(conf)
		RunSwaggerServerWhenEnable(sConf, r)
	})
}
