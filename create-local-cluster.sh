#!/bin/bash -x

make cd-image
make fio-image 
cluster/kind-with-registry.sh
cluster/push-image-local-registry.sh
cluster/install-kubevirt-latest.bash
