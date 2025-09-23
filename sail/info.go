package sail

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/sail/config"
	"github.com/keepchen/go-sail/v3/utils"
)

// 启动终端打印概览信息
func printSummaryInfo(conf config.HttpServerConf, ginEngine *gin.Engine) {
	var (
		protocol     = "http:"
		messages     bytes.Buffer
		localIp, _   = utils.IP().GetLocal()
		delimiter    = strings.Repeat("=", 88)
		subDelimiter = strings.Repeat("-", 88)
		repoLink     = "Repository: https://github.com/keepchen/go-sail"
		blank        = strings.Repeat(" ", len(delimiter)-len(repoLink)-len(constants.GoSailVersion)-11)
	)
	messages.WriteString(delimiter)
	info := fmt.Sprintf("%s\n", constants.GoSailLogo)
	messages.WriteString(info)
	versionInfo := fmt.Sprintf("\n%s%s(version: %s)\n", repoLink, blank, constants.GoSailVersion)
	messages.WriteString(versionInfo)
	messages.WriteString(subDelimiter)
	messages.WriteString("\n")
	sMgs := fmt.Sprintf("[Server] Listening at: {%s}\n", conf.Addr)
	messages.WriteString(sMgs)
	if conf.Debug {
		ginEngine.GET("/go-sail", func(c *gin.Context) {
			c.String(http.StatusOK, fmt.Sprintf("%s\r\n\r\n/** This route only enabled in debug mode **/", constants.GoSailLogo))
		})
		msg := fmt.Sprintf(">\t%s//%s%s%s\n", protocol, localIp, conf.Addr, "/go-sail")
		messages.WriteString(msg)
	}
	if conf.Swagger.Enable {
		msg := fmt.Sprintf("[Swagger] Enabled:\n>\t%s//%s%s%s     (Redocly UI)\n>\t%s//%s%s%s  (Swagger UI)\n",
			protocol, localIp, conf.Addr, "/redoc/docs.html",
			protocol, localIp, conf.Addr, "/swagger/index.html")
		messages.WriteString(msg)
	}
	if conf.Prometheus.Enable {
		if len(conf.Prometheus.AccessPath) == 0 {
			conf.Prometheus.AccessPath = "/metrics"
		}
		msg := fmt.Sprintf("[Prometheus] Enabled:\n>\t%s//%s%s%s\n", protocol, localIp, conf.Prometheus.Addr, conf.Prometheus.AccessPath)
		messages.WriteString(msg)
	}
	if conf.Debug {
		msg := fmt.Sprintf("[Pprof] Enabled:\n>\t%s//%s%s%s\n", protocol, localIp, conf.Addr, "/debug/pprof")
		messages.WriteString(msg)
	}
	messages.WriteString(delimiter)

	fmt.Println(messages.String())
}
