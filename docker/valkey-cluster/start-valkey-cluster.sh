#!/bin/bash
set -e

PORTS=(8000 8001 8002 8003 8004 8005)

for port in "${PORTS[@]}"; do
  mkdir -p /data/$port
  cat > /data/$port/valkey.conf <<EOF
port $port
cluster-enabled yes
cluster-config-file nodes-$port.conf
cluster-node-timeout 5000
appendonly yes
dir /data/$port/
EOF
  valkey-server /data/$port/valkey.conf --daemonize yes
done

sleep 5

# meet
for port in "${PORTS[@]:1}"; do
  valkey-cli -p ${PORTS[0]} cluster meet 127.0.0.1 $port
done

# addslots (3 master)
valkey-cli -p 8000 cluster addslots $(seq 0 5460)
valkey-cli -p 8001 cluster addslots $(seq 5461 10921)
valkey-cli -p 8002 cluster addslots $(seq 10922 16383)

sleep 10  # 给 gossip 更多时间传播 node ID

# replicate with retry
for slave_port in 8003 8004 8005; do
  master_port=$((slave_port - 3))  # 对应 master: 8000->8003, etc.
  for attempt in {1..5}; do
    NODE_ID=$(valkey-cli -p $master_port cluster myid)
    if valkey-cli -p $slave_port cluster replicate "$NODE_ID"; then
      echo "Replicate success for $slave_port -> $master_port"
      break
    else
      echo "Retry replicate $slave_port ($attempt/5)..."
      sleep 3
    fi
  done
done

# replicate
NODE0=$(valkey-cli -p 8000 cluster myid)
NODE1=$(valkey-cli -p 8001 cluster myid)
NODE2=$(valkey-cli -p 8002 cluster myid)

valkey-cli -p 8003 cluster replicate $NODE0
valkey-cli -p 8004 cluster replicate $NODE1
valkey-cli -p 8005 cluster replicate $NODE2

echo "Valkey Cluster ready"
tail -f /dev/null  # 保持容器运行