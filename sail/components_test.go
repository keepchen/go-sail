package sail

import (
	"testing"

	"github.com/keepchen/go-sail/v3/lib/valkey"

	"github.com/keepchen/go-sail/v3/lib/etcd"

	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/kafka"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/lib/nats"
	"github.com/keepchen/go-sail/v3/lib/redis"
	"github.com/keepchen/go-sail/v3/sail/config"
)

func TestGetDB(t *testing.T) {
	t.Run("GetDB", func(t *testing.T) {
		t.Log(GetDB())
	})
}

func TestGetDBR(t *testing.T) {
	t.Run("GetDBR", func(t *testing.T) {
		t.Log(GetDBR())
	})
}

func TestGetDBW(t *testing.T) {
	t.Run("GetDBW", func(t *testing.T) {
		t.Log(GetDBW())
	})
}

func TestNewDB(t *testing.T) {
	t.Run("NewDB", func(t *testing.T) {
		var loggerCfg = logger.Conf{
			Level:    "debug",
			Filename: "../examples/logs/testcase_components.log",
		}
		logger.Init(loggerCfg, "tester")
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
	t.Run("GetRedisUniversal", func(t *testing.T) {
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
}

func TestGetKafkaWriter(t *testing.T) {
	t.Run("GetKafkaWriter", func(t *testing.T) {
		if kafka.GetInstance() == nil {
			t.Log(kafka.GetInstance())
			return
		}
		t.Log(GetKafkaWriter())
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
}

func TestComponentsShutdown(t *testing.T) {
	t.Run("ComponentsShutdown", func(t *testing.T) {
		cfg := &config.Config{}
		componentsShutdown(cfg)
	})
}
