package redis

import (
	"context"
	"crypto/tls"
	"fmt"

	redisLib "github.com/go-redis/redis/v8"
)

// OptionsClusterFields 配置项字段
type OptionsClusterFields struct {
	Addrs    []string //连接地址
	Password string   //密码
}

var redisClusterInstance *redisLib.ClusterClient

// InitRedisCluster 初始化redis集群连接
func InitRedisCluster(conf ClusterConf) {
	rdb := initRedisCluster(conf)

	redisClusterInstance = rdb
}

// GetClusterInstance 获取redis集群连接实例
func GetClusterInstance() *redisLib.ClusterClient {
	return redisClusterInstance
}

func initRedisCluster(conf ClusterConf) *redisLib.ClusterClient {
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

	err := rdb.ForEachShard(context.Background(), func(ctx context.Context, shard *redisLib.Client) error {
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		panic(err)
	}

	return rdb
}

// NewCluster 实例化新的实例
func NewCluster(conf ClusterConf) *redisLib.ClusterClient {
	return initRedisCluster(conf)
}
