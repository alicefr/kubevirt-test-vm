package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"kubevirt.io/client-go/kubecli"
	"kubevirt.io/client-go/log"
)

const (
	programName = "kubevirt-test"
	labelTest   = "test"
)

var (
	vmName string
	pvc    string
)

func usage() string {
	usage := `  # Create a testing VM and SSH access  
  {{ProgramName}} --name vm`
	return usage
}

func NewVirtctlTestCommand() *cobra.Command {
	cobra.AddTemplateFunc(
		"ProgramName", func() string {
			return programName
		},
	)
	var rootCmd = &cobra.Command{
		Use:   programName,
		Short: programName + ": Run and configure KubeVirt VMs for tests",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(cmd.OutOrStderr(), cmd.UsageString())
		},
	}
	cobra.AddTemplateFunc(
		"prepare", func(s string) string {
			// order matters!
			result := strings.Replace(s, "{{ProgramName}}", programName, -1)
			return result
		},
	)
	clientConfig := kubecli.DefaultClientConfig(rootCmd.PersistentFlags())
	rootCmd.AddCommand(
		NewCreateTestVMCommand(clientConfig),
		NewCreateTestPodCommand(clientConfig),
		NewCopyOutputCommand(clientConfig),
		NewDeleteTestCommand(clientConfig),
	)
	return rootCmd
}

func Execute() {
	log.InitializeLogging(programName)
	cmd := NewVirtctlTestCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(cmd.Root().ErrOrStderr(), strings.TrimSpace(err.Error()))
		os.Exit(1)
	}
}
