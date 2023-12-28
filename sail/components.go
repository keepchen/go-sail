package sail

import (
	"github.com/gin-gonic/gin"
	redisLib "github.com/go-redis/redis/v8"
	"github.com/keepchen/go-sail/v3/http/api"
	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/etcd"
	"github.com/keepchen/go-sail/v3/lib/kafka"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/lib/nats"
	"github.com/keepchen/go-sail/v3/lib/redis"
	natsLib "github.com/nats-io/nats.go"
	kafkaLib "github.com/segmentio/kafka-go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GetDB 获取数据库实例
//
// 该方法依次返回读实例、写实例
func GetDB() (read *gorm.DB, write *gorm.DB) {
	read, write = db.GetInstance().R, db.GetInstance().W

	return
}

// GetDBR 获取数据库读实例
func GetDBR() *gorm.DB {
	return db.GetInstance().R
}

// GetDBW 获取数据库写实例
func GetDBW() *gorm.DB {
	return db.GetInstance().W
}

// GetRedis 获取redis连接(standalone)
//
// 单实例模式
func GetRedis() *redisLib.Client {
	return redis.GetInstance()
}

// GetRedisCluster 获取redis连接(cluster)
//
// cluster集群模式
func GetRedisCluster() *redisLib.ClusterClient {
	return redis.GetClusterInstance()
}

// GetNats 获取nats连接实例
func GetNats() *natsLib.Conn {
	return nats.GetInstance()
}

// GetLogger 获取日志实例
func GetLogger(modules ...string) *zap.Logger {
	return logger.GetLogger(modules...)
}

// Response http响应组件
func Response(c *gin.Context) api.API {
	return api.New(c)
}

// GetKafkaInstance 获取kafka完整实例
func GetKafkaInstance() *kafka.Instance {
	return kafka.GetInstance()
}

// GetKafkaConnections 获取kafka连接
func GetKafkaConnections() []*kafkaLib.Conn {
	return kafka.GetConnections()
}

// GetKafkaWriter 获取kafka写实例
func GetKafkaWriter() *kafkaLib.Writer {
	return kafka.GetWriter()
}

// GetKafkaReader 获取kafka读实例
func GetKafkaReader() *kafkaLib.Reader {
	return kafka.GetReader()
}

// GetEtcdInstance 获取etcd连接实例
func GetEtcdInstance() *clientv3.Client {
	return etcd.GetInstance()
}
