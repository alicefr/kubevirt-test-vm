package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	filepath "path"
	"time"

	"github.com/spf13/cobra"
	k8scorev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/client/conditions"
	"kubevirt.io/kubevirt/pkg/virtctl/templates"
)

const copyOutputImage = "quay.io/quay/busybox:latest"

type copyOutputPodCommand struct {
	clientConfig clientcmd.ClientConfig
}

func NewCopyOutputPodCommand(clientConfig clientcmd.ClientConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outputPod",
		Short: "copy the results of the tests locally for the VM",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := copyOutputPodCommand{clientConfig: clientConfig}
			return c.run(cmd, args)
		},
	}
	cmd.PersistentFlags().StringVar(&podName, "name", "", "Name for the testing pod")
	cmd.PersistentFlags().StringVar(&outputDir, "output", "", "Directory where the test output should be copied")
	cmd.SetUsageTemplate(templates.UsageTemplate())
	return cmd
}

func (c *copyOutputPodCommand) run(cmd *cobra.Command, args []string) error {
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
	copyPodName := "copy-" + podName
	// Copy output locally
	if outputDir == "" {
		path, err := os.Getwd()
		if err != nil {
			return err
		}
		outputDir = filepath.Join(path, "pod-output")
	}
	pvcOutput := PvcOutputName(podName)
	dir := "/output"
	pod := &k8scorev1.Pod{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name: copyPodName,
		},
		Spec: k8scorev1.PodSpec{
			RestartPolicy: k8scorev1.RestartPolicyNever,
			Volumes: []k8scorev1.Volume{
				{
					Name: pvcOutput,
					VolumeSource: k8scorev1.VolumeSource{
						PersistentVolumeClaim: &k8scorev1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvcOutput,
							ReadOnly:  false,
						},
					},
				},
			},
			Containers: []k8scorev1.Container{
				{
					Name:    "copy",
					Image:   image,
					Command: []string{"tail", "-f", "/dev/null"},
					VolumeMounts: []k8scorev1.VolumeMount{
						{
							Name:      pvcOutput,
							ReadOnly:  false,
							MountPath: dir,
						},
					},
				},
			},
		},
	}
	_, err = client.CoreV1().Pods(namespace).Create(context.TODO(), pod, k8smetav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
		fmt.Printf("Pod %s already exists \n", podName)
	}
	if err = waitPodRunning(client, copyPodName, namespace); err != nil {
		return err
	}

	opts := []string{"cp",
		fmt.Sprintf("%s/%s:%s", namespace, copyPodName, dir),
		outputDir,
	}
	command := exec.Command("kubectl", opts...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	fmt.Printf("running: %s \n", command.String())
	err = command.Run()
	if err != nil {
		return err
	}
	err = client.CoreV1().Pods(namespace).Delete(context.TODO(), copyPodName, k8smetav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

func isPodRunning(client *kubernetes.Clientset, p, namespace string) wait.ConditionFunc {
	return func() (bool, error) {

		pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), p, k8smetav1.GetOptions{})
		if err != nil {
			return false, err
		}

		switch pod.Status.Phase {
		case v1.PodRunning:
			return true, nil
		case v1.PodFailed, v1.PodSucceeded:
			return false, conditions.ErrPodCompleted
		}
		return false, nil
	}
}

func waitPodRunning(client *kubernetes.Clientset, pod, ns string) error {
	return wait.PollImmediate(2*time.Second, 120*time.Second, isPodRunning(client, pod, ns))
}
