#!/bin/bash -x

DEVICE=$1
OUTPUT_DIR=${OUTPUT_DIR:-/tmp/}
if [ -z "$DEVICE" ]; then
	echo "Provide the device to test as first argument"
	exit 0
fi 

podman -r run --security-opt label=disable -d \
	-v ${OUTPUT_DIR}:/output \
	$flags \
	--privileged \
	-w /output \
	--tls-verify=false \
	-v ${DEVICE}:/dev/device-to-test \
	quay.io/afrosi_rh/fio:latest
