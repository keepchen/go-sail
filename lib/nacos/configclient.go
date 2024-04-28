package nacos

import (
	"log"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// GetConfig 获取配置
//
// groupName 所属分组
//
// dataID 配置文件id
//
// appConfig 解析到目标
//
// format 配置文件格式，支持: json|yaml|toml
func GetConfig(groupName, dataID string, appConfig interface{}, format string) error {
	content, err := GetConfigClient().GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  groupName,
	})
	if err != nil {
		return err
	}

	return ParseConfig([]byte(content), appConfig, format)
}

// ListenConfig 监听配置
//
// groupName 所属分组
//
// dataID 配置文件id
//
// appConfig 解析到目标
//
// format 配置文件格式，支持: json|yaml|toml
//
// printData 打印配置内容(可选项)，为true将打印配置内容，注意隐私保护
func ListenConfig(groupName, dataID string, appConfig interface{}, format string, printData ...bool) error {
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

// ListenConfigWithCallback 监听配置-自定义回调函数方式
//
// groupName 所属分组
//
// dataID 配置文件id
//
// callback 回调函数
func ListenConfigWithCallback(groupName, dataID string, callback func(namespace, group, dataId, data string)) error {
	err := GetConfigClient().ListenConfig(vo.ConfigParam{
		DataId:   dataID,
		Group:    groupName,
		OnChange: callback,
	})

	return err
}
