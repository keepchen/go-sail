#!/bin/bash
set -e

docker compose -f docker/docker-compose.yml up -d

PACKAGES=$(go list ./... | grep -v '/examples/' | grep -v '/static/' | grep -v '/plugins/' | grep -v '/docker/')
go test $PACKAGES -v -coverprofile=coverage.txt

docker compose -f docker/docker-compose.yml down -v