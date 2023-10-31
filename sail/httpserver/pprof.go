package httpserver

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/sail/config"
)

// EnablePProfOnDebugMode 启动pprof检测
func EnablePProfOnDebugMode(conf config.HttpServerConf, r *gin.Engine) {
	if conf.Debug {
		//仅在调试模式下才开始pprof检测
		pprof.Register(r, "/debug/pprof")
	}
}
