package nacos

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	nacosV2Clients "github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	nacosV2Constant "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var (
	iConfigClient config_client.IConfigClient
	iNamingClient naming_client.INamingClient
)

// InitClient 初始化客户端（命名空间客户端和配置客户端）
//
// 注意addrStr支持多个地址，以英文逗号分隔，如：
//
// 192.168.224.2:8848,192.168.224.3:8848
//
// clientCfg参数将覆盖默认clientConfig配置
func InitClient(appName, addrStr, namespace string, clientCfg ...nacosV2Constant.ClientConfig) {
	if len(addrStr) == 0 || len(namespace) == 0 {
		panic(errors.New("[addrStr] or [namespace] is empty"))
	}

	var servers []nacosV2Constant.ServerConfig

	addrSlice := strings.Split(addrStr, ",")
	for _, adr := range addrSlice {
		adrSlice := strings.Split(adr, ":")
		if len(adrSlice) == 2 {
			port, err := strconv.Atoi(adrSlice[1])
			if err != nil {
				panic(err)
			}
			servers = append(servers, nacosV2Constant.ServerConfig{
				IpAddr: adrSlice[0],
				Port:   uint64(port),
			})
		}
	}

	if len(servers) == 0 {
		panic(errors.New("no nacos servers set"))
	}

	var clientConfig = nacosV2Constant.ClientConfig{
		NamespaceId:         namespace, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           10000,
		NotLoadCacheAtStart: true,
		LogDir:              fmt.Sprintf("logs/nacos/log_%s", appName),
		CacheDir:            fmt.Sprintf("logs/nacos/cache_%s", appName),
		LogLevel:            "warn",
	}

	//如果传递了客户端配置，则覆盖默认配置
	if len(clientCfg) > 0 {
		clientConfig = clientCfg[0]
	}

	cc, err := nacosV2Clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: servers,
	})
	if err != nil {
		panic(err)
	}

	iConfigClient = cc

	//create naming client
	nc, err := nacosV2Clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: servers,
	})
	if err != nil {
		panic(err)
	}

	iNamingClient = nc
}

// GetConfigClient 获取配置实例
func GetConfigClient() config_client.IConfigClient {
	return iConfigClient
}

// GetNamingClient 获取服务发现实例
func GetNamingClient() naming_client.INamingClient {
	return iNamingClient
}
