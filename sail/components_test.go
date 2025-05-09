package sail

import (
	"net/url"
	"testing"

	"github.com/keepchen/go-sail/v3/constants"

	"github.com/keepchen/go-sail/v3/lib/jwt"

	"github.com/stretchr/testify/assert"

	"github.com/keepchen/go-sail/v3/lib/valkey"

	"github.com/keepchen/go-sail/v3/lib/etcd"

	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/kafka"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/lib/nats"
	"github.com/keepchen/go-sail/v3/lib/redis"
	"github.com/keepchen/go-sail/v3/sail/config"
)

var (
	loggerConf = logger.Conf{
		Level:    "debug",
		Filename: "../examples/logs/testcase_components.log",
	}
	dbConf = db.Conf{
		Enable:      true,
		DriverName:  "mysql",
		AutoMigrate: true,
		Logger: db.Logger{
			Level: "debug",
		},
		ConnectionPool: db.ConnectionPoolConf{
			MaxOpenConnCount:       10,
			MaxIdleConnCount:       3,
			ConnMaxLifeTimeMinutes: 30,
			ConnMaxIdleTimeMinutes: 10,
		},
		Mysql: db.MysqlConf{
			Read: db.MysqlConfItem{
				Host:      "127.0.0.1",
				Port:      33060,
				Username:  "root",
				Password:  "root",
				Database:  "go_sail",
				Charset:   "utf8mb4",
				ParseTime: true,
				Loc:       url.QueryEscape(constants.TimeZoneUTCPlus7),
			},
			Write: db.MysqlConfItem{
				Host:      "127.0.0.1",
				Port:      33060,
				Username:  "root",
				Password:  "root",
				Database:  "go_sail",
				Charset:   "utf8mb4",
				ParseTime: true,
				Loc:       url.QueryEscape(constants.TimeZoneUTCPlus7),
			},
		},
	}
	sConf = redis.Conf{
		Enable: true,
		Endpoint: redis.Endpoint{
			Host:     "127.0.0.1",
			Port:     6379,
			Username: "",
			Password: "",
		},
		Database: 0,
	}
	cConf = redis.ClusterConf{
		Enable: true,
		Endpoints: []redis.Endpoint{
			{Host: "127.0.0.1", Port: 7000},
			{Host: "127.0.0.1", Port: 7001},
			{Host: "127.0.0.1", Port: 7002},
			//{Host: "127.0.0.1", Port: 7003},
			//{Host: "127.0.0.1", Port: 7004},
			//{Host: "127.0.0.1", Port: 7005},
		},
	}
	vConf = valkey.Conf{
		Enable: true,
		Endpoints: []valkey.Endpoint{
			{Host: "127.0.0.1", Port: 8000},
			{Host: "127.0.0.1", Port: 8001},
			{Host: "127.0.0.1", Port: 8002},
			{Host: "127.0.0.1", Port: 8003},
			{Host: "127.0.0.1", Port: 8004},
			{Host: "127.0.0.1", Port: 8005},
		},
	}
	eConf = etcd.Conf{
		Enable:    true,
		Endpoints: []string{"127.0.0.1:2379"},
	}
)

func TestGetDB(t *testing.T) {
	t.Run("GetDB-Nil", func(t *testing.T) {
		t.Log(GetDB())
	})

	t.Run("GetDB", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail-tester")
		_, e1, _, e2 := db.NewFreshDB(dbConf)
		if e1 != nil || e2 != nil {
			return
		}
		db.Init(dbConf)
		t.Log(GetDB())
	})
}

func TestGetDBR(t *testing.T) {
	t.Run("GetDBR-Nil", func(t *testing.T) {
		t.Log(GetDBR())
	})

	t.Run("GetDBR", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail-tester")
		_, e1, _, e2 := db.NewFreshDB(dbConf)
		if e1 != nil || e2 != nil {
			return
		}
		db.Init(dbConf)
		t.Log(GetDBR())
	})
}

func TestGetDBW(t *testing.T) {
	t.Run("GetDBW-Nil", func(t *testing.T) {
		t.Log(GetDBW())
	})

	t.Run("GetDBW", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail-tester")
		_, e1, _, e2 := db.NewFreshDB(dbConf)
		if e1 != nil || e2 != nil {
			return
		}
		db.Init(dbConf)
		t.Log(GetDBW())
	})
}

func TestNewDB(t *testing.T) {
	t.Run("NewDB", func(t *testing.T) {
		logger.Init(loggerConf, "tester")
		cfg := db.Conf{
			DriverName: "mysql",
		}
		t.Log(NewDB(cfg))
	})
}

func TestGetRedis(t *testing.T) {
	t.Run("GetRedis", func(t *testing.T) {
		t.Log(GetRedis())
	})
}

func TestGetRedisUniversal(t *testing.T) {
	t.Run("GetRedisUniversal-Nil", func(t *testing.T) {
		t.Log(GetRedisUniversal())
	})

	t.Run("GetRedisUniversal-Standalone", func(t *testing.T) {
		_, err := redis.New(sConf)
		if err != nil {
			return
		}
		redis.InitRedis(sConf)
		t.Log(GetRedisUniversal())
	})

	t.Run("GetRedisUniversal-Cluster", func(t *testing.T) {
		_, err := redis.NewCluster(cConf)
		if err != nil {
			return
		}
		redis.InitRedisCluster(cConf)
		t.Log(GetRedisUniversal())
	})
}

func TestGetRedisStandalone(t *testing.T) {
	t.Run("GetRedisStandalone", func(t *testing.T) {
		t.Log(GetRedisStandalone())
	})
}

func TestGetRedisCluster(t *testing.T) {
	t.Run("GetRedisCluster", func(t *testing.T) {
		t.Log(GetRedisCluster())
	})
}

func TestNewRedis(t *testing.T) {
	t.Run("NewRedis", func(t *testing.T) {
		cfg := redis.Conf{}
		t.Log(NewRedis(cfg))
	})
}

func TestNewRedisCluster(t *testing.T) {
	t.Run("NewRedisCluster", func(t *testing.T) {
		cfg := redis.ClusterConf{}
		t.Log(NewRedisCluster(cfg))
	})
}

func TestGetNats(t *testing.T) {
	t.Run("GetNats", func(t *testing.T) {
		t.Log(GetNats())
	})
}

func TestNewNats(t *testing.T) {
	t.Run("NewNats", func(t *testing.T) {
		cfg := nats.Conf{}
		t.Log(NewNats(cfg))
	})
}

func TestGetLogger(t *testing.T) {
	t.Run("GetLogger", func(t *testing.T) {
		cfg := logger.Conf{}
		logger.Init(cfg, "tester")
		t.Log(GetLogger())
	})
}

func TestMarshalInterfaceValue(t *testing.T) {
	t.Run("MarshalInterfaceValue", func(t *testing.T) {
		cfg := logger.Conf{}
		t.Log(MarshalInterfaceValue(cfg))
	})
}

func TestGetKafkaInstance(t *testing.T) {
	t.Run("GetKafkaInstance", func(t *testing.T) {
		t.Log(GetKafkaInstance())
	})
}

func TestGetKafkaConnections(t *testing.T) {
	t.Run("GetKafkaConnections", func(t *testing.T) {
		if kafka.GetInstance() == nil {
			t.Log(kafka.GetInstance())
			return
		}
		t.Log(GetKafkaConnections())
	})

	t.Run("GetKafkaConnections-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			t.Log(GetKafkaConnections())
		})
	})
}

func TestGetKafkaWriter(t *testing.T) {
	t.Run("GetKafkaWriter", func(t *testing.T) {
		if kafka.GetInstance() == nil {
			t.Log(kafka.GetInstance())
			return
		}
		t.Log(GetKafkaWriter())
	})

	t.Run("GetKafkaWriter-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			t.Log(GetKafkaWriter())
		})
	})
}

func TestGetKafkaReader(t *testing.T) {
	t.Run("GetKafkaReader", func(t *testing.T) {
		if kafka.GetInstance() == nil {
			t.Log(kafka.GetInstance())
			return
		}
		t.Log(GetKafkaReader())
	})

	t.Run("GetKafkaReader-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			t.Log(GetKafkaReader())
		})
	})
}

func TestNewKafkaConnections(t *testing.T) {
	t.Run("NewKafkaConnections", func(t *testing.T) {
		cfg := kafka.Conf{}
		t.Log(NewKafkaConnections(cfg))
	})
}

func TestGetEtcdInstance(t *testing.T) {
	t.Run("GetEtcdInstance", func(t *testing.T) {
		t.Log(GetEtcdInstance())
	})
}

func TestNewEtcd(t *testing.T) {
	t.Run("NewEtcd", func(t *testing.T) {
		cfg := etcd.Conf{}
		t.Log(NewEtcd(cfg))
	})
}

func TestGetValKey(t *testing.T) {
	t.Run("GetValKey", func(t *testing.T) {
		t.Log(GetValKey())
	})
}

func TestNewValKey(t *testing.T) {
	t.Run("NewValKey", func(t *testing.T) {
		cfg := valkey.Conf{}
		t.Log(NewValKey(cfg))
	})
}

func TestComponentsStartup(t *testing.T) {
	t.Run("ComponentsStartup", func(t *testing.T) {
		cfg := &config.Config{}
		componentsStartup("tester", cfg)
	})

	t.Run("ComponentsStartup-Enable", func(t *testing.T) {
		cfg := &config.Config{
			RedisConf:        sConf,
			RedisClusterConf: cConf,
			DBConf:           dbConf,
			JwtConf: &jwt.Conf{
				Enable: true,
			},
			NatsConf: nats.Conf{
				Enable: true,
			},
			KafkaConf: config.KafkaExtraConf{
				Conf: kafka.Conf{
					Enable: true,
				},
			},
			EtcdConf:   eConf,
			ValKeyConf: vConf,
		}
		componentsStartup("tester", cfg)
	})
}

func TestComponentsShutdown(t *testing.T) {
	t.Run("ComponentsShutdown", func(t *testing.T) {
		cfg := &config.Config{}
		componentsShutdown(cfg)
	})

	t.Run("ComponentsShutdown-Enable", func(t *testing.T) {
		cfg := &config.Config{
			RedisConf:        sConf,
			RedisClusterConf: cConf,
			DBConf:           dbConf,
			JwtConf: &jwt.Conf{
				Enable: true,
			},
			NatsConf: nats.Conf{
				Enable: true,
			},
			KafkaConf: config.KafkaExtraConf{
				Conf: kafka.Conf{
					Enable: true,
				},
			},
			EtcdConf:   eConf,
			ValKeyConf: vConf,
		}
		componentsShutdown(cfg)
	})
}
