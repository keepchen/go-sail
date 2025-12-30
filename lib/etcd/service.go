package etcd

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/rand"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// RegisterService 注册服务
func RegisterService(ctx context.Context, serviceName, serviceEndpoint string, ttlSeconds int64) (string, error) {
	lease, err := GetInstance().Grant(ctx, ttlSeconds)
	if err != nil {
		return "", err
	}

	// 生成唯一的服务实例ID
	instanceID := fmt.Sprintf("%s/%s", serviceName, generateInstanceID(serviceEndpoint))

	// 存储服务实例信息
	key := fmt.Sprintf("services/%s/instances/%s", serviceName, instanceID)
	value := serviceEndpoint

	_, err = GetInstance().Put(ctx, key, value, clientv3.WithLease(lease.ID))
	if err != nil {
		return "", err
	}

	// 同时存储服务名称用于服务发现
	serviceKey := fmt.Sprintf("services/%s", serviceName)
	_, err = GetInstance().Put(ctx, serviceKey, "", clientv3.WithLease(lease.ID))

	ch, err := GetInstance().KeepAlive(ctx, lease.ID)
	if err != nil {
		return "", err
	}

	// Keep the lease alive
	go func() {
		for range ch {
		}
	}()

	return instanceID, nil
}

// DiscoverService 发现服务
//
// # Note:
//
// 该方法使用 rand.Intn 随机算法获取可用服务节点
func DiscoverService(ctx context.Context, serviceName string) (string, error) {
	// 使用前缀查询获取所有实例
	prefix := fmt.Sprintf("services/%s/instances/", serviceName)
	resp, err := GetInstance().Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return "", err
	}

	if resp == nil || len(resp.Kvs) == 0 {
		return "", fmt.Errorf("service not found: %s", serviceName)
	}

	// 随机选择一个实例
	idx := rand.Intn(len(resp.Kvs))
	return string(resp.Kvs[idx].Value), nil
}

// GetAllServices 获取所有服务实例
func GetAllServices(ctx context.Context, serviceName string) ([]string, error) {
	prefix := fmt.Sprintf("services/%s/instances/", serviceName)
	resp, err := GetInstance().Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	endpoints := make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		endpoints = append(endpoints, string(kv.Value))
	}

	return endpoints, nil
}

// WatchService 动态感知服务
//
// # Note:
//
// 该方法会阻塞进程
func WatchService(ctx context.Context, serviceName string, fn func(k, v []byte)) {
	Watch(ctx, serviceName, fn)
}

func generateInstanceID(endpoint string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(endpoint)))
}
