package cmd

import (
	"context"
	"fmt"
	"strings"

	"io/ioutil"

	"github.com/spf13/cobra"
	k8scorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	kubevirtcorev1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/kubecli"
	"kubevirt.io/kubevirt/pkg/virtctl/templates"
)

const (
	DefaultWorkloadImage = "quay.io/afrosi_rh/fio:latest"
	DefaultVMImage       = "quay.io/afrosi_rh/fedora-podman-cd:latest"
	OutputTestDir        = "/output"
	FioContainerName     = "fio"
	defaultNodePort      = 32756
)

var (
	SSHKeyPath    string
	imageVM       string
	imageWorkload string
	nodePort      int32
	userList      string
)

type createCommand struct {
	clientConfig clientcmd.ClientConfig
}

func NewCreateTestVMCommand(clientConfig clientcmd.ClientConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createVM",
		Short: "create VM for testing",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := createCommand{clientConfig: clientConfig}
			return c.run(cmd, args)
		},
	}
	cmd.PersistentFlags().StringVar(&vmName, "name", "vm-test", "Name for the testing VM")
	cmd.PersistentFlags().StringVar(&SSHKeyPath, "ssh-key", "", "SSH key path to use for accessing the test VM")
	cmd.PersistentFlags().StringVar(&pvc, "pvc", "", "Name of the PVC to run the tests")
	cmd.PersistentFlags().StringVar(&userList, "vm-user", "", "Users to add the ssh key. Specify multiple user separated y ,")
	cmd.PersistentFlags().StringVar(&imageVM, "vm-image", DefaultVMImage, "Name of the image to run the tests")
	cmd.PersistentFlags().StringVar(&imageWorkload, "workload-image", DefaultWorkloadImage, "Name of the image to run the tests")
	cmd.SetUsageTemplate(templates.UsageTemplate())
	return cmd
}

func generatePCIAddress() string {
	return "0000:00:11.0"
}

func pciAddressShell(address string) string {
	return strings.ReplaceAll(address, ":", `\:`)
}

func (c *createCommand) run(cmd *cobra.Command, args []string) error {
	var accessCredential kubevirtcorev1.AccessCredential
	var volumes []kubevirtcorev1.Volume
	var disks []kubevirtcorev1.Disk
	if userList == "" {
		userList = "fedora"
	}
	labels := map[string]string{labelTest: vmName}

	client, err := GetKubernetesClient(c.clientConfig)
	if err != nil {
		return err
	}
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(c.clientConfig)
	if err != nil {
		return fmt.Errorf("cannot obtain KubeVirt client: %v", err)
	}
	namespace, _, err := c.clientConfig.Namespace()
	if err != nil {
		return err
	}
	if SSHKeyPath != "" {
		users := strings.Split(userList, ",")
		// Create a secret out of the ssh key
		data, err := ioutil.ReadFile(SSHKeyPath)
		if err != nil {
			return fmt.Errorf("fail reading ssh file: %v", err)
		}
		secretName := vmName + "-ssh-key"
		secret := &k8scorev1.Secret{
			ObjectMeta: k8smetav1.ObjectMeta{
				Name:   secretName,
				Labels: labels,
			},
			Data: map[string][]byte{"ssh-key": data},
			Type: k8scorev1.SecretTypeOpaque,
		}
		_, err = client.CoreV1().Secrets(namespace).Create(context.TODO(), secret, k8smetav1.CreateOptions{})
		if err != nil {
			if !errors.IsAlreadyExists(err) {
				return err
			}
			fmt.Printf("Secret %s already exists \n", secretName)
		}
		accessCredential = kubevirtcorev1.AccessCredential{
			SSHPublicKey: &kubevirtcorev1.SSHPublicKeyAccessCredential{
				Source: kubevirtcorev1.SSHPublicKeyAccessCredentialSource{
					Secret: &kubevirtcorev1.AccessCredentialSecretSource{
						SecretName: secretName,
					},
				},
				PropagationMethod: kubevirtcorev1.SSHPublicKeyAccessCredentialPropagationMethod{
					QemuGuestAgent: &kubevirtcorev1.QemuGuestAgentSSHPublicKeyAccessCredentialPropagation{
						Users: users,
					},
				},
			},
		}

		if err == nil {
			fmt.Printf("Created secret %s \n", secretName)
		}
	}
	var executeTests string
	if pvc != "" {
		pciAddress := generatePCIAddress()
		volumes = append(volumes, kubevirtcorev1.Volume{
			Name: pvc,
			VolumeSource: kubevirtcorev1.VolumeSource{
				PersistentVolumeClaim: &kubevirtcorev1.PersistentVolumeClaimVolumeSource{
					PersistentVolumeClaimVolumeSource: k8scorev1.PersistentVolumeClaimVolumeSource{
						ClaimName: pvc,
					},
				},
			},
		},
		)
		disks = append(disks, kubevirtcorev1.Disk{
			Name: pvc,
			DiskDevice: kubevirtcorev1.DiskDevice{
				Disk: &kubevirtcorev1.DiskTarget{
					Bus:        "virtio",
					PciAddress: pciAddress,
				},
			},
		})

		executeTests = fmt.Sprintf(`
mkdir -p %s
device=$(ls /sys/bus/pci/devices/%s/virtio*/block/)
[ -z "$device" ] && false
podman run --security-opt label=disable --net=host -d -v %s:/output --name %s --privileged -w /output --tls-verify=false -v /dev/"$device":/dev/device-to-test %s
`, OutputTestDir, pciAddressShell(pciAddress), OutputTestDir, FioContainerName, imageWorkload)
	}
	var order uint
	order = 1
	disks = append(disks, kubevirtcorev1.Disk{
		Name:      "disk0",
		BootOrder: &order,
		DiskDevice: kubevirtcorev1.DiskDevice{
			Disk: &kubevirtcorev1.DiskTarget{
				Bus: "virtio",
			},
		},
	})
	disks = append(disks, kubevirtcorev1.Disk{
		Name: "config-driver",
		DiskDevice: kubevirtcorev1.DiskDevice{
			Disk: &kubevirtcorev1.DiskTarget{
				Bus: "virtio",
			},
		},
	})

	volumes = append(volumes, kubevirtcorev1.Volume{
		Name: "disk0",
		VolumeSource: kubevirtcorev1.VolumeSource{
			ContainerDisk: &kubevirtcorev1.ContainerDiskSource{
				Image: imageVM,
			},
		},
	})
	volumes = append(volumes, kubevirtcorev1.Volume{
		Name: "config-driver",
		VolumeSource: kubevirtcorev1.VolumeSource{
			CloudInitConfigDrive: &kubevirtcorev1.CloudInitConfigDriveSource{
				UserData: fmt.Sprintf(`#!/bin/bash
set -x
sudo systemctl --user enable --now podman.socket
sudo loginctl enable-linger root
 %s `, executeTests),
			},
		},
	})

	requests := map[k8scorev1.ResourceName]resource.Quantity{
		k8scorev1.ResourceMemory: resource.MustParse("1G"),
	}
	// Create VMI
	vmi := &kubevirtcorev1.VirtualMachineInstance{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:   vmName,
			Labels: labels,
		},
		Spec: kubevirtcorev1.VirtualMachineInstanceSpec{
			Domain: kubevirtcorev1.DomainSpec{
				Devices: kubevirtcorev1.Devices{
					Disks: disks,
				},
				Resources: kubevirtcorev1.ResourceRequirements{
					Requests: requests,
				},
			},
			Volumes: volumes,
		},
	}
	if SSHKeyPath != "" {
		vmi.Spec.AccessCredentials = []kubevirtcorev1.AccessCredential{accessCredential}
	}
	vmi, err = virtClient.VirtualMachineInstance(namespace).Create(vmi)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
		fmt.Printf("VMI %s already exists \n", vmName)
	}
	if err == nil {
		fmt.Printf("Created VMI %s \n", vmName)
	}

	return nil
}
