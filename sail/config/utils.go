package config

import (
	"encoding/json"
	"fmt"

	"github.com/keepchen/go-sail/v3/utils"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v2"
)

// PrintTemplateConfig 打印配置信息模板
//
// format 内容格式，支持: json|yaml|toml
//
// writeToFile 写入到目标文件（可选）
func PrintTemplateConfig(format string, writeToFile ...string) {
	var (
		abort      bool
		cfgStr     []byte
		cfg        = Config{}
		formatList = [...]string{"json", "yaml", "toml"}
	)
	switch format {
	case formatList[0]:
		cfgStr, _ = json.MarshalIndent(&cfg, "", "    ")
	case formatList[1]:
		cfgStr, _ = yaml.Marshal(&cfg)
	case formatList[2]:
		cfgStr, _ = toml.Marshal(&cfg)
	default:
		fmt.Printf("[GO-SAIL] <Config> dump config by using unknown format: %s\n", format)
		abort = true
	}

	if abort {
		return
	}

	if len(writeToFile) > 0 {
		err := utils.File().PutContents(cfgStr, writeToFile[0])
		if err != nil {
			fmt.Printf("[GO-SAIL] <Config> dump config to file {%s} error: %s\n", writeToFile[0], err.Error())
		}
	} else {
		fmt.Printf("[GO-SAIL] <Config> dump config (%s) to stdout:\n", format)
		fmt.Println(string(cfgStr))
	}
}

// ParseConfigFromBytes 从字符串解析配置
//
// format 内容格式，支持: json|yaml|toml
//
// source 配置源字符
func ParseConfigFromBytes(format string, source []byte) (*Config, error) {
	var (
		formatList = [...]string{"json", "yaml", "toml"}
		cfg        Config
		err        error
	)

	switch format {
	case formatList[0]:
		err = json.Unmarshal(source, &cfg)
	case formatList[1]:
		err = yaml.Unmarshal(source, &cfg)
	case formatList[2]:
		err = toml.Unmarshal(source, &cfg)
	default:
		fmt.Printf("[GO-SAIL] <Config> dump config by using unknown format: %s\n", format)
		err = fmt.Errorf("[GO-SAIL] <Config> dump config by using unknown format: %s\n", format)
	}

	return &cfg, err
}

// ParseConfigFromBytesToDst 从字符串解析配置到目标结构
//
// format 内容格式，支持: json|yaml|toml
//
// source 配置源字符
//
// dst 目标结构
//
// @return 返回值与参数dst类型相同，需要进行断言
func ParseConfigFromBytesToDst(format string, source []byte, dst interface{}) (interface{}, error) {
	var (
		formatList = [...]string{"json", "yaml", "toml"}
		err        error
	)

	switch format {
	case formatList[0]:
		err = json.Unmarshal(source, dst)
	case formatList[1]:
		err = yaml.Unmarshal(source, dst)
	case formatList[2]:
		err = toml.Unmarshal(source, dst)
	default:
		fmt.Printf("[GO-SAIL] <Config> dump config by using unknown format: %s\n", format)
		err = fmt.Errorf("[GO-SAIL] <Config> dump config by using unknown format: %s\n", format)
	}

	return dst, err
}
