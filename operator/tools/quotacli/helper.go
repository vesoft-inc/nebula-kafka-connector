package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog/v2"
)

var (
	ErrUserNameIsEmpty         = errors.New("user name is empty")
	ErrNamespaceIsEmpty        = errors.New("namespace is empty")
	ErrClusterIsEmpty          = errors.New("cluster is empty")
	ErrResourceRequestsIsEmpty = errors.New("compute resource requests is empty")
	ErrResourceLimitsIsEmpty   = errors.New("compute resource limits is empty")
)

func getRoleName(user string) string {
	return user + "-full-access"
}

func getRoleBindingName(user string) string {
	return user + "-role-binding"
}

func getCsrName(user string) string {
	return user + "-csr"
}

func getContextName(user string) string {
	return user + "-context"
}

func getKeyFileName(user string) string {
	return user + ".key"
}

func getCsrFileName(user string) string {
	return user + ".csr"
}

func getUserConfig(user string) string {
	return user + "-kubeconfig"
}

func encodeToString(fileName string) string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		klog.Errorf("failed to read file: %v", err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(data)
}

func validate(opts *Options) error {
	if opts.QuotaUser == "" {
		return ErrUserNameIsEmpty
	}
	if opts.QuotaNamespace == "" {
		return ErrNamespaceIsEmpty
	}
	if opts.ClusterName == "" {
		return ErrClusterIsEmpty
	}
	if opts.ResourceRequests == "" {
		return ErrResourceRequestsIsEmpty
	}
	if opts.ResourceLimits == "" {
		return ErrResourceLimitsIsEmpty
	}
	return nil
}

func handleResourceRequirementsV1(params map[string]string) (corev1.ResourceRequirements, error) {
	result := corev1.ResourceRequirements{}
	limits, err := populateResourceListV1(params["limits"])
	if err != nil {
		return result, err
	}
	result.Limits = limits
	requests, err := populateResourceListV1(params["requests"])
	if err != nil {
		return result, err
	}
	result.Requests = requests
	return result, nil
}

func populateResourceListV1(spec string) (corev1.ResourceList, error) {
	result := corev1.ResourceList{}
	resourceStatements := strings.Split(spec, ",")
	for _, resourceStatement := range resourceStatements {
		parts := strings.Split(resourceStatement, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid argument syntax %v, expected <resource>=<value>", resourceStatement)
		}
		resourceName := corev1.ResourceName(parts[0])
		resourceQuantity, err := resource.ParseQuantity(parts[1])
		if err != nil {
			return nil, err
		}
		result[resourceName] = resourceQuantity
	}
	return result, nil
}
