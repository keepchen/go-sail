#!/bin/bash

# 生成 Valkey 配置
for port in 8000 8001 8002 8003 8004 8005; do
  mkdir -p valkey-cluster/$port
  mkdir -p valkey-cluster/$port/data
  cat > valkey-cluster/$port/valkey.conf << EOF
port $port
cluster-enabled yes
cluster-config-file nodes.conf
cluster-node-timeout 5000
appendonly yes
bind 0.0.0.0
protected-mode no
daemonize yes
pidfile valkey-cluster/$port/valkey.pid
dir valkey-cluster/$port/data
cluster-announce-ip 127.0.0.1
cluster-announce-port $port
cluster-announce-bus-port $((port + 10000))
EOF
done
