FROM alpine:3.18

# 安装编译依赖 + 运行时工具
RUN apk add --no-cache \
    build-base \
    wget \
    tar \
    bash \
    && wget https://download.redis.io/releases/redis-6.2.6.tar.gz \
    && tar xzf redis-6.2.6.tar.gz \
    && cd redis-6.2.6 \
    && make -j$(nproc) \
    && make install \
    && cd .. \
    && rm -rf redis-6.2.6 redis-6.2.6.tar.gz \
    && apk del build-base wget tar

# 创建根目录（启动脚本会填充子目录和 conf）
RUN mkdir -p /redis-cluster

# 只 COPY 一个启动脚本
COPY start-redis-cluster.sh /usr/local/bin/start-redis-cluster.sh
RUN chmod +x /usr/local/bin/start-redis-cluster.sh

EXPOSE 7000-7005 17000-17005

CMD ["/usr/local/bin/start-redis-cluster.sh"]