#!/bin/bash

set -x

device=/dev/device-to-test
source_device=$(findmnt -n -o SOURCE $device| sed 's/.*\[\/\([^]]*\)\].*/\1/g')
OUTPUT_DIR=/output
lsblk | grep "$source_device"
if [ $? -ne 0 ]; then
    echo "Device $device not found"
    exit 1
fi

mkdir -p $OUTPUT_DIR
# Preallocate disks to have consistent tests
size=$(lsblk -n -o SIZE /dev/$source_device)
if [ -z "$size" ]; then
	echo "Failed parising the size of the device"
	exit 1
fi
fallocate -l $size -x $device

fio /fio-jobs/*.fio --output $OUTPUT_DIR/fio.log
