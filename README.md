# Run benchmark tests in containers and VMs in Kubernetes

This repository aims to collect best practise for running reproducible tests in containers and [KubeVirt](https://kubevirt.io/) VMs in Kubernetes. Currently, we focus on running storage benchmarks using [fio](https://fio.readthedocs.io/en/latest/index.html) but the same process can be applied for other kind of tests such CPU, memory or network benchmarks. This setup uses containers to delivery the test tools and configurations.

## KubeVirt VMs

In order to run the tests in the VM cloud-init, podman and ssh-server need to be installed and enabled in the VMs. In `containerdisk`, we build a VM image with all the required software that can be used as test OS.

The binary `virtctl-test` helps you in creating the VM with the SSH service in order to be able to execute the commands inside the VM. 

Once the SSH access is enabled, you can use [podman-remote](https://docs.podman.io/en/latest/markdown/podman-remote.1.html) in order to controll the podman in the VM and launch the container. You simply needs to add the connection with:
```bash
podman system connection add test --identity <key> ssh://root@<host>:<ssh-port>/run/user/0/podman/podman.sock
```
