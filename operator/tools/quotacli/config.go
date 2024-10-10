package main

import (
	"encoding/base64"
	"os"
	"path"

	certv1 "k8s.io/api/certificates/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type (
	KubeConfig struct {
		ApiVersion     string         `yaml:"apiversion"`
		Clusters       []NamedCluster `yaml:"clusters"`
		Contexts       []NamedContext `yaml:"contexts"`
		CurrentContext string         `yaml:"current-context"`
		Kind           string         `yaml:"kind"`
		Users          []NamedUser    `yaml:"users"`
	}

	NamedCluster struct {
		Name    string  `yaml:"name"`
		Cluster Cluster `yaml:"cluster"`
	}

	Cluster struct {
		CertificateAuthorityData string `yaml:"certificate-authority-data"`
		Server                   string `yaml:"server"`
	}

	NamedContext struct {
		Name    string  `yaml:"name"`
		Context Context `yaml:"context"`
	}

	Context struct {
		Cluster   string `yaml:"cluster"`
		Namespace string `yaml:"namespace"`
		AuthInfo  string `yaml:"user"`
	}

	NamedUser struct {
		Name     string   `yaml:"name"`
		AuthInfo AuthInfo `yaml:"user"`
	}

	AuthInfo struct {
		ClientCertificateData string `yaml:"client-certificate-data"`
		ClientKeyData         string `yaml:"client-key-data"`
	}
)

func LoadConfig(configPath string) (*clientcmdapi.Config, error) {
	configPath = getExplicitFile(configPath)
	config, err := clientcmd.LoadFromFile(configPath)
	if err != nil {
		return nil, err
	}
	if err = clientcmdapi.MinifyConfig(config); err != nil {
		return nil, err
	}
	if err = clientcmdapi.FlattenConfig(config); err != nil {
		return nil, err
	}
	return config, nil
}

func WriteConfigToFile(content []byte, dir, filename string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	fullPath := path.Join(dir, filename)
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return err
	}
	return nil
}

func getExplicitFile(configPath string) string {
	if configPath == "" {
		configPath = os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	}
	if configPath == "" {
		if _, err := os.Stat(clientcmd.RecommendedHomeFile); err == nil || os.IsExist(err) {
			configPath = clientcmd.RecommendedHomeFile
		}
	}
	return configPath
}

func convert(apiConfig *clientcmdapi.Config, csr *certv1.CertificateSigningRequest,
	key, cluster, user, namespace string) *KubeConfig {
	config := &KubeConfig{
		ApiVersion:     "v1",
		Kind:           "Config",
		CurrentContext: getContextName(user),
	}
	config.Clusters = append(config.Clusters, NamedCluster{
		Name: cluster,
		Cluster: Cluster{
			CertificateAuthorityData: base64.StdEncoding.EncodeToString(apiConfig.Clusters[cluster].CertificateAuthorityData),
			Server:                   apiConfig.Clusters[cluster].Server,
		},
	})
	config.Contexts = append(config.Contexts, NamedContext{
		Name: getContextName(user),
		Context: Context{
			Cluster:   cluster,
			Namespace: namespace,
			AuthInfo:  user,
		},
	})
	config.Users = append(config.Users, NamedUser{
		Name: user,
		AuthInfo: AuthInfo{
			ClientCertificateData: base64.StdEncoding.EncodeToString(csr.Status.Certificate),
			ClientKeyData:         key,
		},
	})
	return config
}
