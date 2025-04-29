package nacos

import (
	"fmt"

	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/utils"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"
)

// RegisterService 将服务注册到注册中心
//
// groupName 分组名称
//
// serviceName 服务名称
//
// ip 访问ip地址，如果为空，则使用 utils.IP().GetLocal() 自动获取
//
// port 监听的端口
//
// metadata 元数据信息
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

// UnregisterService 将服务从注册中心下线
//
// serviceName 服务名称
//
// ip 访问ip地址，如果为空，则使用 utils.IP().GetLocal() 自动获取
//
// port 监听的端口
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

// GetHealthyInstanceUrl 获取健康实例url地址
//
// groupName 分组名称
//
// serviceName 服务名称
//
// loggerSvc 日志组件
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

// SubscribeInstances 订阅(监听)服务实例变化
//
// groupName 分组名称
//
// serviceName 服务名称
//
// fn 回调处理函数
//
// loggerSvc 日志组件
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

// UnsubscribeInstances 取消订阅(监听)服务实例变化
//
// groupName 分组名称
//
// serviceName 服务名称
//
// fn 回调处理函数
//
// loggerSvc 日志组件
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
