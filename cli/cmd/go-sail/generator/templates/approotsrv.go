package templates

var AppRoorSrvTpl = `package {{ .ServiceName }}

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/keepchen/go-sail/pkg/lib/nacos"
	"github.com/keepchen/go-sail/pkg/utils"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"{{ .AppName }}/pkg/app/{{ .ServiceName }}/config"
	"github.com/keepchen/go-sail/pkg/lib/logger"
	"go.uber.org/zap"
)

// StartServer 启动服务
func StartServer(wg *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			logger.GetLogger().Error("---- Recovered ----", zap.Any("error", err))
		}
	}()

	//监听退出信号
	errChan := make(chan error)
	_, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%v", <-c)
		cancel()
	}()

	//启动xx服务
	log.Println("======================{{ .ServiceName }} service is ready======================")

	if nacos.GetNamingClient() != nil {
		logger.GetLogger().Sugar().Info("Register services to nacos")
		wg.Add(1)
		go RegisterServicesToNacos(wg)
	}

	//收到退出信号
	logger.GetLogger().Sugar().Warnf("Exit signal: %v", <-errChan)

	wg.Done()

	if nacos.GetNamingClient() != nil {
		UnregisterServiceFromNacos()
	}

	logger.GetLogger().Sugar().Warnf("Shutting down api server at %s", config.GetGlobalConfig().HttpServer.Addr)
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
	portSplit := strings.Split(config.GetGlobalConfig().HttpServer.Addr, ":")
	if len(portSplit) > 0 {
		portInt, err := strconv.Atoi(portSplit[len(portSplit)-1])
		if err == nil {
			param.Port = uint64(portInt)
		}
	}
	param.ServiceName = "{{ .AppName }}-{{ .ServiceName }}"
	param.Weight = 100
	param.Enable = true
	param.Healthy = true
	param.Ephemeral = true
	param.Metadata = map[string]string{"description": "{{ .AppName }}-{{ .ServiceName }}服务"}
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
	portSplit := strings.Split(config.GetGlobalConfig().HttpServer.Addr, ":")
	if len(portSplit) > 0 {
		portInt, err := strconv.Atoi(portSplit[len(portSplit)-1])
		if err == nil {
			param.Port = uint64(portInt)
		}
	}
	param.ServiceName = "{{ .AppName }}-{{ .ServiceName }}"
	param.Ephemeral = true

	_, _ = nc.DeregisterInstance(param)
}
`
