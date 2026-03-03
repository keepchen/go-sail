FROM valkey/valkey:9-alpine AS base

RUN export https_proxy=http://192.168.160.131:7897 http_proxy=http://192.168.160.131:7897 all_proxy=socks5://192.168.160.131:7897

# 安装 redis-cli 等工具（valkey-cli 同 redis-cli）
RUN apk add --no-cache bash curl

# 复制你的 cluster 启动脚本（从你原来的 setup-valkey-cluster.sh 改编）
COPY start-valkey-cluster.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/start-valkey-cluster.sh

# 暴露所有 cluster 端口
EXPOSE 8000-8005 18000-18005

CMD ["start-valkey-cluster.sh"]