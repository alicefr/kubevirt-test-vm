package cmd

import (
	"context"
	"fmt"

	k8scorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	OutputDir = "/output"
)

func PvcOutputName(name string) string {
	if name == "" {
		return "fio-output"
	}
	return fmt.Sprintf("fio-output-%s", name)
}

func CreateOutputPVC(client *kubernetes.Clientset, pvcClaim, namespace string, labels map[string]string) error {
	vMode := k8scorev1.PersistentVolumeFilesystem
	requestPVCs := map[k8scorev1.ResourceName]resource.Quantity{
		k8scorev1.ResourceStorage: resource.MustParse("1G"),
	}
	// Create pvc for the output
	pvc := &k8scorev1.PersistentVolumeClaim{
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:   pvcClaim,
			Labels: labels,
		},
		Spec: k8scorev1.PersistentVolumeClaimSpec{
			VolumeMode:  &vMode,
			AccessModes: []k8scorev1.PersistentVolumeAccessMode{k8scorev1.ReadWriteOnce},
			Resources: k8scorev1.ResourceRequirements{
				Requests: requestPVCs,
			},
		},
	}
	_, err := client.CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(), pvc, k8smetav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
		fmt.Printf("PVC %s already exists \n", pvcName)
	}
	if err != nil {
		fmt.Printf("Created PVC %s \n", pvcName)
	}
	return nil

}

func GetKubernetesClient(clientConfig clientcmd.ClientConfig) (*kubernetes.Clientset, error) {
	conf, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(conf)
}
