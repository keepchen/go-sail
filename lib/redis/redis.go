package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	redisLib "github.com/go-redis/redis/v8"
)

var redisInstance *redisLib.Client

// InitRedis 初始化redis连接
func InitRedis(conf Conf) {
	rdb := mustInitRedis(conf)

	redisInstance = rdb
}

// GetInstance 获取redis连接实例
//
// 获取由InitRedis方法实例化后的连接
func GetInstance() *redisLib.Client {
	return redisInstance
}

func mustInitRedis(conf Conf) *redisLib.Client {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	return rdb
}

func initRedis(conf Conf) (*redisLib.Client, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := rdb.Ping(ctx).Err()

	return rdb, err
}

// New 实例化新的实例
func New(conf Conf) (*redisLib.Client, error) {
	return initRedis(conf)
}
