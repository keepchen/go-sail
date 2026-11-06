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

// Config configuration operations
//
// watcher - watch handler
//
// initClientAsComponent - whether to initialize the client as a component
//
// # Note
//
// When ConfigProvider = Nacos or Etcd and the connection succeeds, initClientAsComponent = true initializes the corresponding client globally
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

// configImpl configuration operations implementation
type configImpl struct {
	watcher               func(configName string, content []byte, viaWatch bool)
	initClientAsComponent bool
	fileLastModifyTime    time.Time
}

// ConfigProvider configuration provider
type ConfigProvider interface {
	// ViaFile manages configuration through files
	ViaFile(filename string) Parser
	// ViaEtcd manages configuration through etcd
	ViaEtcd(conf etcd.Conf, key string) Parser
	// ViaNacos manages configuration through Nacos
	//
	// The clientCfg argument supports overriding defaults, for example for authentication
	ViaNacos(endpoints, namespaceID, groupName, dataID string, clientCfg ...constant.ClientConfig) Parser
}

var _ ConfigProvider = &configImpl{}

// ViaFile manages configuration through files
func (ci *configImpl) ViaFile(filename string) Parser {
	content, err := utils.File().GetContents(filename)
	if err != nil {
		panic(err)
	}

	file, err := os.Stat(filename)
	if err == nil {
		ci.fileLastModifyTime = file.ModTime()
	}

	//Start watching when a watch handler is set
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

// ViaEtcd manages configuration through etcd
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

	//Start watching when a watch handler is set
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

	content := etcdResp.Kvs[0].Value

	//Initialize the client
	if client != nil && ci.initClientAsComponent {
		etcd.Init(conf)
	}

	return &parserImpl{configName: key, content: content}
}

// ViaNacos manages configuration through Nacos
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

	//Start watching when a watch handler is set
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

	//Initialize the client
	if ci.initClientAsComponent {
		nacos.InitClient("Go-Sail", endpoints, namespaceID, clientCfg...)
	}

	return &parserImpl{configName: dataID, content: content}
}

// parserImpl parser implementation
type parserImpl struct {
	configName string
	content    []byte
}

// Parser parser interface
type Parser interface {
	// Parse parsing handler
	Parse(fn func(configName string, content []byte, viaWatch bool))
}

var _ Parser = &parserImpl{}

// Parse parsing handler
func (pi *parserImpl) Parse(fn func(configName string, content []byte, viaWatch bool)) {
	if fn != nil {
		fn(pi.configName, pi.content, false)
	}
}
