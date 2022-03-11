#!/bin/bash -xe

URL=https://download.fedoraproject.org/pub/fedora/linux/releases/35/Cloud/x86_64/images
IMAGE=Fedora-Cloud-Base-35-1.2.x86_64.qcow2
DISK=disk.img
CONTAINER_RUNTIME=${CONTAINER_RUNTIME:-docker}
if ! [ -f "$DISK"  ]; then  
  wget $URL/$IMAGE
  mv $IMAGE $DISK
fi

virt-customize \
	--format qcow2  \
	--run-command "sed -i 's/SELINUX=.*/SELINUX=disabled/g' /etc/selinux/config" \
	-a  $DISK \
	--install podman,openssh-server  \
	--root-password password:test

cat <<EOF > Dockerfile.cd
FROM scratch
COPY disk.img /disk/disk.img
EOF

$CONTAINER_RUNTIME build -t fedora-podman-cd:latest -f Dockerfile.cd .
