#!/bin/bash
set -e

echo "Generating Redis cluster configuration..."

for port in 7000 7001 7002 7003 7004 7005; do
  # 先在 bash 里计算 bus port（避免 conf 里有表达式）
  bus_port=$((port + 10000))

  mkdir -p /redis-cluster/$port

  cat > /redis-cluster/$port/redis.conf <<EOF
port $port
bind 0.0.0.0
protected-mode no
cluster-enabled yes
cluster-config-file nodes.conf
cluster-node-timeout 5000
cluster-announce-ip 127.0.0.1
cluster-announce-port $port
cluster-announce-bus-port $bus_port
appendonly yes
dir /redis-cluster/$port/
logfile /redis-cluster/$port/redis.log
daemonize no
pidfile /redis-cluster/$port/redis.pid
EOF
done

echo "Configuration files generated."

echo "Starting Redis instances..."

for port in 7000 7001 7002 7003 7004 7005; do
  redis-server /redis-cluster/$port/redis.conf &
done

sleep 12   # 多等2秒，确保实例完全启动（尤其是首次）

echo "All instances started."

echo "Initializing cluster..."

# Meet 节点
for port in 7001 7002 7003 7004 7005; do
  redis-cli -p $port cluster meet 127.0.0.1 7000 || echo "Meet $port failed, retrying later if needed"
done

sleep 8   # gossip 传播需要时间

# 分配槽位
redis-cli -p 7000 cluster addslots $(seq 0 5460) || true
redis-cli -p 7001 cluster addslots $(seq 5461 10921) || true
redis-cli -p 7002 cluster addslots $(seq 10922 16383) || true

sleep 5

# 设置 slave + 重试
echo "Setting up replicas..."
for slave_port in 7003 7004 7005; do
  master_port=$((slave_port - 3))
  for attempt in {1..8}; do   # 增加重试次数和间隔
    NODE_ID=$(redis-cli -p $master_port cluster myid 2>/dev/null)
    if [ -n "$NODE_ID" ] && redis-cli -p $slave_port cluster replicate "$NODE_ID"; then
      echo "Replicate success: $slave_port -> $master_port ($NODE_ID)"
      break
    else
      echo "Retry replicate $slave_port ($attempt/8)..."
      sleep 4
    fi
  done
done

# 最终检查（可选，但推荐）
echo "Cluster status:"
redis-cli -p 7000 cluster info
redis-cli -p 7000 cluster nodes

echo "Redis Cluster setup complete"

wait