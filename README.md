# Run benchmark tests in containers and VMs in Kubernetes

This repository aims to collect best practise for running reproducible tests in containers and [KubeVirt](https://kubevirt.io/) VMs in Kubernetes. Currently, we focus on running storage benchmarks using [fio](https://fio.readthedocs.io/en/latest/index.html).
## KubeVirt VMs

In order to run the tests in the VM cloud-init, podman and ssh-server need to be installed and enabled in the VMs. In `containerdisk`, we build a VM image with all the required software that can be used as test OS.

The binary `virtctl-test` helps you in creating the VM with the SSH service in order to be able to execute the commands inside the VM. 

## Create a test VM
You can use the command `virtctl-test createVM` in order to create a test VM and the required service and secret for the SSH access. Example:
```bash
$ virtctl-test createVM --name vm-test --pvc pvc-block \
   --ssh-key test-key.pub  --vm-user root \
   --workload-image quay.io/afrosi_rh/fio-demo:latest
```
This will create a secret for the SSH key injction, a Node Port service in order to access the VM and the KubeVirt. The KubeVirt VM is configured to start the fio container using podman at boot time and running the test on the volume specified with the PVC.

## Create test pod 
The command `virtctl createPod` helps you to create a pod to deploy the same test in a container
```bash
$ virtctl-test createPod --image quay.io/afrosi_rh/fio-demo:latest --pod test-fio-pod --pvc pvc-block-pod
```
## Copy the output locally
Copy the output from the cluster to your locally machine varies if you copy from the pod or from the VM
For the VM you can rely on SSH and the `scp` command.
```bash
$ scp -P 32756 -i test-key root@zeus11.lab.eng.tlv2.redhat.com:/tmp/fio* output/vm
```
or you have configured podman remote:
```bash
$ podman -r cp hardcore_turing:/output .
```
For the pod, you need to create an additional pod because the fio-test-pod has already complete.
You can apply a similar yaml and changing the PVC to point to your fio-output pvc:
```yaml
kind: Pod
metadata:
  name: copy-output
spec:
  containers:
  - image: busybox
    name: fio-test
    volumeMounts:
    - name: output
      mountPath: /output
    command: ["tail", "-f", "/dev/null"]
  dnsPolicy: ClusterFirst
  restartPolicy: Never
  volumes:
  - name: output
    persistentVolumeClaim:
      claimName: fio-output-pod-fio-test
```
```bash
$ kubectl apply -f copy-output
$ kubectl cp copy-output:/output .
```
You can find an example of the output you can get in the [example](https://github.com/alicefr/kubevirt-test-vm/tree/main/example) directory

## Create the fio tests
You can rely on the defaults fio test and create them with `make generate-fio-jobs`. Otherwise, you can simply create your own under the directory `tests-setup/fio-jobs/` and build your custom image with `make FIO_IMAGE_NAME=fio-demo fio-image`. The onyl requirements for the test is that the device needs to be `/dev/device-to-test` as the container will bind the device to test as `/dev/device-to-test`
