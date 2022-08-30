package cmd_test

import (
	"fmt"
	"os"
	filepath "path"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/testing"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"kubevirt.io/client-go/kubecli"
	cmd "kubevirt.io/test-benchmarks/cmd"
)

func mockGetKubernetesClient(_ clientcmd.ClientConfig) (kubernetes.Interface, error) {
	return fake.NewSimpleClientset(), nil
}

func mockGetKubernetesClientReadyPod(_ clientcmd.ClientConfig) (kubernetes.Interface, error) {
	kubeClient := fake.NewSimpleClientset()
	kubeClient.Fake.PrependReactor("get", "pods", func(action testing.Action) (bool, runtime.Object, error) {
		podRunning := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "copy-test",
			},
			Status: v1.PodStatus{
				Phase: v1.PodRunning,
			},
		}
		return true, podRunning, nil
	})

	return kubeClient, nil
}

type mockClientConfig struct{}

func (m *mockClientConfig) RawConfig() (clientcmdapi.Config, error) {
	panic("not implemented")
}

func (m *mockClientConfig) ClientConfig() (*restclient.Config, error) {
	panic("not implemented")
}

func (m *mockClientConfig) Namespace() (string, bool, error) {
	return "test", false, nil
}

func (m *mockClientConfig) ConfigAccess() clientcmd.ConfigAccess {
	panic("not implemented")
}

var _ = Describe("Create", func() {
	var (
		config             clientcmd.ClientConfig
		ctrl               *gomock.Controller
		mockKubevirtClient *kubecli.MockKubevirtClient
		mockVMI            *kubecli.MockVirtualMachineInstanceInterface
	)
	mockGetKubevirtClientFromClientConfig := func(cmdConfig clientcmd.ClientConfig) (kubecli.KubevirtClient, error) {
		return mockKubevirtClient, nil
	}
	goMockVMCreation := func(err error) {
		mockKubevirtClient.EXPECT().VirtualMachineInstance(gomock.Any()).Return(mockVMI)
		mockVMI.EXPECT().Create(gomock.Any()).Return(nil, err)
	}

	Context("Create benchmark pod", func() {
		BeforeEach(func() {
			config = &mockClientConfig{}
			cmd.SetKubernetesClientFunc(mockGetKubernetesClient)
		})
		It("should successfully create the pod", func() {
			cmd := cmd.NewCreateTestPodCommand(config)
			cmd.SetArgs([]string{"--name", "test-pod", "--pvc", "pvc"})
			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("should successfully create the pod with custom image", func() {
			cmd := cmd.NewCreateTestPodCommand(config)
			cmd.SetArgs([]string{"--name", "test-pod",
				"--pvc", "pvc", "--workload-image", "test-image"})
			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("shoud fail when the name isn't specified", func() {
			cmd := cmd.NewCreateTestPodCommand(config)
			cmd.SetArgs([]string{"--pvc", "pvc"})
			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
		})
		It("shoud fail when the pvc isn't specified", func() {
			cmd := cmd.NewCreateTestPodCommand(config)
			cmd.SetArgs([]string{"--name", "test-pod"})
			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
		})
	})
	Context("Create benchmark vm", func() {
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			config = &mockClientConfig{}
			cmd.SetKubernetesClientFunc(mockGetKubernetesClient)
			mockKubevirtClient = kubecli.NewMockKubevirtClient(ctrl)
			mockVMI = kubecli.NewMockVirtualMachineInstanceInterface(ctrl)
			kubecli.GetKubevirtClientFromClientConfig = mockGetKubevirtClientFromClientConfig
		})
		It("should successfully create the VM", func() {
			goMockVMCreation(nil)
			cmd := cmd.NewCreateTestVMCommand(config)
			cmd.SetArgs([]string{"--name", "test-vm", "--pvc", "pvc"})
			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("should be successful with the option vm-image", func() {
			goMockVMCreation(nil)
			cmd := cmd.NewCreateTestVMCommand(config)
			cmd.SetArgs([]string{"--name", "test-vm", "--pvc", "pvc", "--vm-image", "test-image"})
			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("should be successful with the option workload-image", func() {
			goMockVMCreation(nil)
			cmd := cmd.NewCreateTestVMCommand(config)
			cmd.SetArgs([]string{"--name", "test-vm", "--pvc", "pvc", "--workload-image", "test-image"})
			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("should failing create the VM", func() {
			goMockVMCreation(fmt.Errorf("test erro"))
			cmd := cmd.NewCreateTestVMCommand(config)
			cmd.SetArgs([]string{"--name", "test-vm", "--pvc", "pvc"})
			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
		})
		It("should be successful if the VM already exists", func() {
			goMockVMCreation(errors.NewAlreadyExists(schema.GroupResource{}, ""))
			cmd := cmd.NewCreateTestVMCommand(config)
			cmd.SetArgs([]string{"--name", "test-vm", "--pvc", "pvc"})
			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("should failing when the pvc isn't specified", func() {
			cmd := cmd.NewCreateTestVMCommand(config)
			cmd.SetArgs([]string{"--name", "test-vm"})
			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
		})
		It("shoud fail when the name isn't specified", func() {
			cmd := cmd.NewCreateTestVMCommand(config)
			cmd.SetArgs([]string{"--pvc", "pvc"})
			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
		})
		It("should successfully create the VM with the SSH option", func() {
			dir := GinkgoT().TempDir()
			file, err := os.Create(filepath.Join(dir, "test.key"))
			Expect(err).NotTo(HaveOccurred())
			_, err = file.WriteString("ssh-rsa ajdsjadk")
			Expect(err).NotTo(HaveOccurred())

			goMockVMCreation(nil)
			cmd := cmd.NewCreateTestVMCommand(config)
			cmd.SetArgs([]string{"--name", "test-vm", "--pvc", "pvc"})
			err = cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("should fail because ssh file doesn't exis", func() {
			cmd := cmd.NewCreateTestVMCommand(config)
			cmd.SetArgs([]string{"--name", "test-vm", "--pvc", "pvc", "--ssh-key", "not-existing-file"})
			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Copy output", func() {
		var mockKubecltCmdError error
		mockKubectlCmd := func(opts []string) error {
			return mockKubecltCmdError
		}
		BeforeEach(func() {
			config = &mockClientConfig{}
			cmd.SetKubernetesClientFunc(mockGetKubernetesClientReadyPod)
			cmd.KubectlCmd = mockKubectlCmd

		})
		It("should successfully copy the output", func() {
			path, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			cmd := cmd.NewCopyOutputCommand(config)
			cmd.SetArgs([]string{"--name", "test"})
			err = cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			os.Remove(filepath.Join(path, "test-output"))
		})
		It("should successfully copy the output with output option", func() {
			path, err := os.Getwd()
			output := filepath.Join(path, "test-output")
			Expect(err).NotTo(HaveOccurred())
			cmd := cmd.NewCopyOutputCommand(config)
			cmd.SetArgs([]string{"--name", "test", "--output", output})
			err = cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			os.Remove(output)
		})
		It("should fail if kubectl cp return an error", func() {
			mockKubecltCmdError = fmt.Errorf("test-error")
			cmd := cmd.NewCopyOutputCommand(config)
			cmd.SetArgs([]string{"--name", "test"})
			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
		})
	})
})
