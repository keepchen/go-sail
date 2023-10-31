package nacos

import (
	"log"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// ListenConfig 监听配置
//
// @param groupName 所属分组
//
// @param dataID 配置文件id
//
// @param appConfig 解析到目标
//
// @param format 配置文件格式，支持: json|yaml|toml
func ListenConfig(groupName, dataID string, appConfig interface{}, format string) error {
	err := GetConfigClient().ListenConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  groupName,
		OnChange: func(namespace, group, dataId, data string) {
			log.Printf("[GO-SAIL] <Nacos> listen config change,data:\n%s", data)
			err := ParseConfig([]byte(data), &appConfig, format)
			if err != nil {
				log.Printf("[GO-SAIL] <Nacos> listen config {%s:%s} change,but can't be unmarshal: %s\n", group, dataId, err.Error())
				return
			}
		},
	})

	return err
}
