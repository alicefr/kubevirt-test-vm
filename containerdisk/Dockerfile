FROM quay.io/kubevirt/libguestfs-tools:20220331_9da70587e as builder

ENV LIBGUESTFS_BACKEND direct
ENV LIBGUESTFS_PATH /usr/local/lib/guestfs/appliance
USER 0
RUN /entrypoint.sh
RUN virt-builder fedora-35 \
        -v \
        --format qcow2  \
        --run-command "sed -i 's/SELINUX=.*/SELINUX=disabled/g' /etc/selinux/config" \
        -o /tmp/disk.img \
        --install podman,openssh-server  \
        --root-password password:test 
FROM stratch

COPY --from=builder /tmp/disk.img /disk/disk.img
