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
// @param format 内容格式，支持: json|yaml|toml
//
// @param writeToFile 写入到目标文件（可选）
func PrintTemplateConfig(format string, writeToFile ...string) {
	var (
		abort      bool
		configStr  []byte
		config     = Config{}
		formatList = [...]string{"json", "yaml", "toml"}
	)
	switch format {
	case formatList[0]:
		configStr, _ = json.MarshalIndent(&config, "", "    ")
	case formatList[1]:
		configStr, _ = yaml.Marshal(&config)
	case formatList[2]:
		configStr, _ = toml.Marshal(&config)
	default:
		fmt.Printf("[GO-SAIL] <Config> dump config by using unknown format: %s\n", format)
		abort = true
	}

	if abort {
		return
	}

	if len(writeToFile) > 0 {
		err := utils.FilePutContents(configStr, writeToFile[0])
		if err != nil {
			fmt.Printf("[GO-SAIL] <Config> dump config to file {%s} error: %s\n", writeToFile[0], err.Error())
		}
	} else {
		fmt.Printf("[GO-SAIL] <Config> dump config (%s) to stdout:\n", format)
		fmt.Println(string(configStr))
	}
}

// ParseConfigFromBytes 从字符串解析配置
//
// @param format 内容格式，支持: json|yaml|toml
//
// @param writeToFile 写入到目标文件（可选）
func ParseConfigFromBytes(format string, source []byte) (*Config, error) {
	var (
		formatList = [...]string{"json", "yaml", "toml"}
		conf       Config
		err        error
	)

	switch format {
	case formatList[0]:
		err = json.Unmarshal(source, &conf)
	case formatList[1]:
		err = yaml.Unmarshal(source, &conf)
	case formatList[2]:
		err = toml.Unmarshal(source, &conf)
	default:
		fmt.Printf("[GO-SAIL] <Config> dump config by using unknown format: %s\n", format)
		err = fmt.Errorf("[GO-SAIL] <Config> dump config by using unknown format: %s\n", format)
	}

	return &conf, err
}
