package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	kubevirtcorev1 "kubevirt.io/api/core/v1"
	"kubevirt.io/client-go/kubecli"
	"kubevirt.io/kubevirt/pkg/virtctl/templates"
)

type getNodeCommand struct {
	clientConfig clientcmd.ClientConfig
}

func NewGetNodeCommand(clientConfig clientcmd.ClientConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getNode",
		Short: "getNode VM where the VM is running",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getNodeCommand{clientConfig: clientConfig}
			return c.run(cmd, args)
		},
	}
	cmd.PersistentFlags().StringVar(&vmName, "name", "", "Name for the testing VM")
	cmd.SetUsageTemplate(templates.UsageTemplate())
	return cmd
}

func (c *getNodeCommand) run(cmd *cobra.Command, args []string) error {
	if vmName == "" {
		return fmt.Errorf("vm is empty and it si required to be set")
	}

	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(c.clientConfig)
	if err != nil {
		return fmt.Errorf("cannot obtain KubeVirt client: %v", err)
	}
	namespace, _, err := c.clientConfig.Namespace()
	if err != nil {
		return err
	}
	vmi, err := virtClient.VirtualMachineInstance(namespace).Get(vmName, &k8smetav1.GetOptions{})
	if err != nil {
		return err
	}
	if vmi.Status.Phase != kubevirtcorev1.Running {
		return fmt.Errorf("VMI %s isn't in running state: %s", vmName, vmi.Status.Phase)
	}
	node := vmi.ObjectMeta.Labels["kubevirt.io/nodeName"]
	fmt.Println(node)
	return nil
}
