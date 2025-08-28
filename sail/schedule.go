package sail

import (
	"github.com/keepchen/go-sail/v3/schedule"
)

// Schedule 计划任务
//
// # Note
//
// 使用默认的redis连接实例作为竞态检测的实例连接，
// 调用者可以使用 SetRedisClient(client).ForSchedule() 覆盖默认连接
func Schedule(name string, task func()) schedule.Scheduler {
	return schedule.NewJob(name, task)
}
