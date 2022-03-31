package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"kubevirt.io/client-go/kubecli"
	"kubevirt.io/kubevirt/pkg/virtctl/templates"
)

type deleteCommand struct {
	clientConfig clientcmd.ClientConfig
}

func NewDeleteTestVMCommand(clientConfig clientcmd.ClientConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete VM and relative service",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := deleteCommand{clientConfig: clientConfig}
			return c.run(cmd, args)
		},
	}
	//cmd.PersistentFlags().StringVar(&vmName, "name", "", "Name for the testing VM")
	cmd.SetUsageTemplate(templates.UsageTemplate())
	return cmd
}

func (c *deleteCommand) run(cmd *cobra.Command, args []string) error {
	if vmName == "" {
		return fmt.Errorf("vm is empty and it si required to be set")
	}

	conf, err := c.clientConfig.ClientConfig()
	if err != nil {
		return err
	}

	client, err := kubernetes.NewForConfig(conf)
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
	err = virtClient.VirtualMachineInstance(namespace).Delete(vmName, &k8smetav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("Failed deleting VMI %s: %v \n", vmName, err)
	}
	err = client.CoreV1().Secrets(namespace).Delete(context.TODO(), vmName+"-ssh-key", k8smetav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("Failed deleting the secret %s: %v \n", vmName+"-ssh-key", err)
	}
	err = client.CoreV1().Services(namespace).Delete(context.TODO(), vmName+"-svc", k8smetav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("Failed deleting service %s: %v \n", vmName+"-svc", err)
	}
	return nil
}
