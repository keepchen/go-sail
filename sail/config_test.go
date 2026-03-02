package sail

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/keepchen/go-sail/v3/lib/etcd"
	"github.com/keepchen/go-sail/v3/utils"
	"github.com/stretchr/testify/assert"
)

var (
	testConfigFilename = "test-config-file.json"
	testConfigContent  = `{"appName":"go-sail"}`
	watcherFunc        = func(configName string, content []byte, isWatch bool) {
		fmt.Println("watcher: ", configName, string(content), isWatch)
	}
)

func TestConfig(t *testing.T) {
	t.Run("Config", func(t *testing.T) {
		t.Log(Config(nil))
	})
}

func TestConfigViaFile(t *testing.T) {
	t.Run("ConfigViaFile-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Config(nil).ViaFile("this-file-not-exist.json")
		})
	})
	t.Run("ConfigViaFile-Truly-File", func(t *testing.T) {
		_ = utils.File().PutContents([]byte(testConfigContent), testConfigFilename)
		defer func() {
			_ = os.Remove(testConfigFilename)
		}()
		Config(watcherFunc).ViaFile(testConfigFilename)

		Config(watcherFunc).ViaFile(testConfigFilename).Parse(watcherFunc)
	})
}

var etcdConf = etcd.Conf{
	Enable:    true,
	Endpoints: []string{"127.0.0.1:2379"},
	//Username:  "root",
	//Password:  "changeMe",
}

func TestConfigViaEtcd(t *testing.T) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s", etcdConf.Endpoints[0]))
	if err != nil {
		return
	}
	_ = conn.Close()

	//写入测试数据
	instance, err := etcd.New(etcdConf)
	t.Log(instance, err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	t.Log("----Put----")
	t.Log(instance.Put(ctx, testConfigFilename, testConfigContent))

	t.Run("ConfigViaEtcd-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Config(nil).ViaEtcd(etcdConf, "this-file-not-exist.json")
		})
	})
	t.Run("ConfigViaEtcd-Truly", func(t *testing.T) {
		Config(watcherFunc).ViaEtcd(etcdConf, testConfigFilename)

		Config(watcherFunc).ViaEtcd(etcdConf, testConfigFilename).Parse(watcherFunc)

		Config(watcherFunc, true).ViaEtcd(etcdConf, testConfigFilename)
	})
}

func TestConfigViaNacos(t *testing.T) {
	t.Run("ConfigViaNacos-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Config(nil).ViaNacos("127.0.0.1:8848", "", "", "this-file-not-exist.json")
		})
	})
}
