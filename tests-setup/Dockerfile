FROM fedora:35

RUN dnf install -y fio python3 util-linux && dnf clean all -y
RUN mkdir -p /output
COPY ./fio-jobs /fio-jobs
COPY ./fio.sh /fio.sh
ENTRYPOINT ["/fio.sh"]
