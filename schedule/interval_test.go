package schedule

import (
	"fmt"
	"github.com/keepchen/go-sail/v3/lib/redis"
	"net"
	"testing"
	"time"
)

func TestEvery(t *testing.T) {
	t.Run("Every", func(t *testing.T) {
		cancel := NewJob("Every", func() {
			fmt.Println("Every...")
		}).Every(time.Second)

		cancel()
	})
}

func TestEverySecond(t *testing.T) {
	t.Run("EverySecond", func(t *testing.T) {
		cancel := NewJob("EverySecond", func() {
			fmt.Println("EverySecond...")
		}).EverySecond()

		cancel()
	})
}

func TestEveryFiveSeconds(t *testing.T) {
	t.Run("EveryFiveSeconds", func(t *testing.T) {
		cancel := NewJob("EveryFiveSeconds", func() {
			fmt.Println("EveryFiveSeconds...")
		}).EveryFiveSeconds()

		cancel()
	})
}

func TestEveryTenSeconds(t *testing.T) {
	t.Run("EveryTenSeconds", func(t *testing.T) {
		cancel := NewJob("EveryTenSeconds", func() {
			fmt.Println("EveryTenSeconds...")
		}).EveryTenSeconds()

		cancel()
	})
}

func TestEveryFifteenSeconds(t *testing.T) {
	t.Run("EveryFifteenSeconds", func(t *testing.T) {
		cancel := NewJob("EveryFifteenSeconds", func() {
			fmt.Println("EveryFifteenSeconds...")
		}).EveryFifteenSeconds()

		cancel()
	})
}

func TestEveryTwentySeconds(t *testing.T) {
	t.Run("EveryTwentySeconds", func(t *testing.T) {
		cancel := NewJob("EveryTwentySeconds", func() {
			fmt.Println("EveryTwentySeconds...")
		}).EveryTwentySeconds()

		cancel()
	})
}

func TestEveryThirtySeconds(t *testing.T) {
	t.Run("EveryThirtySeconds", func(t *testing.T) {
		cancel := NewJob("EveryThirtySeconds", func() {
			fmt.Println("EveryThirtySeconds...")
		}).EveryThirtySeconds()

		cancel()
	})
}

func TestEveryMinute(t *testing.T) {
	t.Run("EveryMinute", func(t *testing.T) {
		cancel := NewJob("EveryMinute", func() {
			fmt.Println("EveryMinute...")
		}).EveryMinute()

		cancel()
	})
}

func TestEveryFiveMinutes(t *testing.T) {
	t.Run("EveryFiveMinutes", func(t *testing.T) {
		cancel := NewJob("EveryFiveMinutes", func() {
			fmt.Println("EveryFiveMinutes...")
		}).EveryFiveMinutes()

		cancel()
	})
}

func TestEveryTenMinutes(t *testing.T) {
	t.Run("EveryTenMinutes", func(t *testing.T) {
		cancel := NewJob("EveryTenMinutes", func() {
			fmt.Println("EveryTenMinutes...")
		}).EveryTenMinutes()

		cancel()
	})
}

func TestEveryFifteenMinutes(t *testing.T) {
	t.Run("EveryFifteenMinutes", func(t *testing.T) {
		cancel := NewJob("EveryFifteenMinutes", func() {
			fmt.Println("EveryFifteenMinutes...")
		}).EveryFifteenMinutes()

		cancel()
	})
}

func TestEveryTwentyMinutes(t *testing.T) {
	t.Run("EveryTwentyMinutes", func(t *testing.T) {
		cancel := NewJob("EveryTwentyMinutes", func() {
			fmt.Println("EveryTwentyMinutes...")
		}).EveryTwentyMinutes()

		cancel()
	})
}

func TestEveryThirtyMinutes(t *testing.T) {
	t.Run("EveryThirtyMinutes", func(t *testing.T) {
		cancel := NewJob("EveryThirtyMinutes", func() {
			fmt.Println("EveryThirtyMinutes...")
		}).EveryThirtyMinutes()

		cancel()
	})
}

func TestHourly(t *testing.T) {
	t.Run("Hourly", func(t *testing.T) {
		cancel := NewJob("Hourly", func() {
			fmt.Println("Hourly...")
		}).Hourly()

		cancel()
	})
}

func TestEveryFiveHours(t *testing.T) {
	t.Run("EveryFiveHours", func(t *testing.T) {
		cancel := NewJob("EveryFiveHours", func() {
			fmt.Println("EveryFiveHours...")
		}).EveryFiveHours()

		cancel()
	})
}

func TestEveryTenHours(t *testing.T) {
	t.Run("EveryTenHours", func(t *testing.T) {
		cancel := NewJob("EveryTenHours", func() {
			fmt.Println("EveryTenHours...")
		}).EveryTenHours()

		cancel()
	})
}

func TestEveryTwentyHours(t *testing.T) {
	t.Run("EveryTwentyHours", func(t *testing.T) {
		cancel := NewJob("EveryTwentyHours", func() {
			fmt.Println("EveryTwentyHours...")
		}).EveryTwentyHours()

		cancel()
	})
}

func TestDaily(t *testing.T) {
	t.Run("Daily", func(t *testing.T) {
		cancel := NewJob("Daily", func() {
			fmt.Println("Daily...")
		}).Daily()

		cancel()
	})
}

func TestWeekly(t *testing.T) {
	t.Run("Weekly", func(t *testing.T) {
		cancel := NewJob("Weekly", func() {
			fmt.Println("Weekly...")
		}).Weekly()

		cancel()
	})
}

func TestMonthly(t *testing.T) {
	t.Run("Monthly", func(t *testing.T) {
		cancel := NewJob("Monthly", func() {
			fmt.Println("Monthly...")
		}).Monthly()

		cancel()
	})
}

func TestYearly(t *testing.T) {
	t.Run("Yearly", func(t *testing.T) {
		cancel := NewJob("Yearly", func() {
			fmt.Println("Yearly...")
		}).Yearly()

		cancel()
	})
}

func TestWithoutOverlapping(t *testing.T) {
	t.Run("WithoutOverlapping", func(t *testing.T) {
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
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", sConf.Host, sConf.Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		redis.InitRedis(sConf)

		cancel := NewJob("Yearly", func() {
			fmt.Println("Yearly...")
		}).WithoutOverlapping().EverySecond()

		cancel()
	})
}
