package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Watch 监听key
func Watch(key string, fn func(k, v []byte)) {
	watchChan := GetInstance().Watch(context.Background(), key)
	for watchResp := range watchChan {
		for _, value := range watchResp.Events {
			fn(value.Kv.Key, value.Kv.Value)
		}
	}
}

// WatchWithPrefix 监听key，带前缀
func WatchWithPrefix(key string, fn func(k, v []byte)) {
	watchChan := GetInstance().Watch(context.Background(), key, clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, value := range watchResp.Events {
			fn(value.Kv.Key, value.Kv.Value)
		}
	}
}
