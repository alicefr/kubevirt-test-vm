package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	k8scorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	podName string
	pvcName string
	image   string
)

type createPodCommand struct {
	clientConfig clientcmd.ClientConfig
}

func NewCreateTestPodCommand(clientConfig clientcmd.ClientConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createPod",
		Short: "create pod for testing",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := createPodCommand{clientConfig: clientConfig}
			return c.run(cmd, args)
		},
	}
	cmd.PersistentFlags().StringVar(&podName, "name", "", "Name for the testing pod")
	cmd.PersistentFlags().StringVar(&pvcName, "pvc", "", "Name of the PVC to run the tests")
	cmd.PersistentFlags().StringVar(&image, "workload-image", DefaultWorkloadImage, "Name of the image to run the tests")
	return cmd
}

func (c *createPodCommand) run(cmd *cobra.Command, args []string) error {
	labels := map[string]string{labelTest: podName}
	client, err := GetKubernetesClient(c.clientConfig)
	if err != nil {
		return err
	}
	namespace, _, err := c.clientConfig.Namespace()
	if err != nil {
		return err
	}
	if podName == "" {
		return fmt.Errorf("the pod cannot be empty")
	}
	if pvcName == "" {
		return fmt.Errorf("the pvc cannot be empty")
	}

	// Create pvc for the output
	pvcOutputName := PvcOutputName(podName)
	err = CreateOutputPVC(client, pvcOutputName, namespace, labels)

	if err != nil {
		return err
	}

	resources := map[k8scorev1.ResourceName]resource.Quantity{
		k8scorev1.ResourceMemory: resource.MustParse("1G"),
		k8scorev1.ResourceCPU:    resource.MustParse("1.0"),
	}
	var volumes []k8scorev1.Volume
	volumes = append(volumes, k8scorev1.Volume{
		Name: pvcOutputName,
		VolumeSource: k8scorev1.VolumeSource{
			PersistentVolumeClaim: &k8scorev1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvcOutputName,
				ReadOnly:  false,
			},
		},
	})
	volumes = append(volumes, k8scorev1.Volume{
		Name: pvcName,
		VolumeSource: k8scorev1.VolumeSource{
			PersistentVolumeClaim: &k8scorev1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvcName,
				ReadOnly:  false,
			},
		},
	})
	privileged := true
	pod := &k8scorev1.Pod{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:   podName,
			Labels: labels,
		},
		Spec: k8scorev1.PodSpec{
			RestartPolicy: k8scorev1.RestartPolicyNever,
			Volumes:       volumes,
			Containers: []k8scorev1.Container{
				{
					Name:       podName,
					Image:      image,
					Command:    []string{"/fio.sh"},
					WorkingDir: OutputDir,
					Stdin:      true,
					TTY:        true,
					Resources: k8scorev1.ResourceRequirements{
						Requests: resources,
					},
					SecurityContext: &k8scorev1.SecurityContext{
						Privileged: &privileged,
					},
					VolumeMounts: []k8scorev1.VolumeMount{
						{
							Name:      pvcOutputName,
							ReadOnly:  false,
							MountPath: OutputDir,
						},
					},
					VolumeDevices: []k8scorev1.VolumeDevice{
						{
							Name:       pvcName,
							DevicePath: "/dev/device-to-test",
						},
					},
				},
			},
			HostNetwork: true,
		},
	}
	_, err = client.CoreV1().Pods(namespace).Create(context.TODO(), pod, k8smetav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
		fmt.Printf("Pod %s already exists \n", podName)
	}
	return nil
}
