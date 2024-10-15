package main

import "github.com/spf13/pflag"

type CreateOptions struct {
	QuotaUser string

	QuotaNamespace string

	ClusterName string

	UserConfigPath string

	ResourceRequests string

	ResourceLimits string

	CertPath string
}

func (o *CreateOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.QuotaUser, "quota-user", "", "Username for basic authentication to the API server.")
	fs.StringVar(&o.QuotaNamespace, "quota-namespace", "", "The namespace for the nebula cluster request.")
	fs.StringVar(&o.ClusterName, "cluster-name", "", "The name of the kubeconfig cluster to use, run cmd 'kubectl config current-context' to confirm.")
	fs.StringVar(&o.UserConfigPath, "user-config-path", "kube", "Path to the saved kubeconfig file to use for the quota user.")
	fs.StringVar(&o.ResourceRequests, "resource-requests", "", "The compute resource requests total for all Pods in the namespace for the quota user, For example, 'cpu=4,memory=8Gi'.")
	fs.StringVar(&o.ResourceLimits, "resource-limits", "", "The compute resource limits total for all Pods in the namespace for the quota user, For example, 'cpu=8,memory=16Gi'.")
	fs.StringVar(&o.CertPath, "cert-path", "certs", "The directory that the certificate signing request (CSR) and private key will be written.")
}

func (o *CreateOptions) Validate() error {
	if o.QuotaUser == "" {
		return ErrUserNameIsEmpty
	}
	if o.QuotaNamespace == "" {
		return ErrNamespaceIsEmpty
	}
	if o.ClusterName == "" {
		return ErrClusterIsEmpty
	}
	if o.ResourceRequests == "" {
		return ErrResourceRequestsIsEmpty
	}
	if o.ResourceLimits == "" {
		return ErrResourceLimitsIsEmpty
	}
	return nil
}
