package kubernetesClient

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewClient() (*kubernetes.Clientset, *rest.Config) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load in-cluster config: %v", err))
	}

	c, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to create clientset: %v", err))
	}

	return c, cfg
}
