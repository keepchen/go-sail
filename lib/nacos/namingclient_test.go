package nacos

import (
	"fmt"
	"testing"

	"github.com/nacos-group/nacos-sdk-go/v2/model"

	"github.com/keepchen/go-sail/v3/lib/logger"
)

var loggerConf = logger.Conf{
	Filename: "../../examples/logs/nacos_tester.log",
}

func TestRegisterService(t *testing.T) {
	t.Run("RegisterService", func(t *testing.T) {
		t.Log(RegisterService("go-sail", "tester", "", 0, nil))
		t.Log(RegisterService("go-sail", "tester", "127.0.0.1", 8848, nil))
	})
}

func TestUnregisterService(t *testing.T) {
	t.Run("UnregisterService", func(t *testing.T) {
		t.Log(UnregisterService("go-sail", "tester", "", 0))
		t.Log(UnregisterService("go-sail", "tester", "127.0.0.1", 8848))
	})
}

func TestGetHealthyInstanceUrl(t *testing.T) {
	t.Run("GetHealthyInstanceUrl", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		t.Log(GetHealthyInstanceUrl("go-sail", "tester", logger.GetLogger()))
	})
}

func TestSubscribeInstances(t *testing.T) {
	t.Run("SubscribeInstances", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		fn := func(instances []model.Instance, err error) {
			fmt.Println(instances, err)
		}
		SubscribeInstances("go-sail", "tester", fn, logger.GetLogger())
	})
}

func TestUnsubscribeInstances(t *testing.T) {
	t.Run("UnsubscribeInstances", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		fn := func(instances []model.Instance, err error) {
			fmt.Println(instances, err)
		}
		UnsubscribeInstances("go-sail", "tester", fn, logger.GetLogger())
	})
}
