#!/bin/bash

set -x

device=device-to-test
OUTPUT_DIR=/output
if [  $(lsblk | grep "$device") -ne 0 ]; then 
    echo "Device $device not found"
    exit 1
fi

mkdir -p $OUTPUT_DIR
# Preallocate disks to have consistent tests
size=$(lsblk |grep $device |awk '{ print $4 }')
fallocate -l $size -x $device

fio /fio-jobs/*.fio --output $OUTPUT_DIR/fio.log
