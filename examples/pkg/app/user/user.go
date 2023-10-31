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

	"github.com/keepchen/go-sail/v3/lib/logger"

	"github.com/keepchen/go-sail/v3/examples/pkg/app/user/http/routes"

	"github.com/keepchen/go-sail/v3/constants"

	"github.com/keepchen/go-sail/v3/http/api"
	"github.com/keepchen/go-sail/v3/lib/nacos"
	"github.com/keepchen/go-sail/v3/sail"
	"github.com/keepchen/go-sail/v3/utils"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/keepchen/go-sail/v3/config"
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
			},
		}
		apiOption = &api.Option{
			EmptyDataStruct:  api.DefaultEmptyDataStructObject,
			ErrNoneCode:      constants.CodeType(200),
			ErrNoneCodeMsg:   "SUCCEED",
			ForceHttpCode200: true,
		}
		fn = func() {
			fmt.Println("call user function to do something...")
		}
	)

	sail.WakeupHttp("go-sail", conf, apiOption).Launch(routes.RegisterRoutes, fn)
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
