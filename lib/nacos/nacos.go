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

// InitClient initializes clients (naming and configuration clients)
//
// Note that endpoints supports multiple comma-separated addresses, for example:
//
// 192.168.224.2:8848,192.168.224.3:8848
//
// The clientCfg argument overrides the default clientConfig configuration
func InitClient(appName, endpoints, namespace string, clientCfg ...nacosV2Constant.ClientConfig) {
	//create config client
	cc, err := NewConfigClient(appName, endpoints, namespace, clientCfg...)
	if err != nil {
		panic(err)
	}

	iConfigClient = cc

	//create naming client
	nc, err := NewNamingClient(appName, endpoints, namespace, clientCfg...)
	if err != nil {
		panic(err)
	}

	iNamingClient = nc
}

// NewConfigClient creates a configuration client
//
// Note that endpoints supports multiple comma-separated addresses, for example:
//
// 192.168.224.2:8848,192.168.224.3:8848
//
// The clientCfg argument overrides the default clientConfig configuration
func NewConfigClient(appName, endpoints, namespace string, clientCfg ...nacosV2Constant.ClientConfig) (config_client.IConfigClient, error) {
	if len(endpoints) == 0 || len(namespace) == 0 {
		return nil, errors.New("[endpoints] or [namespace] is empty")
	}

	var servers []nacosV2Constant.ServerConfig

	addrSlice := strings.Split(endpoints, ",")
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
		return nil, errors.New("no nacos servers set")
	}

	var clientConfig = nacosV2Constant.ClientConfig{
		NamespaceId:         namespace, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           10000,
		NotLoadCacheAtStart: true,
		LogDir:              fmt.Sprintf("logs/nacos/log_%s", appName),
		CacheDir:            fmt.Sprintf("logs/nacos/cache_%s", appName),
		LogLevel:            "warn",
	}

	//Override defaults if a client configuration is provided
	if len(clientCfg) > 0 {
		clientConfig = clientCfg[0]
	}

	cc, err := nacosV2Clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: servers,
	})
	if err != nil {
		return nil, err
	}

	return cc, nil
}

// NewNamingClient creates a naming client
//
// Note that endpoints supports multiple comma-separated addresses, for example:
//
// 192.168.224.2:8848,192.168.224.3:8848
//
// The clientCfg argument overrides the default clientConfig configuration
func NewNamingClient(appName, endpoints, namespace string, clientCfg ...nacosV2Constant.ClientConfig) (naming_client.INamingClient, error) {
	if len(endpoints) == 0 || len(namespace) == 0 {
		return nil, errors.New("[endpoints] or [namespace] is empty")
	}

	var servers []nacosV2Constant.ServerConfig

	addrSlice := strings.Split(endpoints, ",")
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
		return nil, errors.New("no nacos servers set")
	}

	var clientConfig = nacosV2Constant.ClientConfig{
		NamespaceId:         namespace, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           10000,
		NotLoadCacheAtStart: true,
		LogDir:              fmt.Sprintf("logs/nacos/log_%s", appName),
		CacheDir:            fmt.Sprintf("logs/nacos/cache_%s", appName),
		LogLevel:            "warn",
	}

	//Override defaults if a client configuration is provided
	if len(clientCfg) > 0 {
		clientConfig = clientCfg[0]
	}

	//create naming client
	nc, err := nacosV2Clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: servers,
	})
	if err != nil {
		return nil, err
	}

	return nc, nil
}

// GetConfigClient retrieves the configuration instance
func GetConfigClient() config_client.IConfigClient {
	return iConfigClient
}

// GetNamingClient retrieves the service discovery instance
func GetNamingClient() naming_client.INamingClient {
	return iNamingClient
}
