package templates

var DockerfileTpl = `FROM golang:latest AS builder

ARG COMMIT_ID
ARG VERSION=""
ARG VCS_BRANCH=""
ARG GRPC_STUB_REVISION=""
ARG PROJECT_NAME={{ .AppName }}
ARG DOCKER_PROJECT_DIR=/build
ARG EXTRA_BUILD_ARGS=""
ARG GOCACHE=""
ARG GOPROXY=https://goproxy.cn,direct
ARG GOSUMDB=off

WORKDIR $DOCKER_PROJECT_DIR
COPY . $DOCKER_PROJECT_DIR

ENV GOPROXY=$GOPROXY
ENV GOSUMDB=$GOSUMDB

RUN mkdir -p /output \
    && make build -e GOCACHE=$GOCACHE \
    -e COMMIT_ID=$COMMIT_ID -e OUTPUT_FILE=/output/{{ .AppName }} \
    -e VERSION=$VERSION -e VCS_BRANCH=$VCS_BRANCH -e EXTRA_BUILD_ARGS=$EXTRA_BUILD_ARGS

FROM alpine

RUN apk update && apk add tzdata

ENV TZ=Asia/Shanghai
ARG SERVE_MODE
ENV SUB_CMD=SERVE_MODE

COPY --from=builder /output/{{ .AppName }} /usr/bin/{{ .AppName }}
# COPY {{ .AppName }}/pkg/app/[x]/config/config.yaml /data

CMD {{ .AppName }} ${SUB_CMD}
`
