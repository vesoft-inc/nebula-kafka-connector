package main

import (
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func InitClient(configPath string) (kubernetes.Interface, error) {
	c, err := BuildConfigFromPath(configPath)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(c)
}

func BuildConfigFromPath(configPath string) (*restclient.Config, error) {
	configPath = getExplicitFile(configPath)
	if configPath != "" {
		return clientcmd.BuildConfigFromFlags("", configPath)
	}
	return restclient.InClusterConfig()
}
