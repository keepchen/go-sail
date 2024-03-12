// ----------- api doc definition -----------------

// @title          user - <go-sail>
// @version        1.0
// @description    This is an api document of go-sail.
// @termsOfService https://blog.keepchen.com

// @contact.name  keepchen
// @contact.url   https://blog.keepchen.com
// @contact.email keepchen2016@gmail.com

// @license.name MIT
// @license.url  https://github.com/keepchen/go-sail/LICENSE

// @Host     localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
// @description                Access Token protects our entity endpoints

// ----------- api doc definition -----------------

package user

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/keepchen/go-sail/v3/schedule"

	"github.com/keepchen/go-sail/v3/examples/pkg/app/user/http/routes"

	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/sail/config"

	"github.com/keepchen/go-sail/v3/constants"

	"github.com/keepchen/go-sail/v3/http/api"
	"github.com/keepchen/go-sail/v3/lib/nacos"
	"github.com/keepchen/go-sail/v3/sail"
	"github.com/keepchen/go-sail/v3/utils"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	userConfig "github.com/keepchen/go-sail/v3/examples/pkg/app/user/config"
)

// StartServer 启动服务
//
//go:generate swag init --dir ./ --output ./http/docs --parseDependency --parseInternal --generalInfo user.go
//go:generate redoc-cli bundle ./http/docs/*.yaml -o ./http/docs/docs.html
//go:generate node ../../../../plugins/redocly/redocly-copy.js ./http/docs/docs.html
func StartServer(wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		conf = &config.Config{
			LoggerConf: logger.Conf{
				Filename: "examples/logs/running.log",
			},
			HttpServer: config.HttpServerConf{
				Debug: true,
				Addr:  ":8000",
				Swagger: config.SwaggerConf{
					Enable:      true,
					RedocUIPath: "examples/pkg/app/user/http/docs/docs.html",
					JsonPath:    "examples/pkg/app/user/http/docs/swagger.json",
				},
				Prometheus: config.PrometheusConf{
					Enable:     true,
					Addr:       ":19100",
					AccessPath: "/metrics",
				},
				WebSocketRoutePath: "go-sail-ws",
			},
		}
		apiOption = &api.Option{
			EmptyDataStruct:  api.DefaultEmptyDataStructObject,
			ErrNoneCode:      constants.CodeType(200),
			ErrNoneCodeMsg:   "SUCCEED",
			ForceHttpCode200: true,
		}
		beforeFunc = func() {
			fmt.Println("call user function [before] to do something...")
		}
		afterFunc = func() {
			fmt.Println("call user function [after] to do something...")
			job0 := "print now datetime"
			cancel0 := schedule.NewJob(job0, func() {
				fmt.Println("now: ", utils.FormatDate(time.Now(), utils.YYYY_MM_DD_HH_MM_SS_EN))
			}).RunAt(schedule.EveryMinute)
			time.AfterFunc(time.Minute*3, cancel0)

			job1 := "print hello"
			cancel1 := schedule.NewJob(job1, func() {
				time.Sleep(time.Second * 10)
				fmt.Println(utils.FormatDate(time.Now(), utils.YYYY_MM_DD_HH_MM_SS_EN), "hello")
			}).EverySecond()
			time.AfterFunc(time.Second*33, cancel1)

			ticker := time.NewTicker(time.Second)
			times := 0
		LOOP:
			for range ticker.C {
				times++
				fmt.Printf("job: {%s} is running: %t | job: {%s} is running: %t\n",
					job0, schedule.JobIsRunning(job0), job1, schedule.JobIsRunning(job1))
				if times > 50 {
					break LOOP
				}
			}
		}
	)

	//挂载处理方法后启动
	sail.WakeupHttp("go-sail", conf).
		SetupApiOption(apiOption).
		EnableWebsocket(nil, nil).
		Hook(routes.RegisterRoutes, beforeFunc, afterFunc).
		Launch()
}

// RegisterServicesToNacos 将服务注册到注册中心
func RegisterServicesToNacos(wg *sync.WaitGroup) {
	defer wg.Done()

	nc := nacos.GetNamingClient()
	var param vo.RegisterInstanceParam
	localIp, err := utils.GetLocalIP()
	if err == nil {
		param.Ip = localIp
	}
	portSplit := strings.Split(userConfig.GetGlobalConfig().HttpServer.Addr, ":")
	if len(portSplit) > 0 {
		portInt, err := strconv.Atoi(portSplit[len(portSplit)-1])
		if err == nil {
			param.Port = uint64(portInt)
		}
	}
	param.ServiceName = "go-sail-user"
	param.Weight = 100
	param.Enable = true
	param.Healthy = true
	param.Ephemeral = true
	param.Metadata = map[string]string{"description": "go-sail-user服务"}
	_, err = nc.RegisterInstance(param)
	if err != nil {
		panic(err)
	}
}

// UnregisterServiceFromNacos 将服务从注册中心下线
func UnregisterServiceFromNacos() {
	nc := nacos.GetNamingClient()
	var param vo.DeregisterInstanceParam
	localIp, err := utils.GetLocalIP()
	if err == nil {
		param.Ip = localIp
	}
	portSplit := strings.Split(userConfig.GetGlobalConfig().HttpServer.Addr, ":")
	if len(portSplit) > 0 {
		portInt, err := strconv.Atoi(portSplit[len(portSplit)-1])
		if err == nil {
			param.Port = uint64(portInt)
		}
	}
	param.ServiceName = "go-sail-user"
	param.Ephemeral = true

	_, _ = nc.DeregisterInstance(param)
}
