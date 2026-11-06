//@doc https://github.com/etcd-io/etcd/tree/main/client/v3

package etcd

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var client *clientv3.Client

// Init initializes a connection instance
func Init(conf Conf) {
	if conf.Timeout == 0 {
		conf.Timeout = 10000
	}

	cfg := clientv3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: time.Duration(conf.Timeout) * time.Millisecond,
	}

	if conf.Tls != nil {
		cfg.TLS = conf.Tls
	}
	if len(conf.Username) != 0 && len(conf.Password) != 0 {
		cfg.Username = conf.Username
		cfg.Password = conf.Password
	}

	cli, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}

	client = cli
}

// GetInstance retrieves the connection instance
func GetInstance() *clientv3.Client {
	return client
}

// New creates a connection instance
func New(conf Conf) (*clientv3.Client, error) {
	if conf.Timeout == 0 {
		conf.Timeout = 10000
	}

	cfg := clientv3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: time.Duration(conf.Timeout) * time.Millisecond,
	}

	if conf.Tls != nil {
		cfg.TLS = conf.Tls
	}
	if len(conf.Username) != 0 && len(conf.Password) != 0 {
		cfg.Username = conf.Username
		cfg.Password = conf.Password
	}

	return clientv3.New(cfg)
}
