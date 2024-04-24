package sail

import (
	"fmt"

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
//
// 注意，使用前请确保db组件已初始化成功。
func GetDB() (read *gorm.DB, write *gorm.DB) {
	read, write = db.GetInstance().R, db.GetInstance().W

	return
}

// GetDBR 获取数据库读实例
//
// 注意，使用前请确保db组件已初始化成功。
func GetDBR() *gorm.DB {
	return db.GetInstance().R
}

// GetDBW 获取数据库写实例
//
// 注意，使用前请确保db组件已初始化成功。
func GetDBW() *gorm.DB {
	return db.GetInstance().W
}

// GetRedis 获取redis连接(standalone)
//
// 单实例模式
//
// 注意，使用前请确保redis(standalone)组件已初始化成功。
func GetRedis() *redisLib.Client {
	return redis.GetInstance()
}

// GetRedisCluster 获取redis连接(cluster)
//
// cluster集群模式
//
// 注意，使用前请确保redis(cluster)组件已初始化成功。
func GetRedisCluster() *redisLib.ClusterClient {
	return redis.GetClusterInstance()
}

// GetNats 获取nats连接实例
//
// 注意，使用前请确保nats组件已初始化成功。
func GetNats() *natsLib.Conn {
	return nats.GetInstance()
}

// GetLogger 获取日志实例
//
// 注意，使用前请确保logger组件已初始化成功。
func GetLogger(module ...string) *zap.Logger {
	return logger.GetLogger(module...)
}

// Response http响应组件
func Response(c *gin.Context) api.Responder {
	return api.New(c)
}

// GetKafkaInstance 获取kafka完整实例
//
// 注意，使用前请确保kafka组件已初始化成功。
func GetKafkaInstance() *kafka.Instance {
	return kafka.GetInstance()
}

// GetKafkaConnections 获取kafka连接
//
// 注意，使用前请确保kafka组件已初始化成功。
func GetKafkaConnections() []*kafkaLib.Conn {
	return kafka.GetConnections()
}

// GetKafkaWriter 获取kafka写实例
//
// 注意，使用前请确保kafka组件已初始化成功。
func GetKafkaWriter() *kafkaLib.Writer {
	return kafka.GetWriter()
}

// GetKafkaReader 获取kafka读实例
//
// 注意，使用前请确保kafka组件已初始化成功。
func GetKafkaReader() *kafkaLib.Reader {
	return kafka.GetReader()
}

// GetEtcdInstance 获取etcd连接实例
//
// 注意，使用前请确保etcd组件已初始化成功。
func GetEtcdInstance() *clientv3.Client {
	return etcd.GetInstance()
}

// 根据配置依次初始化组件
func componentsStartup(appName string, conf *config.Config) {
	//- logger
	logger.Init(conf.LoggerConf, appName)
	fmt.Println("[GO-SAIL] <Components> initialize [logger] successfully")

	//- redis(standalone)
	if conf.RedisConf.Enable {
		redis.InitRedis(conf.RedisConf)
		fmt.Println("[GO-SAIL] <Components> initialize [redis standalone] successfully")
	}

	//- redis(cluster)
	if conf.RedisClusterConf.Enable {
		redis.InitRedisCluster(conf.RedisClusterConf)
		fmt.Println("[GO-SAIL] <Components> initialize [redis cluster] successfully")
	}

	//- database
	if conf.DBConf.Enable {
		db.Init(conf.DBConf)
		fmt.Println("[GO-SAIL] <Components> initialize [db] successfully")
	}

	//- jwt
	if conf.JwtConf != nil && conf.JwtConf.Enable {
		conf.JwtConf.Load()
		fmt.Println("[GO-SAIL] <Components> initialize [jwt] successfully")
	}

	//- nats
	if conf.NatsConf.Enable {
		nats.Init(conf.NatsConf)
		fmt.Println("[GO-SAIL] <Components> initialize [nats] successfully")
	}

	//- kafka
	if conf.KafkaConf.Conf.Enable {
		kafka.Init(conf.KafkaConf.Conf, conf.KafkaConf.Topic, conf.KafkaConf.GroupID)
		fmt.Println("[GO-SAIL] <Components> initialize [kafka] successfully")
	}

	//- etcd
	if conf.EtcdConf.Enable {
		etcd.Init(conf.EtcdConf)
		fmt.Println("[GO-SAIL] <Components> initialize [etcd] successfully")
	}
}

// 根据配置依次停止组件服务
func componentsShutdown(conf *config.Config) {
	//- redis(standalone)
	if conf.RedisConf.Enable && redis.GetInstance() != nil {
		_ = redis.GetInstance().Close()
		fmt.Println("[GO-SAIL] <Components> shutdown [redis standalone] successfully")
	}

	//- redis(cluster)
	if conf.RedisClusterConf.Enable && redis.GetClusterInstance() != nil {
		_ = redis.GetClusterInstance().Close()
		fmt.Println("[GO-SAIL] <Components> shutdown [redis cluster] successfully")
	}

	//- database
	if conf.DBConf.Enable && db.GetInstance() != nil {
		if db.GetInstance().R != nil {
			if rawDB, err := db.GetInstance().R.DB(); err == nil {
				_ = rawDB.Close()
			}
		}
		if db.GetInstance().W != nil {
			if rawDB, err := db.GetInstance().W.DB(); err == nil {
				_ = rawDB.Close()
			}
		}
		fmt.Println("[GO-SAIL] <Components> shutdown [db] successfully")
	}

	//- nats
	if conf.NatsConf.Enable && nats.GetInstance() != nil {
		nats.GetInstance().Close()
		fmt.Println("[GO-SAIL] <Components> shutdown [nats] successfully")
	}

	//- kafka
	if conf.KafkaConf.Conf.Enable && kafka.GetInstance() != nil {
		if reader := kafka.GetInstance().Reader; reader != nil {
			_ = reader.Close()
		}
		if writer := kafka.GetInstance().Reader; writer != nil {
			_ = writer.Close()
		}
		fmt.Println("[GO-SAIL] <Components> shutdown [kafka] successfully")
	}

	//- etcd
	if conf.EtcdConf.Enable && etcd.GetInstance() != nil {
		_ = etcd.GetInstance().Close()
		fmt.Println("[GO-SAIL] <Components> shutdown [etcd] successfully")
	}
}
