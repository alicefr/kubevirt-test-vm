package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"kubevirt.io/client-go/kubecli"
	"kubevirt.io/kubevirt/pkg/virtctl/templates"
)

type deleteCommand struct {
	clientConfig clientcmd.ClientConfig
}

func NewDeleteTestCommand(clientConfig clientcmd.ClientConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "clean up the pod/VM and the pvc for the tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := deleteCommand{clientConfig: clientConfig}
			return c.run(cmd, args)
		},
	}
	cmd.PersistentFlags().StringVar(&name, "name", "", "Name for the testing pod/VM")
	cmd.SetUsageTemplate(templates.UsageTemplate())
	return cmd
}

func (c *deleteCommand) run(cmd *cobra.Command, args []string) error {
	if name == "" {
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
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(c.clientConfig)
	if err != nil {
		return fmt.Errorf("cannot obtain KubeVirt client: %v", err)
	}

	// Delete the pod if it exists
	err = client.CoreV1().Pods(namespace).Delete(context.TODO(), name, k8smetav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if err == nil {
		fmt.Printf("Deleted pod %s \n", name)
	}

	// Delete the VM if it exists
	err = virtClient.VirtualMachineInstance(namespace).Delete(name, &k8smetav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if err == nil {
		fmt.Printf("Deleted VM %s \n", name)
	}

	// Delete the output PVC if exists
	pvcOutput := PvcOutputName(name)
	err = client.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), pvcOutput, k8smetav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if err == nil {
		fmt.Printf("Deleted PVC %s \n", pvcOutput)
	}

	// Delete the SSH secret if it exists
	err = client.CoreV1().Secrets(namespace).Delete(context.TODO(), name+"-ssh-key", k8smetav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err

	}
	if err == nil {
		fmt.Printf("Deleted secret %s \n", name+"-ssh-key")
	}

	return nil
}
