package nacos

import (
	"fmt"

	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/utils"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"
)

// RegisterService registers a service with the service registry
//
// groupName group name
//
// serviceName service name
//
// ip access IP address; if empty, obtain it automatically with utils.IP().GetLocal()
//
// port listen port
//
// metadata metadata information
func RegisterService(groupName, serviceName string, ip string, port uint64, metadata map[string]string) (bool, error) {
	var param vo.RegisterInstanceParam
	param.Ip = ip
	if len(param.Ip) == 0 {
		localIp, err := utils.IP().GetLocal()
		if err == nil {
			param.Ip = localIp
		}
	}
	param.Port = port
	param.GroupName = groupName
	param.ServiceName = serviceName
	param.Weight = 100
	param.Enable = true
	param.Healthy = true
	param.Ephemeral = true
	param.Metadata = metadata

	return GetNamingClient().RegisterInstance(param)
}

// UnregisterService deregisters a service from the service registry
//
// serviceName service name
//
// ip access IP address; if empty, obtain it automatically with utils.IP().GetLocal()
//
// port listen port
func UnregisterService(groupName, serviceName string, ip string, port uint64) (bool, error) {
	var param vo.DeregisterInstanceParam
	param.Ip = ip
	if len(param.Ip) == 0 {
		localIp, err := utils.IP().GetLocal()
		if err == nil {
			param.Ip = localIp
		}
	}
	param.Port = port
	param.GroupName = groupName
	param.ServiceName = serviceName
	param.Ephemeral = true

	return GetNamingClient().DeregisterInstance(param)
}

// GetHealthyInstanceUrl retrieves the URL of a healthy instance
//
// groupName group name
//
// serviceName service name
//
// loggerSvc logging component
func GetHealthyInstanceUrl(groupName, serviceName string, loggerSvc *zap.Logger) string {
	var param vo.SelectOneHealthInstanceParam
	param.GroupName = groupName
	param.ServiceName = serviceName
	instance, err := GetNamingClient().SelectOneHealthyInstance(param)
	if err != nil {
		loggerSvc.Error("GetHealthyInstanceUrl",
			zap.Any("value", logger.MarshalInterfaceValue([]string{param.GroupName, param.ServiceName})),
			zap.Errors("errors", []error{err}))
		return ""
	}

	return fmt.Sprintf("%s:%d", instance.Ip, instance.Port)
}

// SubscribeInstances subscribes to (watches) service instance changes
//
// groupName group name
//
// serviceName service name
//
// fn callback handler
//
// loggerSvc logging component
func SubscribeInstances(groupName, serviceName string, fn func(instances []model.Instance, err error), loggerSvc *zap.Logger) {
	var param = &vo.SubscribeParam{
		GroupName:         groupName,
		ServiceName:       serviceName,
		SubscribeCallback: fn,
	}

	err := GetNamingClient().Subscribe(param)
	if err != nil {
		loggerSvc.Warn("Subscribe Instances failed", zap.String("err", err.Error()))
	}
}

// UnsubscribeInstances unsubscribes from (stops watching) service instance changes
//
// groupName group name
//
// serviceName service name
//
// fn callback handler
//
// loggerSvc logging component
func UnsubscribeInstances(groupName, serviceName string, fn func(instances []model.Instance, err error), loggerSvc *zap.Logger) {
	var param = &vo.SubscribeParam{
		GroupName:         groupName,
		ServiceName:       serviceName,
		SubscribeCallback: fn,
	}

	err := GetNamingClient().Unsubscribe(param)
	if err != nil {
		loggerSvc.Warn("Subscribe Instances failed", zap.String("err", err.Error()))
	}
}
