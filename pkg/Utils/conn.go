package utils

import (
	"fmt"

	crclient "github.com/vegito11/AWSAuthSync/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubeConn(kubeconfig_path *string) (*crclient.Clientset, *kubernetes.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig_path)

	if err != nil {
		fmt.Printf(" Error %s building from flags \n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf(" Error %s getting cluster config \n", err.Error())
		}
	}

	crclientset, err := crclient.NewForConfig(config)

	if err != nil {
		fmt.Printf("Error %s, creating CR clientset \n ", err.Error())
		return nil, nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		fmt.Printf("Error %s, creating clientset \n ", err.Error())
		return nil, nil, err
	}

	return crclientset, clientset, nil
}
