package etcd

import (
	"context"
	"fmt"
	"math/rand"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// RegisterService 注册服务
func RegisterService(ctx context.Context, serviceName, serviceEndpoint string, ttlSeconds int64) error {
	lease, err := GetInstance().Grant(ctx, ttlSeconds)
	if err != nil {
		return err
	}

	_, err = GetInstance().Put(ctx, serviceName, serviceEndpoint, clientv3.WithLease(lease.ID))

	if err != nil {
		return err
	}

	ch, err := GetInstance().KeepAlive(ctx, lease.ID)
	if err != nil {
		return err
	}

	// Keep the lease alive
	go func() {
		for {
			<-ch
		}
	}()

	return nil
}

// DiscoverService 发现服务
//
// # Note:
//
// 该方法使用 rand.Intn 随机算法获取可用服务节点
func DiscoverService(ctx context.Context, serviceName string) (string, error) {
	resp, err := GetInstance().Get(ctx, serviceName)
	if err != nil {
		return "", err
	}

	if resp == nil || len(resp.Kvs) == 0 {
		return "", fmt.Errorf("service not found")
	}

	return string(resp.Kvs[rand.Intn(len(resp.Kvs))].Value), nil
}

// WatchService 动态感知服务
//
// # Note:
//
// 该方法会阻塞进程
func WatchService(ctx context.Context, serviceName string, fn func(k, v []byte)) {
	Watch(ctx, serviceName, fn)
}
