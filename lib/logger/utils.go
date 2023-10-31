package logger

import (
	"encoding/json"
)

//MarshalInterfaceValue 将interface序列化成字符串
//
//主要用于日志记录
func MarshalInterfaceValue(obj interface{}) string {
	js, _ := json.Marshal(obj)

	return string(js)
}
