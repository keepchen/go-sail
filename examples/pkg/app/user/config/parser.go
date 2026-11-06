package config

import (
	"fmt"
	"log"

	"github.com/jinzhu/configor"
	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/lib/nacos"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	yaml "gopkg.in/yaml.v2"
)

// ParseConfig parses configuration
//
// Parse configuration from a file when cfgPath is nonempty; otherwise read from Nacos
func ParseConfig(cfgPath string) {
	if len(cfgPath) != 0 {
		fmt.Printf("Parse config from local file: %s\n", cfgPath)
		ParseConfigFromFile(cfgPath)
	} else {
		fmt.Println("Parse config from nacos")
		ParseConfigFromNacos()
	}

	//Test loading the time zone
	TestLoadTimeZonePanicWhileFailure()
}

func ParseConfigFromFile(cfgPath string) {
	if len(cfgPath) != 0 {
		if err := configor.New(&configor.Config{AutoReload: true}).Load(globalConfig, cfgPath); err != nil {
			panic(err)
		}
	} else {
		if err := configor.New(&configor.Config{AutoReload: true}).Load(globalConfig); err != nil {
			panic(err)
		}
	}

	//Parse JWT configuration
	globalConfig.JWT.Load()

	//Set the default time zone
	if len(globalConfig.Timezone) == 0 {
		globalConfig.Timezone = constants.DefaultTimeZone
	}

	if globalConfig.Debug {
		log.Printf("loaded config: %#v", globalConfig)
	}
}

// ParseConfigFromNacos reads and parses configuration from Nacos
//
// Note: the nacosAddrs and nacosNamespaceID environment variables must be set
//
// For example:
//
// nacosAddrs=192.168.224.3:8848,192.168.224.3:8848
//
// nacosNamespaceID=f0c2f58d-54e3-45df-94cd-a6fb3a8fa534
func ParseConfigFromNacos() {
	var (
		dataID = "go-sail-user.yml"
		group  = "go-sail"
	)
	cc := nacos.GetConfigClient()

	content, err := cc.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
	})
	if err != nil {
		panic(err)
	}

	var conf Config
	err = yaml.Unmarshal([]byte(content), &conf)
	if err != nil {
		panic(err)
	}

	globalConfig = &conf

	//Set the default time zone
	if len(globalConfig.Timezone) == 0 {
		globalConfig.Timezone = constants.DefaultTimeZone
	}

	if globalConfig.Debug {
		log.Printf("loaded config: %#v", globalConfig)
	}

	err = cc.ListenConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			log.Printf("listen config change,data:\n%s", data)
			var cfg Config
			err = yaml.Unmarshal([]byte(data), &cfg)
			if err != nil {
				log.Printf("listen config change,but can't be yaml unmarshal: %s", err.Error())
				return
			}
			globalConfig = &cfg

			//Set the default time zone
			if len(globalConfig.Timezone) == 0 {
				globalConfig.Timezone = constants.DefaultTimeZone
			}
		},
	})
}
