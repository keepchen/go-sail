package httpserver

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/sail/config"
)

// EnablePProfOnDebugMode starts pprof profiling
func EnablePProfOnDebugMode(conf config.HttpServerConf, r *gin.Engine) {
	if conf.Debug {
		//Enable pprof profiling only in debug mode
		pprof.Register(r, "/debug/pprof")
	}
}
