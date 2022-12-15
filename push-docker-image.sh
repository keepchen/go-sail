#!/bin/sh

service_name=go-sail

pull_images(){
  docker info
  docker login -u "${local_harbor_username}" -p "${local_harbor_password}" "${local_harbor_host}"
  echo "|--[pull local] current service_name name: ${service_name}"
  docker pull "${local_harbor_host}"/keepchen/"${service_name}":"${SHORT_COMMIT_HASH}"
  docker logout
}

push_images_prod(){
  docker login -u "${prod_registry_username}" -p "${prod_registry_password}" "${prod_registry_host}"
  echo "|--[push prod] current service_name name: ${service_name}"
  docker tag "${local_harbor_host}"/keepchen/"${service_name}":"${SHORT_COMMIT_HASH}" "${prod_registry_host}"/keepchen/"${service_name}":"${SHORT_COMMIT_HASH}"
  docker push "${prod_registry_host}"/keepchen/"${service_name}":"${SHORT_COMMIT_HASH}"
  docker logout
}

case "${DRONE_PROMOTE_TARGET}" in
  prod|production)
    pull_images "$@"
    push_images_prod "$@"
    ;;
  pre|preview)
    pull_images "$@"
    push_images_preview "$@"
    ;;
  test)
    pull_images "$@"
    push_images_test "$@"
    ;;
  *)
    echo "[!] Not match any target( ${DRONE_PROMOTE_TARGET} ),short commit hash is: [ ${SHORT_COMMIT_HASH} ],task quit."
    exit 0
    ;;
esac

# ==== clear images ====
dirty_images=$(docker images | grep "${SHORT_COMMIT_HASH}" | awk '{print $3}')
echo "|--Tips: dirty images: ${dirty_images}"
if [ -n "${dirty_images}" ]; then
    echo "clear dirty images..."
    docker images | grep "${SHORT_COMMIT_HASH}" | awk '{print $1 ":" $2}' | xargs docker rmi
fi
echo "Task completed!,image tag is:[ ${SHORT_COMMIT_HASH} ]. Enjoy it. :)"
