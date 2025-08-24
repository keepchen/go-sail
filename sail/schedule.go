package sail

import (
	redisLib "github.com/go-redis/redis/v8"
	"github.com/keepchen/go-sail/v3/schedule"
)

// SetScheduleRedisClient 设置计划任务使用的redis客户端
func SetScheduleRedisClient(client redisLib.UniversalClient) {
	schedule.SetRedisClientOnce(client)
}

// Schedule 计划任务
func Schedule(name string, task func()) schedule.Scheduler {
	return schedule.NewJob(name, task)
}
