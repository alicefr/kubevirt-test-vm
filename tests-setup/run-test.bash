#!/bin/bash -x

DEVICE=$1
REMOTE=${REMOTE:-}
CONTAINER_RUNTIME=${CONTAINER_RUNTIME:-docker}
OUTPUT_DIR=${OUTPUT_DIR:-/tmp/output-fio}
remote=""
mkdir -p ${OUTPUT_DIR}
rm -rf ${OUTPUT_DIR}/*

if [ -z "$DEVICE" ]; then
	echo "Provide the device to test as first argument"
	exit 0
fi 
if [ $REMOTE -eq 1 ] & [ "$CONTAINER_RUNTIME" == "podman" ];  then 
	remote="-r"
	flags="--tls-verify=false"
fi
${CONTAINER_RUNTIME} $remote run --security-opt label=disable -ti \
	-v ${OUTPUT_DIR}:/output \
	$flags \
	--privileged \
	-w /output \
	-v ${DEVICE}:/dev/device-to-test \
	localhost:5000/fio:latest \
	"fio /fio-jobs/*.fio --output /output/fio.log"
