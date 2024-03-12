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
	"github.com/keepchen/go-sail/v3/sail/config"
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
func GetLogger(module ...string) *zap.Logger {
	return logger.GetLogger(module...)
}

// Response http响应组件
func Response(c *gin.Context) api.Responder {
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

// 根据配置依次初始化组件
func componentsStartup(appName string, conf *config.Config) {
	//- logger
	logger.Init(conf.LoggerConf, appName)

	//- redis(standalone)
	if conf.RedisConf.Enable {
		redis.InitRedis(conf.RedisConf)
	}

	//- redis(cluster)
	if conf.RedisClusterConf.Enable {
		redis.InitRedisCluster(conf.RedisClusterConf)
	}

	//- database
	if conf.DBConf.Enable {
		db.Init(conf.DBConf)
	}

	//- jwt
	if conf.JwtConf.Enable {
		conf.JwtConf.Load()
	}

	//- nats
	if conf.NatsConf.Enable {
		nats.Init(conf.NatsConf)
	}

	//- kafka
	if conf.KafkaConf.Conf.Enable {
		kafka.Init(conf.KafkaConf.Conf, conf.KafkaConf.Topic, conf.KafkaConf.GroupID)
	}

	//- etcd
	if conf.EtcdConf.Enable {
		etcd.Init(conf.EtcdConf)
	}
}
