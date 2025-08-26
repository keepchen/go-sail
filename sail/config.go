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
func Config(panicWhileErr bool, watcher func(content []byte, viaWatch bool)) ConfigProvider {
	return &configImpl{
		panicWhileErr: panicWhileErr,
		watcher:       watcher,
	}
}

// configImpl 配置操作方法实现
type configImpl struct {
	panicWhileErr      bool
	watcher            func(content []byte, viaWatch bool)
	fileLastModifyTime time.Time
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
	if err != nil && ci.panicWhileErr {
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
							fmt.Println("[Go-Sail] <Config> content has been modified (via File)")
							data, gErr := utils.File().GetContents(filename)
							if gErr == nil {
								ci.watcher(data, true)
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

	return &parserImpl{content: content}
}

// ViaEtcd 通过etcd操作配置
func (ci *configImpl) ViaEtcd(conf etcd.Conf, key string) Parser {
	client, err := etcd.New(conf)
	if err != nil && ci.panicWhileErr {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	etcdResp, err := client.Get(ctx, key)
	if err != nil && ci.panicWhileErr {
		panic(err)
	}
	if etcdResp == nil && ci.panicWhileErr {
		panic(err)
	}
	if len(etcdResp.Kvs) == 0 && ci.panicWhileErr {
		panic(err)
	}

	//当监听函数被设置时，需要执行监听
	if ci.watcher != nil && client != nil {
		watcher := func() {
			watchChan := client.Watch(context.Background(), key)
			for watchResp := range watchChan {
				for _, value := range watchResp.Events {
					if string(value.Kv.Key) == key {
						fmt.Println("[Go-Sail] <Config> content has been modified (via Etcd)")
						ci.watcher(value.Kv.Value, true)
					}
				}
			}
		}

		go watcher()
	}

	var content []byte
	if len(etcdResp.Kvs) > 0 {
		content = etcdResp.Kvs[0].Value
	}

	return &parserImpl{content: content}
}

// ViaNacos 通过nacos操作配置
func (ci *configImpl) ViaNacos(endpoints, namespaceID, groupName, dataID string, clientCfg ...constant.ClientConfig) Parser {
	client, err := nacos.NewConfigClient("Go-Sail", endpoints, namespaceID, clientCfg...)
	if err != nil && ci.panicWhileErr {
		panic(err)
	}
	if client == nil && ci.panicWhileErr {
		panic("client is nil")
	}
	var content []byte
	if client != nil {
		data, err := client.GetConfig(vo.ConfigParam{
			DataId: dataID,
			Group:  groupName,
		})
		content = []byte(data)
		if err != nil && ci.panicWhileErr {
			panic(err)
		}
	}

	//当监听函数被设置时，需要执行监听
	if ci.watcher != nil && client != nil {
		watcher := func(namespace string, group string, dataId string, data string) {
			fmt.Println("[Go-Sail] <Config> content has been modified (via Nacos)")
			ci.watcher([]byte(data), true)
		}
		_ = client.ListenConfig(vo.ConfigParam{
			DataId:   dataID,
			Group:    groupName,
			OnChange: watcher,
		})
	}

	return &parserImpl{content: content}
}

// parserImpl 解析器实现
type parserImpl struct {
	content []byte
}

// Parser 解析器接口
type Parser interface {
	// Parse 解析处理函数
	Parse(fn func(content []byte, viaWatch bool))
}

var _ Parser = &parserImpl{}

// Parse 解析处理函数
func (pi *parserImpl) Parse(fn func(content []byte, viaWatch bool)) {
	if fn != nil {
		fn(pi.content, false)
	}
}
