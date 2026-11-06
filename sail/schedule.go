package sail

import (
	"github.com/keepchen/go-sail/v3/schedule"
)

// Schedule scheduled tasks
//
// # Note
//
// Use the default Redis connection instance for contention detection;
// callers can override it with SetRedisClient(client).ForSchedule()
func Schedule(name string, task func()) schedule.Scheduler {
	return schedule.NewJob(name, task)
}
