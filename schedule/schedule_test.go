package schedule

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/keepchen/go-sail/v3/lib/redis"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJobNameKey(t *testing.T) {
	t.Log(generateJobNameKey("tester"))
}

func TestNewJob(t *testing.T) {
	t.Run("NewJob", func(t *testing.T) {
		scheduler := NewJob("NewJob", func() {
			fmt.Println("NewJob...")
		})
		t.Log(scheduler)
	})
}

func TestJob(t *testing.T) {
	t.Run("Job", func(t *testing.T) {
		scheduler := Job("Job", func() {
			fmt.Println("Job...")
		})
		t.Log(scheduler)
	})

	t.Run("Job-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			scheduler := Job("Job-Panic", func() {
				fmt.Println("Job-Panic...")
			})
			t.Log(scheduler)
			Job("Job-Panic", func() {
				fmt.Println("Job-Panic...")
			})

		})
	})
}

func TestJobIsRunning(t *testing.T) {
	t.Run("JobIsRunning", func(t *testing.T) {
		Job("JobIsRunning", func() {
			fmt.Println("JobIsRunning...")
		})
		t.Log(JobIsRunning("JobIsRunning"))
	})
}

func TestCall(t *testing.T) {
	t.Run("Call", func(t *testing.T) {
		Job("Call", func() {
			fmt.Println("Call...")
		})
		Call("Call", false)
		time.Sleep(time.Second * 2)
	})

	t.Run("Call-NoExist", func(t *testing.T) {
		Call("Call-NoExist", false)
	})

	t.Run("Call-mandatory-true", func(t *testing.T) {
		Job("Call-mandatory-true", func() {
			fmt.Println("Call-mandatory-true...")
		})
		Call("Call-mandatory-true", true)
		time.Sleep(time.Second * 2)
	})

	t.Run("Call-mandatory-false", func(t *testing.T) {
		Job("Call-mandatory-false", func() {
			fmt.Println("Call-mandatory-false...")
		})
		Call("Call-mandatory-false", false)
		time.Sleep(time.Second * 2)
	})
}

func TestMustCall(t *testing.T) {
	t.Run("MustCall", func(t *testing.T) {
		scheduler := NewJob("MustCall", func() {
			fmt.Println("MustCall...")
		})
		t.Log(scheduler)
		MustCall("MustCall", true)
		time.Sleep(time.Second * 2)
	})
}

func TestSetRedisClientOnce(t *testing.T) {
	t.Run("Run-SetRedisClientOnce", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", sConf.Host, sConf.Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		redisClient, err := redis.New(sConf)

		assert.NoError(t, err)

		SetRedisClientOnce(redisClient)
		//multiple set
		SetRedisClientOnce(redisClient)
		SetRedisClientOnce(redisClient)
	})
}
