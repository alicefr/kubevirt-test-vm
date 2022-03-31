package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"kubevirt.io/client-go/kubecli"
	"kubevirt.io/kubevirt/pkg/virtctl/templates"
)

var volume string

type getDiskCommand struct {
	clientConfig clientcmd.ClientConfig
}

func NewGetDiskCommand(clientConfig clientcmd.ClientConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getDisk",
		Short: "getDisk VM and relative service",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := getDiskCommand{clientConfig: clientConfig}
			return c.run(cmd, args)
		},
	}
	cmd.PersistentFlags().StringVar(&volume, "volume", "", "Name of the volume to run the tests")
	cmd.SetUsageTemplate(templates.UsageTemplate())
	return cmd
}

func (c *getDiskCommand) run(cmd *cobra.Command, args []string) error {
	if vmName == "" {
		return fmt.Errorf("vm is empty and it si required to be set")
	}
	if volume == "" {
		return fmt.Errorf("pvc is empty and it si required to be set")
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

	for _, vs := range vmi.Status.VolumeStatus {
		if vs.Name == volume {
			fmt.Println(vs.Target)
			break
		}

	}
	return nil
}
