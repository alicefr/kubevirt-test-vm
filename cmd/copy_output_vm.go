package cmd

import (
	"fmt"
	"os"
	"os/exec"
	filepath "path"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"kubevirt.io/kubevirt/pkg/virtctl/templates"
)

var (
	waitForOutput bool
	vmUser        string
	outputDir     string
	key           string
)

const (
	waitCommandOutput = "while [ ! -f /tmp/done ]; do sleep 1; done"
	commandOutput     = "ls /tmp/done"
	virtctlCommand    = "virtctl"
)

type copyOutputVMCommand struct {
	clientConfig clientcmd.ClientConfig
}

func NewCopyOutputVMCommand(clientConfig clientcmd.ClientConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outputVM",
		Short: "copy the results of the tests locally for the VM",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := copyOutputVMCommand{clientConfig: clientConfig}
			return c.run(cmd, args)
		},
	}
	cmd.PersistentFlags().StringVar(&vmName, "name", "vm-test", "Name for the testing VM")
	cmd.PersistentFlags().StringVar(&outputDir, "output", "", "Directory where the test output should be copied")
	cmd.PersistentFlags().BoolVar(&waitForOutput, "wait", false, "Wait until the test finish")
	cmd.PersistentFlags().StringVar(&vmUser, "vm-user", "root", "User for loggin into the VM")
	cmd.PersistentFlags().StringVar(&key, "ssh-key", "", "SSH key to access the VM")
	cmd.SetUsageTemplate(templates.UsageTemplate())
	return cmd
}

func (c *copyOutputVMCommand) run(cmd *cobra.Command, args []string) error {
	commandCheckTest := commandOutput
	if waitForOutput {
		commandCheckTest = waitCommandOutput
	}
	opts := []string{"ssh", "--known-hosts", "", "--command", commandCheckTest}
	if key != "" {
		opts = append(opts, "-i", key)
	}
	opts = append(opts, fmt.Sprintf("%s@%s", vmUser, vmName))
	// Verify that the test has finished
	command := exec.Command("virtctl", opts...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	fmt.Printf("running: %s \n", command.String())
	err := command.Run()
	if err != nil {
		return err
	}

	// Copy output locally
	if outputDir == "" {
		path, err := os.Getwd()
		if err != nil {
			return err
		}
		outputDir = filepath.Join(path, "vm-output")
	}
	// Create output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	opts = []string{"scp",
		"--recursive",
		"--known-hosts", "",
		fmt.Sprintf("%s@%s:%s", vmUser, vmName, OutputTestDir),
		outputDir,
	}
	command = exec.Command("virtctl", opts...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	fmt.Printf("running: %s \n", command.String())
	err = command.Run()
	if err != nil {
		return err
	}
	return nil
}
