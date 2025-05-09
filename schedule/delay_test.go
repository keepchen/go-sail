package schedule

import (
	"fmt"
	"testing"
	"time"

	"github.com/keepchen/go-sail/v3/lib/redis"
)

var sConf = redis.Conf{
	Enable: true,
	Endpoint: redis.Endpoint{
		Host:     "127.0.0.1",
		Port:     6379,
		Username: "",
		Password: "",
	},
	Database: 0,
}

func TestRunAfter(t *testing.T) {
	t.Run("RunAfter", func(t *testing.T) {
		NewJob("RunAfter", func() {
			fmt.Println("RunAfter...")
		}).RunAfter(time.Second)

		time.Sleep(3 * time.Second)

		_, err := redis.New(sConf)
		if err == nil {
			redis.InitRedis(sConf)

			NewJob("RunAfter-WithoutOverlapping", func() {
				fmt.Println("RunAfter...")
			}).WithoutOverlapping().RunAfter(time.Second)
		}

		time.Sleep(3 * time.Second)

		cancel := NewJob("RunAfter2", func() {
			fmt.Println("RunAfter2...")
		}).RunAfter(time.Minute)

		cancel()
	})
}
