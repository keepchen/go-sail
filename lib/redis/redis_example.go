package redis

import (
	"context"
	"log"
	"time"

	redisLib "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

//ExampleUsage 使用示例
//
//@see https://github.com/go-redis/redis
func ExampleUsage() {
	conf := Conf{
		Addr: Addr{
			Host:     "localhost",
			Port:     6379,
			Password: "",
		},
		Database: 0,
	}
	InitRedis(conf)

	err := GetInstance().Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := GetInstance().Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	log.Println("key", val)

	val2, err := GetInstance().Get(ctx, "key2").Result()
	if err == redisLib.Nil {
		log.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		log.Println("key2", val2)
	}

	// SET key value EX 10 NX
	set, err := GetInstance().SetNX(ctx, "key", "value", 10*time.Second).Result()
	log.Println("SetNX", set, "err", err)

	// SET key value keep ttl NX
	set, err = GetInstance().SetNX(ctx, "key", "value", redisLib.KeepTTL).Result()
	log.Println("SetNX", set, "err", err)

	// SORT list LIMIT 0 2 ASC
	val3, err := GetInstance().Sort(ctx, "list", &redisLib.Sort{Offset: 0, Count: 2, Order: "ASC"}).Result()
	log.Println("Sort", val3, "err", err)

	// ZRANGEBYSCORE zset -inf +inf WITHSCORES LIMIT 0 2
	val4, err := GetInstance().ZRangeByScoreWithScores(ctx, "zset", &redisLib.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  2,
	}).Result()
	log.Println("ZRangeByScoreWithScores", val4, "err", err)

	// ZINTERSTORE out 2 zset1 zset2 WEIGHTS 2 3 AGGREGATE SUM
	val5, err := GetInstance().ZInterStore(ctx, "out", &redisLib.ZStore{
		Keys:    []string{"zset1", "zset2"},
		Weights: []float64{2, 3},
	}).Result()
	log.Println("ZInterStore", val5, "err", err)

	// EVAL "return {KEYS[1],ARGV[1]}" 1 "key" "hello"
	val6, err := GetInstance().Eval(ctx, "return {KEYS[1],ARGV[1]}", []string{"key"}, "hello").Result()
	log.Println("Eval", val6, "err", err)

	// custom command
	res, err := GetInstance().Do(ctx, "set", "key", "value").Result()
	log.Println("Do", res, "err", err)
}
