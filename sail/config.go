package sail

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/keepchen/go-sail/v3/lib/etcd"
	"github.com/keepchen/go-sail/v3/lib/nacos"
	"github.com/keepchen/go-sail/v3/utils"

	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// Config 配置操作方法
//
// watcher - 监听处理函数
//
// initClientAsComponent - 是否初始化客户端作为组件
//
// # Note
//
// 当ConfigProvider = Nacos 或 Etcd 且连接成功时，initClientAsComponent = true会全局初始化对应的客户端
func Config(watcher func(configName string, content []byte, viaWatch bool), initClientAsComponent ...bool) ConfigProvider {
	var initClient = false
	if len(initClientAsComponent) > 0 {
		initClient = initClientAsComponent[0]
	}
	return &configImpl{
		watcher:               watcher,
		initClientAsComponent: initClient,
	}
}

// configImpl 配置操作方法实现
type configImpl struct {
	watcher               func(configName string, content []byte, viaWatch bool)
	initClientAsComponent bool
	fileLastModifyTime    time.Time
}

// ConfigProvider 配置操作提供方案
type ConfigProvider interface {
	// ViaFile 通过文件操作配置
	ViaFile(filename string) Parser
	// ViaEtcd 通过etcd操作配置
	ViaEtcd(conf etcd.Conf, key string) Parser
	// ViaNacos 通过nacos操作配置
	//
	// clientCfg 参数支持覆盖默认配置，例如认证操作
	ViaNacos(endpoints, namespaceID, groupName, dataID string, clientCfg ...constant.ClientConfig) Parser
}

var _ ConfigProvider = &configImpl{}

// ViaFile 通过文件操作配置
func (ci *configImpl) ViaFile(filename string) Parser {
	content, err := utils.File().GetContents(filename)
	if err != nil {
		panic(err)
	}

	file, err := os.Stat(filename)
	if err == nil {
		ci.fileLastModifyTime = file.ModTime()
	}

	//当监听函数被设置时，需要执行监听
	if ci.watcher != nil {
		watcher := func() {
			ticker := time.NewTicker(time.Second * 3)
			defer ticker.Stop()
			errChan := make(chan error)
			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
				errChan <- fmt.Errorf("%v", <-c)
			}()
		LISTEN:
			for {
				select {
				case <-ticker.C:
					{
						f, e := os.Stat(filename)
						if e == nil && f.ModTime().After(ci.fileLastModifyTime) {
							ci.fileLastModifyTime = f.ModTime()
							fmt.Println("[Go-Sail] <Config> content has been modified (via File)")
							data, gErr := utils.File().GetContents(filename)
							if gErr == nil {
								ci.watcher(filename, data, true)
							}
						}
					}
				case <-errChan:
					break LISTEN
				}
			}
		}

		go watcher()
	}

	return &parserImpl{configName: filename, content: content}
}

// ViaEtcd 通过etcd操作配置
func (ci *configImpl) ViaEtcd(conf etcd.Conf, key string) Parser {
	client, err := etcd.New(conf)
	if err != nil {
		panic(err)
	}
	if client == nil {
		panic("etcd client is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	etcdResp, err := client.Get(ctx, key)
	if err != nil {
		panic(err)
	}
	if etcdResp == nil {
		panic(err)
	}
	if len(etcdResp.Kvs) == 0 {
		panic(err)
	}

	//当监听函数被设置时，需要执行监听
	if ci.watcher != nil {
		watcher := func() {
			watchChan := client.Watch(context.Background(), key)
			for watchResp := range watchChan {
				for _, value := range watchResp.Events {
					if string(value.Kv.Key) == key {
						fmt.Println("[Go-Sail] <Config> content has been modified (via Etcd)")
						ci.watcher(key, value.Kv.Value, true)
					}
				}
			}
		}

		go watcher()
	}

	var content []byte
	content = etcdResp.Kvs[0].Value

	//初始化客户端
	if client != nil && ci.initClientAsComponent {
		etcd.Init(conf)
	}

	return &parserImpl{configName: key, content: content}
}

// ViaNacos 通过nacos操作配置
func (ci *configImpl) ViaNacos(endpoints, namespaceID, groupName, dataID string, clientCfg ...constant.ClientConfig) Parser {
	client, err := nacos.NewConfigClient("Go-Sail", endpoints, namespaceID, clientCfg...)
	if err != nil {
		panic(err)
	}
	if client == nil {
		panic("client is nil")
	}
	var content []byte
	data, err := client.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  groupName,
	})
	content = []byte(data)
	if err != nil {
		panic(err)
	}

	//当监听函数被设置时，需要执行监听
	if ci.watcher != nil {
		watcher := func(namespace string, group string, dataId string, data string) {
			fmt.Println("[Go-Sail] <Config> content has been modified (via Nacos)")
			ci.watcher(dataId, []byte(data), true)
		}
		_ = client.ListenConfig(vo.ConfigParam{
			DataId:   dataID,
			Group:    groupName,
			OnChange: watcher,
		})
	}

	//初始化客户端
	if ci.initClientAsComponent {
		nacos.InitClient("Go-Sail", endpoints, namespaceID, clientCfg...)
	}

	return &parserImpl{configName: dataID, content: content}
}

// parserImpl 解析器实现
type parserImpl struct {
	configName string
	content    []byte
}

// Parser 解析器接口
type Parser interface {
	// Parse 解析处理函数
	Parse(fn func(configName string, content []byte, viaWatch bool))
}

var _ Parser = &parserImpl{}

// Parse 解析处理函数
func (pi *parserImpl) Parse(fn func(configName string, content []byte, viaWatch bool)) {
	if fn != nil {
		fn(pi.configName, pi.content, false)
	}
}
