package nacos

import (
	"encoding/json"
	"fmt"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v2"
)

// ParseConfig parses configuration content
//
// configBytes configuration content as a byte slice
//
// appConfig target to parse into
//
// format content format; supports: json|yaml|toml
func ParseConfig(configBytes []byte, appConfig any, format string) error {
	var (
		err        error
		formatList = [...]string{"json", "yaml", "toml"}
	)
	switch format {
	case formatList[0]:
		err = json.Unmarshal(configBytes, appConfig)
	case formatList[1]:
		err = yaml.Unmarshal(configBytes, appConfig)
	case formatList[2]:
		err = toml.Unmarshal(configBytes, appConfig)
	default:
		err = fmt.Errorf("[GO-SAIL] <Config> dump config by using unknown format: %s\n", format)
	}

	return err
}
