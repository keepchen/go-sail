package redis

import (
	"context"
	"crypto/tls"
	"fmt"

	redisLib "github.com/go-redis/redis/v8"
)

var redisInstance *redisLib.Client

// InitRedis 初始化redis连接
func InitRedis(conf Conf) {
	rdb := initRedis(conf)

	redisInstance = rdb
}

// GetInstance 获取redis连接实例
//
// 获取由InitRedis方法实例化后的连接
func GetInstance() *redisLib.Client {
	return redisInstance
}

func initRedis(conf Conf) *redisLib.Client {
	opts := &redisLib.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Username: conf.Username,
		Password: conf.Password,
		DB:       conf.Database,
	}
	if conf.SSLEnable {
		//https://redis.uptrace.dev/guide/go-redis.html#using-tls
		//
		//To enable TLS/SSL, you need to provide an empty tls.Config.
		//If you're using private certs, you need to specify them in the tls.Config
		opts.TLSConfig = &tls.Config{}
	}
	rdb := redisLib.NewClient(opts)

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return rdb
}

// New 实例化新的实例
func New(conf Conf) *redisLib.Client {
	return initRedis(conf)
}
