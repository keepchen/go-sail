FROM golang:latest AS builder

ARG COMMIT_ID
ARG VERSION=""
ARG VCS_BRANCH=""
ARG GRPC_STUB_REVISION=""
ARG PROJECT_NAME=go-sail
ARG DOCKER_PROJECT_DIR=/build
ARG EXTRA_BUILD_ARGS=""
ARG GOCACHE=""
ARG GOPROXY=https://goproxy.cn,direct
ARG GOSUMDB=off

WORKDIR $DOCKER_PROJECT_DIR
COPY .. $DOCKER_PROJECT_DIR

ENV GOPROXY=$GOPROXY
ENV GOSUMDB=$GOSUMDB

RUN mkdir -p /output \
    && make build -e GOCACHE=$GOCACHE \
    -e COMMIT_ID=$COMMIT_ID -e OUTPUT_FILE=/output/go-sail \
    -e VERSION=$VERSION -e VCS_BRANCH=$VCS_BRANCH -e EXTRA_BUILD_ARGS=$EXTRA_BUILD_ARGS

FROM alpine:3.23.3

RUN apk update && apk add tzdata

ENV TZ=Asia/Shanghai
ARG SERVE_MODE
ENV SUB_CMD=SERVE_MODE

COPY --from=builder /output/go-sail /usr/bin/go-sail
# COPY go-sail/pkg/app/[x]/config/config.yaml /data

CMD go-sail ${SUB_CMD}