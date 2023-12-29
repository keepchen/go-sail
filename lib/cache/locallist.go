package cache

import "github.com/keepchen/message-queue/queue"

var listInstance *queue.Instance

// InitList 初始化列表
func InitList() {
	listInstance = queue.GetDBInstance("<localList>")
	listInstance.SetDebugMode(false)
}

// GetListInstance 获取列表实例
func GetListInstance() *queue.Instance {
	return listInstance
}

// NewList 新建列表
func NewList(listName string) *queue.Instance {
	instance := queue.GetDBInstance(listName)
	instance.SetDebugMode(false)

	return instance
}
