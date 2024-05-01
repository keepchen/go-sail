package sail

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/keepchen/go-sail/v3/sail/config"

	"github.com/gin-gonic/gin"

	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/utils"
)

// 启动终端打印概览信息
func printSummaryInfo(conf config.HttpServerConf, ginEngine *gin.Engine) {
	var (
		protocol     = "http:"
		messages     bytes.Buffer
		localIp, _   = utils.GetLocalIP()
		delimiter    = []byte(strings.Repeat("=", 88))
		subDelimiter = []byte(strings.Repeat("-", 88))
		blank        = strings.Repeat(" ", len(delimiter)-len(constants.GoSailVersion)-11)
	)
	messages.Write(delimiter)
	info := fmt.Sprintf("%s\n", constants.GoSailLogo)
	messages.Write([]byte(info))
	versionInfo := fmt.Sprintf("%s(version: %s)\n", blank, constants.GoSailVersion)
	messages.Write([]byte(versionInfo))
	messages.Write(subDelimiter)
	messages.Write([]byte("\n"))
	sMgs := fmt.Sprintf("[Server] Listening at: {%s}\n", conf.Addr)
	messages.Write([]byte(sMgs))
	if conf.Debug {
		ginEngine.GET("/go-sail", func(c *gin.Context) {
			c.String(http.StatusOK, fmt.Sprintf("%s\r\n\r\n/** This route only enabled in debug mode **/", constants.GoSailLogo))
		})
		msg := fmt.Sprintf(">\t%s//%s%s%s\n", protocol, localIp, conf.Addr, "/go-sail")
		messages.Write([]byte(msg))
	}
	if conf.Swagger.Enable {
		msg := fmt.Sprintf("[Swagger] Enabled:\n>\t%s//%s%s%s     (Redocly UI)\n>\t%s//%s%s%s  (Swagger UI)\n",
			protocol, localIp, conf.Addr, "/redoc/docs.html",
			protocol, localIp, conf.Addr, "/swagger/index.html")
		messages.Write([]byte(msg))
	}
	if conf.Prometheus.Enable {
		msg := fmt.Sprintf("[Prometheus] Enabled:\n>\t%s//%s%s%s\n", protocol, localIp, conf.Prometheus.Addr, conf.Prometheus.AccessPath)
		messages.Write([]byte(msg))
	}
	if conf.Debug {
		msg := fmt.Sprintf("[Pprof] Enabled:\n>\t%s//%s%s%s\n", protocol, localIp, conf.Addr, "/debug/pprof")
		messages.Write([]byte(msg))
	}
	messages.Write(delimiter)

	fmt.Println(messages.String())
}
