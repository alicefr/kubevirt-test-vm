#!/bin/bash
set -x

LOCAL_REGISTRY=localhost:5001
IMAGES=("fedora-podman-cd:latest", "fio:latest")
CONTAINER_RUNTIME=${CONTAINER_RUNTIME:-docker}
for i in ${IMAGES[@]}; do
	$CONTAINER_RUNTIME tag $i $LOCAL_REGISTRY/$i
	$CONTAINER_RUNTIME push --tls-verify=false $LOCAL_REGISTRY/$i
done
