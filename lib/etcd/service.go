package etcd

import (
	"context"
	"crypto/md5"
	"fmt"
	"math/rand"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// RegisterService registers a service
func RegisterService(ctx context.Context, groupName, serviceName, serviceEndpoint string, ttlSeconds int64) (string, error) {
	lease, err := GetInstance().Grant(ctx, ttlSeconds)
	if err != nil {
		return "", err
	}

	// Generate a unique service instance ID
	instanceID := fmt.Sprintf("%s/%s/%s", groupName, serviceName, generateInstanceID(serviceEndpoint))

	// Store service instance information
	key := fmt.Sprintf("services/%s/%s/instances/%s", groupName, serviceName, instanceID)
	value := serviceEndpoint

	_, err = GetInstance().Put(ctx, key, value, clientv3.WithLease(lease.ID))
	if err != nil {
		return "", err
	}

	// Also store the service name for service discovery
	serviceKey := fmt.Sprintf("services/%s/%s", groupName, serviceName)
	_, err = GetInstance().Put(ctx, serviceKey, "", clientv3.WithLease(lease.ID))
	if err != nil {
		return "", err
	}

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

// DiscoverService discovers a service
//
// # Note:
//
// This method uses rand.Intn to randomly select an available service node
func DiscoverService(ctx context.Context, groupName, serviceName string) (string, error) {
	// Use a prefix query to retrieve all instances
	prefix := fmt.Sprintf("services/%s/%s/instances/", groupName, serviceName)
	resp, err := GetInstance().Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return "", err
	}

	if resp == nil || len(resp.Kvs) == 0 {
		return "", fmt.Errorf("service not found: %s", serviceName)
	}

	// Randomly select an instance
	idx := rand.Intn(len(resp.Kvs))
	return string(resp.Kvs[idx].Value), nil
}

// GetAllServices retrieves all service instances
func GetAllServices(ctx context.Context, groupName, serviceName string) ([]string, error) {
	prefix := fmt.Sprintf("services/%s/%s/instances/", groupName, serviceName)
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

// WatchService watches for service changes
//
// # Note:
//
// This method blocks execution
func WatchService(ctx context.Context, groupName, serviceName string, fn func(k, v []byte)) {
	key := fmt.Sprintf("services/%s/%s/instances/", groupName, serviceName)
	WatchWithPrefix(ctx, key, fn)
}

func generateInstanceID(endpoint string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(endpoint)))
}
