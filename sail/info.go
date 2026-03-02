package sail

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/lib/nacos"
	"github.com/keepchen/go-sail/v3/sail/config"
	"github.com/keepchen/go-sail/v3/utils"
)

// 启动终端打印概览信息
func printSummaryInfo(conf config.Config, ginEngine *gin.Engine) {
	var (
		protocol   = "http:"
		messages   bytes.Buffer
		localIp, _ = utils.IP().GetLocal()
		delimiter  = strings.Repeat("=", 120)
		docLink    = "Documentation: https://go-sail.dev"
		repoLink   = "Repository:    https://github.com/keepchen/go-sail"
		blank      = strings.Repeat(" ", len(delimiter)-len(repoLink)-len(constants.GoSailVersion)-11)
	)
	//- 基础信息开始位置
	messages.WriteString(delimiter)
	info := fmt.Sprintf("%s\n", constants.GoSailLogo)
	messages.WriteString(info)
	versionInfo := fmt.Sprintf("\n%s\n%s%s(version: %s)\n", docLink, repoLink, blank, constants.GoSailVersion)
	messages.WriteString(versionInfo)
	messages.WriteString(delimiter)
	//- 基础信息终止位置

	//- 服务及组件信息开始位置
	tw := table.NewWriter()
	style := table.StyleLight
	style.Title.Align = text.AlignCenter
	tw.SetStyle(style)
	tw.SetTitle("%s", "SUMMARY INFO")
	//tw.AppendHeader(table.Row{constants.GoSailLogo, "", "", "", ""}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
	tw.AppendRow(table.Row{"<<Routes>>", "", "", "", ""}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
	tw.AppendSeparator()
	//第一行标题
	tw.AppendRow(table.Row{
		"<Name>", "<Url / Port>", "<Description>", "", "",
	}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
	tw.AppendSeparator()
	listeningPorts := []string{conf.HttpServer.Addr}
	if conf.HttpServer.Prometheus.Enable {
		listeningPorts = append(listeningPorts, conf.HttpServer.Prometheus.Addr)
	}
	tw.AppendRow(table.Row{
		"Server", strings.Join(listeningPorts, " & "), "Listening port(s)", "", "",
	}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
	tw.AppendSeparator()
	if conf.HttpServer.Debug {
		ginEngine.GET("/go-sail", func(c *gin.Context) {
			c.String(http.StatusOK, fmt.Sprintf("%s\r\n\r\n/** This route only enabled in debug mode **/", constants.GoSailLogo))
		})
		welcomeRoutes := fmt.Sprintf("%s//%s%s%s", protocol, localIp, conf.HttpServer.Addr, "/go-sail")
		tw.AppendRow(table.Row{
			"Welcome", welcomeRoutes, "Display go-sail logo\n(Only enable in debug mode)", "", "",
		}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
		tw.AppendSeparator()
	}
	if conf.HttpServer.Swagger.Enable {
		swaggerRoutes := fmt.Sprintf("%s//%s%s%s", protocol, localIp, conf.HttpServer.Addr, "/swagger/index.html")
		tw.AppendRow(table.Row{
			"API Docs", swaggerRoutes, "Swagger UI", "", "",
		}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
		redoclyRoutes := fmt.Sprintf("%s//%s%s%s", protocol, localIp, conf.HttpServer.Addr, "/redoc/docs.html")
		tw.AppendRow(table.Row{
			"", redoclyRoutes, "Redocly UI", "", "",
		}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
		tw.AppendSeparator()
	}
	if conf.HttpServer.Prometheus.Enable {
		if len(conf.HttpServer.Prometheus.AccessPath) == 0 {
			conf.HttpServer.Prometheus.AccessPath = "/metrics"
		}
		prometheusRoutes := fmt.Sprintf("%s//%s%s%s", protocol, localIp, conf.HttpServer.Prometheus.Addr, conf.HttpServer.Prometheus.AccessPath)
		tw.AppendRow(table.Row{
			"Monitor", prometheusRoutes, "Prometheus", "", "",
		}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
	}
	if conf.HttpServer.Debug {
		pprofRoutes := fmt.Sprintf("%s//%s%s%s", protocol, localIp, conf.HttpServer.Addr, "/debug/pprof")
		tw.AppendRow(table.Row{
			"", pprofRoutes, "PProf", "", "",
		}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
	}
	tw.AppendSeparator()
	tw.AppendRow(table.Row{"<<Components>>", "", "", "", ""}, table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter})
	tw.AppendSeparator()
	//第一行标题
	tw.AppendRow(table.Row{
		"<Name>", "<Library>", "<Description>", "<Enable>", "<Available>",
	})
	tw.AppendSeparator()
	level := "debug"
	if len(conf.LoggerConf.Level) > 0 {
		level = conf.LoggerConf.Level
	}
	tw.AppendRow(table.Row{
		"Logger", "go.uber.org/zap", fmt.Sprintf("Level: %s", level), "Yes", "Yes",
	})
	if conf.DBConf.Enable {
		txt := "No"
		if GetDBR() != nil && GetDBW() != nil {
			txt = "Yes"
		}
		logLevel := "info"
		if len(conf.LoggerConf.Level) > 0 {
			logLevel = conf.LoggerConf.Level
		}
		tw.AppendRow(table.Row{
			"Database", "gorm.io/gorm", fmt.Sprintf("Driver: %s | Log level: %s", conf.DBConf.DriverName, logLevel), "Yes", txt,
		})
	}
	if conf.RedisConf.Enable {
		txt := "No"
		if GetRedisStandalone() != nil {
			txt = "Yes"
		}
		tw.AppendRow(table.Row{
			"Redis", "go-redis/redis", fmt.Sprintf("Mode: %s", "standalone"), "Yes", txt,
		})
	}
	if conf.RedisClusterConf.Enable {
		txt := "No"
		if GetRedisCluster() != nil {
			txt = "Yes"
		}
		tw.AppendRow(table.Row{
			"Redis", "go-redis/redis", fmt.Sprintf("Mode: %s", "cluster"), "Yes", txt,
		})
	}
	if GetEtcdInstance() != nil {
		mode := "standalone"
		if len(GetEtcdInstance().Endpoints()) > 1 {
			mode = "cluster"
		}
		tw.AppendRow(table.Row{
			"Etcd", "go.etcd.io/etcd", fmt.Sprintf("Mode: %s", mode), "Yes", "Yes",
		})
	}
	if conf.NatsConf.Enable {
		txt := "No"
		if GetNats() != nil {
			txt = "Yes"
		}
		mode := "standalone"
		if len(conf.NatsConf.Endpoints) > 1 {
			mode = "cluster"
		}
		tw.AppendRow(table.Row{
			"Nats", "nats-io/nats.go", fmt.Sprintf("Mode: %s", mode), "Yes", txt,
		})
	}
	if conf.KafkaConf.Conf.Enable {
		txt := "No"
		if GetKafkaInstance() != nil {
			txt = "Yes"
		}
		mode := "standalone"
		if len(conf.KafkaConf.Conf.Endpoints) > 0 {
			mode = "cluster"
		}
		tw.AppendRow(table.Row{
			"Kafka", "segmentio/kafka-go", fmt.Sprintf("Mode: %s", mode), "Yes", txt,
		})
	}
	if conf.ValKeyConf.Enable {
		txt := "No"
		if GetValKey() != nil {
			txt = "Yes"
		}
		mode := "standalone"
		if len(conf.ValKeyConf.Endpoints) > 0 {
			mode = "cluster"
		}
		tw.AppendRow(table.Row{
			"Kafka", "valkey-io/valkey-go", fmt.Sprintf("Mode: %s", mode), "Yes", txt,
		})
	}
	if conf.JwtConf != nil && conf.JwtConf.Enable {
		tw.AppendRow(table.Row{
			"JWT", "golang-jwt/jwt", fmt.Sprintf("Algorithm: %s", conf.JwtConf.Algorithm), "Yes", "Yes",
		})
	}
	if nacos.GetConfigClient() != nil {
		tw.AppendRow(table.Row{
			"Nacos", "nacos-group/nacos-sdk-go", "Client: config", "Yes", "Yes",
		})
	}
	if nacos.GetNamingClient() != nil {
		tw.AppendRow(table.Row{
			"Nacos", "nacos-group/nacos-sdk-go", "Client: naming", "Yes", "Yes",
		})
	}
	//- 服务及组件信息结束位置
	tw.AppendSeparator()
	tw.SetCaption(fmt.Sprintf("Launched at: %s\n\n", time.Now().Format(time.DateTime)))

	fmt.Println(messages.String())
	fmt.Println(tw.Render())
}
