#!/bin/bash

set -x

device=/dev/device-to-test
OUTPUT_DIR=/output
lsblk | grep "$device"
if [ $? -ne 0 ]; then
    echo "Device $device not found"
    exit 1
fi

mkdir -p $OUTPUT_DIR
# Preallocate disks to have consistent tests
size=$(lsblk -n -o SIZE $device)
if [ -z "$size" ]; then
	echo "Failed parising the size of the device"
	exit 1
fi
fallocate -l $size -x $device

fio /fio-jobs/*.fio --output $OUTPUT_DIR/fio.log
