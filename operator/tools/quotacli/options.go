package main

type CreateOptions struct {
	QuotaUser string

	QuotaNamespace string

	ClusterName string

	UserConfigPath string

	ResourceRequests string

	ResourceLimits string

	CertPath string
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
