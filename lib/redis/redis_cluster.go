package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	redisLib "github.com/go-redis/redis/v8"
)

var redisClusterInstance *redisLib.ClusterClient

// InitRedisCluster 初始化redis集群连接
func InitRedisCluster(conf ClusterConf) {
	rdb := mustInitRedisCluster(conf)

	redisClusterInstance = rdb
}

// GetClusterInstance 获取redis集群连接实例
func GetClusterInstance() *redisLib.ClusterClient {
	return redisClusterInstance
}

func mustInitRedisCluster(conf ClusterConf) *redisLib.ClusterClient {
	var (
		endpoints = make([]string, len(conf.Endpoints))
		username  string
		password  string
	)
	for i := 0; i < len(conf.Endpoints); i++ {
		endpoints[i] = fmt.Sprintf("%s:%d", conf.Endpoints[i].Host, conf.Endpoints[i].Port)
		if len(conf.Endpoints[i].Password) != 0 {
			password = conf.Endpoints[i].Password
		}
		if len(conf.Endpoints[i].Username) != 0 {
			username = conf.Endpoints[i].Username
		}
	}
	opts := &redisLib.ClusterOptions{
		Addrs:        endpoints,
		Username:     username,
		Password:     password,
		MaxRedirects: len(conf.Endpoints) - 1,
	}
	if opts.MaxRedirects < 3 {
		opts.MaxRedirects = 3
	}
	if conf.SSLEnable {
		//https://redis.uptrace.dev/guide/go-redis.html#using-tls
		//
		//To enable TLS/SSL, you need to provide an empty tls.Config.
		//If you're using private certs, you need to specify them in the tls.Config
		opts.TLSConfig = &tls.Config{}
	}
	rdb := redisLib.NewClusterClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := rdb.ForEachShard(ctx, func(ctx context.Context, shard *redisLib.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		panic(err)
	}

	return rdb
}

func initRedisCluster(conf ClusterConf) (*redisLib.ClusterClient, error) {
	var (
		endpoints = make([]string, len(conf.Endpoints))
		username  string
		password  string
	)
	for i := 0; i < len(conf.Endpoints); i++ {
		endpoints[i] = fmt.Sprintf("%s:%d", conf.Endpoints[i].Host, conf.Endpoints[i].Port)
		if len(conf.Endpoints[i].Password) != 0 {
			password = conf.Endpoints[i].Password
		}
		if len(conf.Endpoints[i].Username) != 0 {
			username = conf.Endpoints[i].Username
		}
	}
	opts := &redisLib.ClusterOptions{
		Addrs:        endpoints,
		Username:     username,
		Password:     password,
		MaxRedirects: len(conf.Endpoints) - 1,
	}
	if opts.MaxRedirects < 3 {
		opts.MaxRedirects = 3
	}
	if conf.SSLEnable {
		//https://redis.uptrace.dev/guide/go-redis.html#using-tls
		//
		//To enable TLS/SSL, you need to provide an empty tls.Config.
		//If you're using private certs, you need to specify them in the tls.Config
		opts.TLSConfig = &tls.Config{}
	}
	rdb := redisLib.NewClusterClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := rdb.ForEachShard(ctx, func(ctx context.Context, shard *redisLib.Client) error {
		return shard.Ping(ctx).Err()
	})

	return rdb, err
}

// NewCluster 实例化新的实例
func NewCluster(conf ClusterConf) (*redisLib.ClusterClient, error) {
	return initRedisCluster(conf)
}
