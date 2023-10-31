#!/bin/bash

set -e

env

push_tag="${REGISTRY_HOST}/keepchen/go-sail:${SHORT_COMMIT_ID}"

echo "${push_tag}"

docker build --tag ${push_tag} \
  --build-arg EXTRA_BUILD_ARGS=-mod=vendor \
  --build-arg COMMIT_ID=${CI_COMMIT_SHA} \
  --build-arg COMMIT_TAG=${SHORT_COMMIT_ID} \
  --build-arg VCS_BRANCH=${CI_COMMIT_BRANCH} \
  --build-arg VERSION=${SHORT_COMMIT_ID} .

docker image ls

docker push ${push_tag}

dirty_images=$(docker images | grep "${SHORT_COMMIT_ID}" | awk '{print $3}')
echo "|--Tips: dirty images: ${dirty_images}"
if [ -n "${dirty_images}" ]; then
    echo "clear dirty images..."
    docker images | grep "${SHORT_COMMIT_ID}" | awk '{print $1 ":" $2}' | xargs docker rmi
fi
docker logout

echo "Task completed!,image tag is:[ $(echo ${SHORT_COMMIT_ID}|cut -c1-8) ]. Enjoy it. :)"