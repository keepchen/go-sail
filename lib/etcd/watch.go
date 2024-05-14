package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Watch 监听key
func Watch(ctx context.Context, key string, fn func(k, v []byte)) {
	watchChan := GetInstance().Watch(ctx, key)
	for watchResp := range watchChan {
		for _, value := range watchResp.Events {
			fn(value.Kv.Key, value.Kv.Value)
		}
	}
}

// WatchWithPrefix 监听key，带前缀
func WatchWithPrefix(ctx context.Context, key string, fn func(k, v []byte)) {
	watchChan := GetInstance().Watch(ctx, key, clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, value := range watchResp.Events {
			fn(value.Kv.Key, value.Kv.Value)
		}
	}
}

// WatchWith 监听key，带选项
func WatchWith(ctx context.Context, key string, fn func(k, v []byte), opts ...clientv3.OpOption) {
	watchChan := GetInstance().Watch(ctx, key, opts...)
	for watchResp := range watchChan {
		for _, value := range watchResp.Events {
			fn(value.Kv.Key, value.Kv.Value)
		}
	}
}
