package nats

import (
	"strings"

	natsLib "github.com/nats-io/nats.go"
)

var natsInstance *natsLib.Conn

// Init 初始化
func Init(conf Conf) {
	conn := mustInitNats(conf)

	natsInstance = conn
}

func mustInitNats(conf Conf) *natsLib.Conn {
	var opts natsLib.Option
	if len(conf.Username) != 0 && len(conf.Password) != 0 {
		opts = natsLib.UserInfo(conf.Username, conf.Password)
	}

	conn, err := natsLib.Connect(strings.Join(conf.Endpoints, ","), opts)

	if err != nil {
		panic(err)
	}

	return conn
}

func initNats(conf Conf) (*natsLib.Conn, error) {
	var opts natsLib.Option
	if len(conf.Username) != 0 && len(conf.Password) != 0 {
		opts = natsLib.UserInfo(conf.Username, conf.Password)
	}

	conn, err := natsLib.Connect(strings.Join(conf.Endpoints, ","), opts)

	return conn, err
}

// GetInstance 获取链接实例
func GetInstance() *natsLib.Conn {
	return natsInstance
}

// New 初始化新的nats实例
func New(conf Conf) (*natsLib.Conn, error) {
	return initNats(conf)
}
