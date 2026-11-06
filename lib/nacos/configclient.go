package nacos

import (
	"log"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// GetConfig retrieves configuration
//
// groupName group the configuration belongs to
//
// dataID configuration file ID
//
// appConfig target to parse into
//
// format configuration file format; supports: json|yaml|toml
func GetConfig(groupName, dataID string, appConfig any, format string) error {
	content, err := GetConfigClient().GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  groupName,
	})
	if err != nil {
		return err
	}

	return ParseConfig([]byte(content), appConfig, format)
}

// ListenConfig watches configuration
//
// groupName group the configuration belongs to
//
// dataID configuration file ID
//
// appConfig target to parse into
//
// format configuration file format; supports: json|yaml|toml
//
// printData prints configuration content (optional); true enables printing, so protect sensitive information
func ListenConfig(groupName, dataID string, appConfig any, format string, printData ...bool) error {
	err := GetConfigClient().ListenConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  groupName,
		OnChange: func(namespace, group, dataId, data string) {
			if len(printData) > 0 && printData[0] {
				log.Printf("[GO-SAIL] <Nacos> listen config change,data:\n%s", data)
			}
			err := ParseConfig([]byte(data), &appConfig, format)
			if err != nil {
				log.Printf("[GO-SAIL] <Nacos> listen config {%s:%s} change,but can't be unmarshal: %s\n", group, dataId, err.Error())
				return
			}
		},
	})

	return err
}

// ListenConfigWithCallback watches configuration with a custom callback
//
// groupName group the configuration belongs to
//
// dataID configuration file ID
//
// callback callback function
func ListenConfigWithCallback(groupName, dataID string, callback func(namespace, group, dataId, data string)) error {
	err := GetConfigClient().ListenConfig(vo.ConfigParam{
		DataId:   dataID,
		Group:    groupName,
		OnChange: callback,
	})

	return err
}
