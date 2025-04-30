package valkey

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
)

var instance valkey.Client

// Init 初始化连接实例
func Init(conf Conf) {
	client := mustInitVK(conf)
	instance = client
}

func mustInitVK(conf Conf) valkey.Client {
	initAddress := make([]string, 0, len(conf.Endpoints))

	for _, endpoint := range conf.Endpoints {
		initAddress = append(initAddress, fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port))
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: initAddress,
		Username:    conf.Username,
		Password:    conf.Password,
		ShuffleInit: true,
	})

	if err != nil {
		panic(err)
	}

	//ping
	err = client.Do(context.Background(), client.B().Ping().Build()).Error()
	if err != nil {
		panic(err)
	}

	return client
}

// GetValKey 获取连接实例
func GetValKey() valkey.Client {
	return instance
}

// New 新建连接实例
func New(conf Conf) (valkey.Client, error) {
	return initVK(conf)
}

func initVK(conf Conf) (valkey.Client, error) {
	initAddress := make([]string, 0, len(conf.Endpoints))

	for _, endpoint := range conf.Endpoints {
		initAddress = append(initAddress, fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port))
	}

	return valkey.NewClient(valkey.ClientOption{
		InitAddress: initAddress,
		Username:    conf.Username,
		Password:    conf.Password,
		ShuffleInit: true,
	})
}
