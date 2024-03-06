package redis

import (
	"context"
	"log"
	"time"

	redisLib "github.com/go-redis/redis/v8"
)

var ctx2 = context.Background()

// ExampleClusterUsage 使用示例
//
// @see https://github.com/go-redis/redis
func ExampleClusterUsage() {
	conf := ClusterConf{
		SSLEnable: false,
		Endpoints: []Endpoint{
			{
				Host:     "localhost",
				Port:     6379,
				Password: "",
			},
			{
				Host:     "localhost",
				Port:     6380,
				Password: "",
			},
			{
				Host:     "localhost",
				Port:     6381,
				Password: "",
			},
		},
	}
	InitRedisCluster(conf)

	err := GetClusterInstance().Set(ctx2, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := GetClusterInstance().Get(ctx2, "key").Result()
	if err != nil {
		panic(err)
	}
	log.Println("key", val)

	val2, err := GetClusterInstance().Get(ctx2, "key2").Result()
	if err == redisLib.Nil {
		log.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		log.Println("key2", val2)
	}

	// SET key value EX 10 NX
	set, err := GetClusterInstance().SetNX(ctx2, "key", "value", 10*time.Second).Result()
	log.Println("SetNX", set, "err", err)

	// SET key value keep ttl NX
	set, err = GetClusterInstance().SetNX(ctx2, "key", "value", redisLib.KeepTTL).Result()
	log.Println("SetNX", set, "err", err)

	// SORT list LIMIT 0 2 ASC
	val3, err := GetClusterInstance().Sort(ctx2, "list", &redisLib.Sort{Offset: 0, Count: 2, Order: "ASC"}).Result()
	log.Println("Sort", val3, "err", err)

	// ZRANGEBYSCORE zset -inf +inf WITHSCORES LIMIT 0 2
	val4, err := GetClusterInstance().ZRangeByScoreWithScores(ctx2, "zset", &redisLib.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  2,
	}).Result()
	log.Println("ZRangeByScoreWithScores", val4, "err", err)

	// ZINTERSTORE out 2 zset1 zset2 WEIGHTS 2 3 AGGREGATE SUM
	val5, err := GetClusterInstance().ZInterStore(ctx2, "out", &redisLib.ZStore{
		Keys:    []string{"zset1", "zset2"},
		Weights: []float64{2, 3},
	}).Result()
	log.Println("ZInterStore", val5, "err", err)

	// EVAL "return {KEYS[1],ARGV[1]}" 1 "key" "hello"
	val6, err := GetClusterInstance().Eval(ctx2, "return {KEYS[1],ARGV[1]}", []string{"key"}, "hello").Result()
	log.Println("Eval", val6, "err", err)

	// custom command
	res, err := GetClusterInstance().Do(ctx2, "set", "key", "value").Result()
	log.Println("Do", res, "err", err)
}
