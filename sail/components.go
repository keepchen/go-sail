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

// GetDB 获取数据库实例
//
// 该方法依次返回读实例、写实例
//
// 注意，使用前请确保db组件已初始化成功。
func GetDB() (read *gorm.DB, write *gorm.DB) {
	if db.GetInstance() == nil {
		read, write = nil, nil
	} else {
		read, write = db.GetInstance().R, db.GetInstance().W
	}

	return
}

// NewDB 创建新的数据实例
func NewDB(conf db.Conf) (read *gorm.DB, rErr error, write *gorm.DB, wErr error) {
	return db.New(conf)
}

// GetDBR 获取数据库读实例
//
// 注意，使用前请确保db组件已初始化成功。
func GetDBR() *gorm.DB {
	if db.GetInstance() == nil {
		return nil
	}

	return db.GetInstance().R
}

// GetDBW 获取数据库写实例
//
// 注意，使用前请确保db组件已初始化成功。
func GetDBW() *gorm.DB {
	if db.GetInstance() == nil {
		return nil
	}

	return db.GetInstance().W
}

// GetRedis 获取通用的redis连接
//
// # GetRedisUniversal 方法的语法糖
//
// 自动检测当前已实例化的redis连接
//
// 如果同时存在standalone实例和cluster实例，优先返回standalone实例
//
// 如果期望获取指定的连接类型，请单独使用 GetRedisStandalone 或 GetRedisCluster
//
// ----
//
// 提示：当你能确定你的连接类型时，可以通过断言的方式调用进阶方法，例如：
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

// GetRedisUniversal 获取通用的redis连接
//
// 自动检测当前已实例化的redis连接
//
// 如果同时存在standalone实例和cluster实例，优先返回standalone实例
//
// 如果期望获取指定的连接类型，请单独使用 GetRedisStandalone 或 GetRedisCluster
//
// ----
//
// 提示：当你能确定你的连接类型时，可以通过断言的方式调用进阶方法，例如：
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

// GetRedisStandalone 获取redis连接(standalone)
//
// 单实例模式
//
// 注意，使用前请确保redis(standalone)组件已初始化成功。
func GetRedisStandalone() *redisLib.Client {
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

// NewRedis 创建新的redis (standalone)连接实例
func NewRedis(conf redis.Conf) (*redisLib.Client, error) {
	return redis.New(conf)
}

// NewRedisCluster  创建新的redis (cluster)连接实例
func NewRedisCluster(conf redis.ClusterConf) (*redisLib.ClusterClient, error) {
	return redis.NewCluster(conf)
}

// GetNats 获取nats连接实例
//
// 注意，使用前请确保nats组件已初始化成功。
func GetNats() *natsLib.Conn {
	return nats.GetInstance()
}

// NewNats 创建新的nats 连接实例
func NewNats(conf nats.Conf) (*natsLib.Conn, error) {
	return nats.New(conf)
}

// GetLogger 获取日志实例
//
// 注意，使用前请确保logger组件已初始化成功。
func GetLogger(module ...string) *zap.Logger {
	return logger.GetLogger(module...)
}

// MarshalInterfaceValue 将interface序列化成字符串
//
// 主要用于日志记录
func MarshalInterfaceValue(obj any) string {
	return logger.MarshalInterfaceValue(obj)
}

// Response http响应组件
func Response(c *gin.Context) api.Responder {
	return api.Response(c)
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

// NewKafkaConnections 创建新的kafka连接实例
func NewKafkaConnections(conf kafka.Conf) []*kafkaLib.Conn {
	return kafka.NewConnections(conf)
}

// GetEtcdInstance 获取etcd连接实例
//
// 注意，使用前请确保etcd组件已初始化成功。
func GetEtcdInstance() *clientv3.Client {
	return etcd.GetInstance()
}

// NewEtcd 创建新的etcd连接实例
func NewEtcd(conf etcd.Conf) (*clientv3.Client, error) {
	return etcd.New(conf)
}

// GetValKey 获取valkey连接实例
//
// 注意，使用前请确保valkey组件已初始化成功。
func GetValKey() valkeyLib.Client {
	return valkey.GetValKey()
}

// NewValKey 创建新的valkey连接实例
func NewValKey(conf valkey.Conf) (valkeyLib.Client, error) {
	return valkey.New(conf)
}

// GetNacosConfigClient 获取nacos配置连接实例
//
// 注意，使用前请确保nacos组件已初始化成功。
func GetNacosConfigClient() config_client.IConfigClient {
	return nacos.GetConfigClient()
}

// GetNacosNamingClient 获取nacos命名空间连接实例
//
// 注意，使用前请确保nacos组件已初始化成功。
func GetNacosNamingClient() naming_client.INamingClient {
	return nacos.GetNamingClient()
}

// NewNacosConfigClient 新建nacos配置连接实例
func NewNacosConfigClient(appName string, endpoints string, namespace string,
	clientCfg ...constant.ClientConfig) (config_client.IConfigClient, error) {
	return nacos.NewConfigClient(appName, endpoints, namespace, clientCfg...)
}

// NewNacosNamingClient 新建nacos命名空间连接实例
func NewNacosNamingClient(appName string, endpoints string, namespace string,
	clientCfg ...constant.ClientConfig) (naming_client.INamingClient, error) {
	return nacos.NewNamingClient(appName, endpoints, namespace, clientCfg...)
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
		// 可能通过config进行了初始化，因此先检测是否有全局可用变量，没有才进行初始化操作
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
	if etcd.GetInstance() != nil {
		_ = etcd.GetInstance().Close()
		fmt.Println("[GO-SAIL] <Components> shutdown [etcd] successfully")
	}

	//- valkey
	if conf.ValKeyConf.Enable && valkey.GetValKey() != nil {
		valkey.GetValKey().Close()
		fmt.Println("[GO-SAIL] <Components> shutdown [valkey] successfully")
	}
}
