#!/bin/bash
set -x

LOCAL_REGISTRY=localhost:5000
IMAGES=("fedora-podman-cd:latest" "fio:latest")
CONTAINER_RUNTIME=${CONTAINER_RUNTIME:-docker}
if [ $CONTAINER_RUNTIME == podman  ]; then
	CONTAINER_RUNTIME_FLAGS=--tls-verify=false
fi

for i in ${IMAGES[@]}; do
	$CONTAINER_RUNTIME tag $i $LOCAL_REGISTRY/$i
	$CONTAINER_RUNTIME push $CONTAINER_RUNTIME_FLAGS  $LOCAL_REGISTRY/$i
done
