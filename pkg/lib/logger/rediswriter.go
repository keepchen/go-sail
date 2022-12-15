package logger

import (
	"context"

	redisLib "github.com/go-redis/redis/v8"
)

type redisWriterStd struct {
	cli     *redisLib.Client
	listKey string
}

func (w *redisWriterStd) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(context.Background(), w.listKey, p).Result()

	return int(n), err
}

type redisClusterWriterStd struct {
	cli     *redisLib.ClusterClient
	listKey string
}

func (w *redisClusterWriterStd) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(context.Background(), w.listKey, p).Result()

	return int(n), err
}
