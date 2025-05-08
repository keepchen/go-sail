#!/bin/bash

# 生成 Redis 配置
for port in 7000 7001 7002 7003 7004 7005; do
  mkdir -p redis-cluster/$port
  mkdir -p redis-cluster/$port/data
  cat > redis-cluster/$port/redis.conf << EOF
port $port
cluster-enabled yes
cluster-config-file nodes.conf
cluster-node-timeout 5000
appendonly yes
bind 0.0.0.0
protected-mode no
daemonize yes
pidfile redis-cluster/$port/redis.pid
dir redis-cluster/$port/data
EOF
done
