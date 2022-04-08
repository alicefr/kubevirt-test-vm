CONTAINER_RUNTIME ?= podman
REGISTRY ?= quay.io/afrosi_rh
DEVICE_IN_CONTAINER ?= device-to-test
TIME_RUNNING_TEST ?= 300
TEST_FLAVOR ?= write,read,randread,randwrite
BLOCKSIZE ?= 4k,1m
FIO_JOBS_DIR=tests-setup/fio-jobs
TAG=latest
CD_IMAGE=$(REGISTRY)/fedora-podman-cd:$(TAG)
FIO_IMAGE=$(REGISTRY)/fio:$(TAG)

build:
	mkdir -p bin
	go build -o bin/virtctl-test main.go

fio-image:
	$(CONTAINER_RUNTIME) build -t $(FIO_IMAGE) tests-setup

cd-image:
	$(CONTAINER_RUNTIME) build -t $(CD_IMAGE) containerdisk

generate-fio-jobs:
	mkdir -p fio-jobs
	$(CONTAINER_RUNTIME) run -ti --security-opt label=disable \
	-v $(PWD)/$(FIO_JOBS_DIR):/fio-jobs \
	-w /fio-jobs \
	--hostname fio \
	--entrypoint genfio \
	fio \
	-d /dev/$(DEVICE_IN_CONTAINER) -r $(TIME_RUNNING_TEST) -m $(TEST_FLAVOR) -b $(BLOCKSIZE) -s -x fio

clean:
	rm -rf $(OUTPUT_DIR) $(FIO_JOBS_DIR)
