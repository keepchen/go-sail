package logger

import (
	"context"

	redisLib "github.com/go-redis/redis/v8"
)

type redisWriterStd struct {
	cli     *redisLib.Client
	listKey string
}

func (w *redisWriterStd) Sync() error {
	//TODO implement me

	return nil
}

func (w *redisWriterStd) Write(p []byte) (int, error) {
	_, err := w.cli.RPush(context.Background(), w.listKey, p).Result()

	return len(p), err
}

type redisClusterWriterStd struct {
	cli     *redisLib.ClusterClient
	listKey string
}

func (w *redisClusterWriterStd) Sync() error {
	//TODO implement me

	return nil
}

func (w *redisClusterWriterStd) Write(p []byte) (int, error) {
	_, err := w.cli.RPush(context.Background(), w.listKey, p).Result()

	return len(p), err
}
