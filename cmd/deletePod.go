package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"kubevirt.io/kubevirt/pkg/virtctl/templates"
)

type deletePodCommand struct {
	clientConfig clientcmd.ClientConfig
}

func NewDeleteTestPodCommand(clientConfig clientcmd.ClientConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deletePod",
		Short: "clean up the pod and the pvc for the tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := deletePodCommand{clientConfig: clientConfig}
			return c.run(cmd, args)
		},
	}
	cmd.PersistentFlags().StringVar(&podName, "name", "", "Name for the testing pod")
	cmd.SetUsageTemplate(templates.UsageTemplate())
	return cmd
}

func (c *deletePodCommand) run(cmd *cobra.Command, args []string) error {
	if podName == "" {
		return fmt.Errorf("provide a name for the test")
	}

	client, err := GetKubernetesClient(c.clientConfig)
	if err != nil {
		return err
	}
	namespace, _, err := c.clientConfig.Namespace()
	if err != nil {
		return err
	}
	err = client.CoreV1().Pods(namespace).Delete(context.TODO(), podName, k8smetav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	pvcOutput := PvcOutputName(podName)
	err = client.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), pvcOutput, k8smetav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}
