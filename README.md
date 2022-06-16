# Run benchmark tests in containers and VMs in Kubernetes

This repository aims to collect best practise for running reproducible tests in containers and [KubeVirt](https://kubevirt.io/) VMs in Kubernetes. Currently, it focuses on running storage benchmarks using [fio](https://fio.readthedocs.io/en/latest/index.html).

## KubeVirt VMs
For running the tests in the VM, cloud-init, podman, qemu-guest-agent and ssh-server need to be installed and enabled in the VM. We build [containerdisk](https://github.com/alicefr/kubevirt-test-vm/tree/main/containerdisk) with all the required software needed to setup the VM.

The binary `virtctl-test` helps you in creating the VM with the SSH service in order to be able to execute the commands inside the VM. You need to have install on your system `kubectl` and `virtctl` binaries as `virtctl-test` call them in order to copy the output from the VMs and the pods.

## Create a test VM
You can use the command `virtctl-test createVM` in order to create a test VM and the required service and secret for the SSH access. Example:
```bash
$ virtctl-test createVM --name vm-test --pvc pvc \
   --ssh-key test-key.pub  --vm-user fedora
```
This will create a secret for the SSH key injction and a running VMI. The KubeVirt VM is configured to start the fio container using podman at boot time and running the test on the volume specified with the PVC.

## Create test pod 
The command `virtctl createPod` helps you to create a pod to deploy the same test in a container
```bash
$ virtctl-test createPod --pod test-pod \
  --pvc pvc
```
Running fio tests inside with standard containers can be compared to the fio test run in the VM, if the PVC is a block.

## Copy the output locally
Copy the test output for the VM:
```bash
$ virtctl-test outputVM --name vm-test \
  --ssh-key test-key.pub  --vm-user fedora
```
If the test hasn't finished yet, the command will fail. If you want to block and wait for the test to finish you can add the flag `--wait`.

Copy the test output for the pod:
```bash
$ virtctl-test outputPod --name fio-test
```

## Create the fio tests
You can rely on the defaults fio test and create them with `make generate-fio-jobs`. Otherwise, you can simply create your own tests under the directory `tests-setup/fio-jobs/` and build your custom image with `make FIO_IMAGE_NAME=fio-demo fio-image`. The onyl requirements for the test is that the device needs to be `/dev/device-to-test` as the container will bind the device to test as `/dev/device-to-test`
