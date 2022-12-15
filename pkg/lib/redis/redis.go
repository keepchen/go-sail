package redis

import (
	"context"
	"crypto/tls"
	"fmt"

	redisLib "github.com/go-redis/redis/v8"
)

var redisInstance *redisLib.Client

//InitRedis 初始化redis连接
func InitRedis(conf Conf) {
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

	redisInstance = rdb
}

//GetInstance 获取redis连接实例
func GetInstance() *redisLib.Client {
	return redisInstance
}
