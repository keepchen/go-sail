package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Watch watches a key
func Watch(ctx context.Context, key string, fn func(k, v []byte)) {
	watchChan := GetInstance().Watch(ctx, key)
	for watchResp := range watchChan {
		for _, value := range watchResp.Events {
			fn(value.Kv.Key, value.Kv.Value)
		}
	}
}

// WatchWithPrefix watches keys with a prefix
func WatchWithPrefix(ctx context.Context, key string, fn func(k, v []byte)) {
	watchChan := GetInstance().Watch(ctx, key, clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, value := range watchResp.Events {
			fn(value.Kv.Key, value.Kv.Value)
		}
	}
}

// WatchWith watches a key with options
func WatchWith(ctx context.Context, key string, fn func(k, v []byte), opts ...clientv3.OpOption) {
	watchChan := GetInstance().Watch(ctx, key, opts...)
	for watchResp := range watchChan {
		for _, value := range watchResp.Events {
			fn(value.Kv.Key, value.Kv.Value)
		}
	}
}
