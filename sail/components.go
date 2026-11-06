package sail

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	redisLib "github.com/go-redis/redis/v8"
	"github.com/keepchen/go-sail/v3/http/api"
	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/etcd"
	"github.com/keepchen/go-sail/v3/lib/kafka"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/lib/nacos"
	"github.com/keepchen/go-sail/v3/lib/nats"
	"github.com/keepchen/go-sail/v3/lib/redis"
	"github.com/keepchen/go-sail/v3/lib/valkey"
	"github.com/keepchen/go-sail/v3/sail/config"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	natsLib "github.com/nats-io/nats.go"
	kafkaLib "github.com/segmentio/kafka-go"
	valkeyLib "github.com/valkey-io/valkey-go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GetDB retrieves database instances
//
// This method returns the read instance followed by the write instance
//
// Note: ensure the database component has been initialized before use.
func GetDB() (read *gorm.DB, write *gorm.DB) {
	if db.GetInstance() == nil {
		read, write = nil, nil
	} else {
		read, write = db.GetInstance().R, db.GetInstance().W
	}

	return
}

// NewDB creates new database instances
func NewDB(conf db.Conf) (read *gorm.DB, rErr error, write *gorm.DB, wErr error) {
	return db.New(conf)
}

// GetDBR retrieves the database read instance
//
// Note: ensure the database component has been initialized before use.
func GetDBR() *gorm.DB {
	if db.GetInstance() == nil {
		return nil
	}

	return db.GetInstance().R
}

// GetDBW retrieves the database write instance
//
// Note: ensure the database component has been initialized before use.
func GetDBW() *gorm.DB {
	if db.GetInstance() == nil {
		return nil
	}

	return db.GetInstance().W
}

// GetRedis retrieves a universal Redis connection
//
// # Shorthand for GetRedisUniversal
//
// Automatically detect currently initialized Redis connections
//
// Prefer the standalone instance if both standalone and cluster instances exist
//
// Use GetRedisStandalone or GetRedisCluster to retrieve a specific connection type
//
// ----
//
// Tip: when the connection type is known, use a type assertion to call advanced methods, for example:
//
// # - Conn (on standalone client)
//
// GetRedis().(*redis.Client).Conn()
//
// # - ForEachShard (on cluster client)
//
// GetRedis().(*redis.ClusterClient).ForEachShard()
func GetRedis() redisLib.UniversalClient {
	return GetRedisUniversal()
}

// GetRedisUniversal retrieves a universal Redis connection
//
// Automatically detect currently initialized Redis connections
//
// Prefer the standalone instance if both standalone and cluster instances exist
//
// Use GetRedisStandalone or GetRedisCluster to retrieve a specific connection type
//
// ----
//
// Tip: when the connection type is known, use a type assertion to call advanced methods, for example:
//
// # - Conn (on standalone client)
//
// GetRedisUniversal().(*redis.Client).Conn()
//
// # - ForEachShard (on cluster client)
//
// GetRedisUniversal().(*redis.ClusterClient).ForEachShard()
func GetRedisUniversal() redisLib.UniversalClient {
	switch true {
	default:
		return nil
	case redis.GetInstance() != nil:
		return redis.GetInstance()
	case redis.GetClusterInstance() != nil:
		return redis.GetClusterInstance()
	}
}

// GetRedisStandalone retrieves a Redis connection (standalone)
//
// Standalone mode
//
// Note: ensure the Redis (standalone) component has been initialized before use.
func GetRedisStandalone() *redisLib.Client {
	return redis.GetInstance()
}

// GetRedisCluster retrieves a Redis connection (cluster)
//
// Cluster mode
//
// Note: ensure the Redis (cluster) component has been initialized before use.
func GetRedisCluster() *redisLib.ClusterClient {
	return redis.GetClusterInstance()
}

// NewRedis creates a new Redis (standalone) connection instance
func NewRedis(conf redis.Conf) (*redisLib.Client, error) {
	return redis.New(conf)
}

// NewRedisCluster  creates a new Redis (cluster) connection instance
func NewRedisCluster(conf redis.ClusterConf) (*redisLib.ClusterClient, error) {
	return redis.NewCluster(conf)
}

// GetNats retrieves the NATS connection instance
//
// Note: ensure the NATS component has been initialized before use.
func GetNats() *natsLib.Conn {
	return nats.GetInstance()
}

// NewNats creates a new NATS connection instance
func NewNats(conf nats.Conf) (*natsLib.Conn, error) {
	return nats.New(conf)
}

// GetLogger retrieves the logger instance
//
// Note: ensure the logging component has been initialized before use.
func GetLogger(module ...string) *zap.Logger {
	return logger.GetLogger(module...)
}

// MarshalInterfaceValue serializes an interface to a string
//
// Mainly used for logging
func MarshalInterfaceValue(obj any) string {
	return logger.MarshalInterfaceValue(obj)
}

// Response HTTP response component
func Response(c *gin.Context) api.Responder {
	return api.Response(c)
}

// GetKafkaInstance retrieves the complete Kafka instance
//
// Note: ensure the Kafka component has been initialized before use.
func GetKafkaInstance() *kafka.Instance {
	return kafka.GetInstance()
}

// GetKafkaConnections retrieves Kafka connections
//
// Note: ensure the Kafka component has been initialized before use.
func GetKafkaConnections() []*kafkaLib.Conn {
	return kafka.GetConnections()
}

// GetKafkaWriter retrieves the Kafka writer instance
//
// Note: ensure the Kafka component has been initialized before use.
func GetKafkaWriter() *kafkaLib.Writer {
	return kafka.GetWriter()
}

// GetKafkaReader retrieves the Kafka reader instance
//
// Note: ensure the Kafka component has been initialized before use.
func GetKafkaReader() *kafkaLib.Reader {
	return kafka.GetReader()
}

// NewKafkaConnections creates new Kafka connection instances
func NewKafkaConnections(conf kafka.Conf) []*kafkaLib.Conn {
	return kafka.NewConnections(conf)
}

// GetEtcdInstance retrieves the etcd connection instance
//
// Note: ensure the etcd component has been initialized before use.
func GetEtcdInstance() *clientv3.Client {
	return etcd.GetInstance()
}

// NewEtcd creates a new etcd connection instance
func NewEtcd(conf etcd.Conf) (*clientv3.Client, error) {
	return etcd.New(conf)
}

// GetValKey retrieves the Valkey connection instance
//
// Note: ensure the Valkey component has been initialized before use.
func GetValKey() valkeyLib.Client {
	return valkey.GetValKey()
}

// NewValKey creates a new Valkey connection instance
func NewValKey(conf valkey.Conf) (valkeyLib.Client, error) {
	return valkey.New(conf)
}

// GetNacosConfigClient retrieves the Nacos configuration connection instance
//
// Note: ensure the Nacos component has been initialized before use.
func GetNacosConfigClient() config_client.IConfigClient {
	return nacos.GetConfigClient()
}

// GetNacosNamingClient retrieves the Nacos naming connection instance
//
// Note: ensure the Nacos component has been initialized before use.
func GetNacosNamingClient() naming_client.INamingClient {
	return nacos.GetNamingClient()
}

// NewNacosConfigClient creates a new Nacos configuration connection instance
func NewNacosConfigClient(appName string, endpoints string, namespace string,
	clientCfg ...constant.ClientConfig) (config_client.IConfigClient, error) {
	return nacos.NewConfigClient(appName, endpoints, namespace, clientCfg...)
}

// NewNacosNamingClient creates a new Nacos naming connection instance
func NewNacosNamingClient(appName string, endpoints string, namespace string,
	clientCfg ...constant.ClientConfig) (naming_client.INamingClient, error) {
	return nacos.NewNamingClient(appName, endpoints, namespace, clientCfg...)
}

var (
	monitoringCtx, monitoringCancel = context.WithCancel(context.Background())
)

// Initialize components in order according to configuration
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
	if conf.EtcdConf.Enable || etcd.GetInstance() != nil {
		// Configuration may have already initialized this component; first check for an available global variable and initialize only if absent
		if etcd.GetInstance() == nil {
			etcd.Init(conf.EtcdConf)
		}
		fmt.Println("[GO-SAIL] <Components> initialize [etcd] successfully")
	}

	//- valkey
	if conf.ValKeyConf.Enable {
		valkey.Init(conf.ValKeyConf)
		fmt.Println("[GO-SAIL] <Components> initialize [valkey] successfully")
	}

	//- nacos
	// Configuration may have already initialized this component; first check for an available global variable and print its information if present
	if nacos.GetConfigClient() != nil {
		fmt.Println("[GO-SAIL] <Components> initialize [nacos - configClient] successfully")
	}
	if nacos.GetNamingClient() != nil {
		fmt.Println("[GO-SAIL] <Components> initialize [nacos - namingClient] successfully")
	}

	//-monitor/distribution
	if conf.MonitorDistributorConf.Enabled {
		if GetRedis() != nil {
			Monitoring(monitoringCtx, GetLogger()).Distribution(&conf.MonitorDistributorConf, GetRedis())
			fmt.Println("[GO-SAIL] <Components> initialize [monitor-distribution] successfully")
		} else {
			fmt.Println("[GO-SAIL] <Components> initialize [monitor-distribution] canceled, because the redis is disabled")
		}
	}
}

// Stop component services in order according to configuration
func componentsShutdown(conf *config.Config) {
	//- monitor/distribution
	if conf.MonitorDistributorConf.Enabled {
		monitoringCancel()
	}
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
	if etcd.GetInstance() != nil {
		_ = etcd.GetInstance().Close()
		fmt.Println("[GO-SAIL] <Components> shutdown [etcd] successfully")
	}

	//- valkey
	if conf.ValKeyConf.Enable && valkey.GetValKey() != nil {
		valkey.GetValKey().Close()
		fmt.Println("[GO-SAIL] <Components> shutdown [valkey] successfully")
	}

	//- nacos
	if nacos.GetConfigClient() != nil {
		nacos.GetConfigClient().CloseClient()
		fmt.Println("[GO-SAIL] <Components> shutdown [nacos - configClient] successfully")
	}
	if nacos.GetNamingClient() != nil {
		nacos.GetNamingClient().CloseClient()
		fmt.Println("[GO-SAIL] <Components> shutdown [nacos - namingClient] successfully")
	}
}
