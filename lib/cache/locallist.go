package cache

import "github.com/keepchen/message-queue/queue"

var listInstance *queue.Instance

// InitList initializes the list
func InitList() {
	listInstance = queue.GetDBInstance("<localList>")
	listInstance.SetDebugMode(false)
}

// GetListInstance retrieves the list instance
func GetListInstance() *queue.Instance {
	return listInstance
}

// NewList creates a list
func NewList(listName string) *queue.Instance {
	instance := queue.GetDBInstance(listName)
	instance.SetDebugMode(false)

	return instance
}
