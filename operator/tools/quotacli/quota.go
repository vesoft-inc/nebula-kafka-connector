package main

import (
	"context"
	"fmt"

	certv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

type QuotaClient struct {
	client kubernetes.Interface
}

func (r *QuotaClient) CreateNamespace(ctx context.Context, namespace string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err := r.client.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func (r *QuotaClient) CreateRole(ctx context.Context, user, namespace string) error {
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getRoleName(user),
			Namespace: namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
		},
	}
	_, err := r.client.RbacV1().Roles(namespace).Create(ctx, role, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func (r *QuotaClient) CreateRoleBinding(ctx context.Context, user, namespace string) error {
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getRoleBindingName(user),
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "User",
				Name:      user,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     getRoleName(user),
		},
	}
	_, err := r.client.RbacV1().RoleBindings(namespace).Create(ctx, roleBinding, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func (r *QuotaClient) CreateCSR(ctx context.Context, user string, request []byte) error {
	csr := &certv1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: getCsrName(user),
		},
		Spec: certv1.CertificateSigningRequestSpec{
			Groups:     []string{"system:authenticated"},
			Request:    request,
			SignerName: "kubernetes.io/kube-apiserver-client",
			Usages: []certv1.KeyUsage{
				certv1.UsageDigitalSignature,
				certv1.UsageKeyEncipherment,
				certv1.UsageClientAuth,
			},
		},
	}
	_, err := r.client.CertificatesV1().CertificateSigningRequests().Create(ctx, csr, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func (r *QuotaClient) ApproveCSR(ctx context.Context, csrName string) error {
	var csr runtime.Object
	var err error
	csr, err = r.GetCSR(ctx, csrName)
	if apierrors.IsNotFound(err) {
		return fmt.Errorf("could not find v1 version of %s: %v", csrName, err)
	}
	if err != nil {
		return err
	}

	for i := 0; ; i++ {
		modifiedCSR, hasCondition, err := addConditionIfNeeded(csr, string(certv1.CertificateDenied), string(certv1.CertificateApproved), "QuotaCliApprove", "This CSR was approved by quota-cli certificate approve.")
		if err != nil {
			return err
		}
		if !hasCondition {
			if mCSR, ok := modifiedCSR.(*certv1.CertificateSigningRequest); ok {
				_, err = r.client.CertificatesV1().CertificateSigningRequests().UpdateApproval(context.TODO(), mCSR.Name, mCSR, metav1.UpdateOptions{})
			} else {
				return fmt.Errorf("can only handle certificates.k8s.io CertificateSigningRequest objects, got %T", mCSR)
			}

			if apierrors.IsConflict(err) && i < 10 {
				continue
			}
			if err != nil {
				return err
			}
		}
		break
	}
	return nil
}

func (r *QuotaClient) DeleteCSR(ctx context.Context, csrName string) error {
	err := r.client.CertificatesV1().CertificateSigningRequests().Delete(ctx, csrName, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}

func (r *QuotaClient) GetCSR(ctx context.Context, csrName string) (*certv1.CertificateSigningRequest, error) {
	csr, err := r.client.CertificatesV1().CertificateSigningRequests().Get(ctx, csrName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return csr, nil
}

func (r *QuotaClient) CreateResourceQuota(ctx context.Context, namespace string, resourceRequirements corev1.ResourceRequirements) error {
	quota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "compute-resource",
			Namespace: namespace,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				"requests.cpu":    *resourceRequirements.Requests.Cpu(),
				"requests.memory": *resourceRequirements.Requests.Memory(),
				"limits.cpu":      *resourceRequirements.Limits.Cpu(),
				"limits.memory":   *resourceRequirements.Limits.Memory(),
			},
		},
	}
	_, err := r.client.CoreV1().ResourceQuotas(namespace).Create(ctx, quota, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

func addConditionIfNeeded(obj runtime.Object, mustNotHaveConditionType, conditionType, reason, message string) (runtime.Object, bool, error) {
	if csr, ok := obj.(*certv1.CertificateSigningRequest); ok {
		var alreadyHasCondition bool
		for _, c := range csr.Status.Conditions {
			if string(c.Type) == mustNotHaveConditionType {
				return nil, false, fmt.Errorf("certificate signing request %q is already %s", csr.Name, c.Type)
			}
			if string(c.Type) == conditionType {
				alreadyHasCondition = true
			}
		}
		if alreadyHasCondition {
			return csr, true, nil
		}
		csr.Status.Conditions = append(csr.Status.Conditions, certv1.CertificateSigningRequestCondition{
			Type:           certv1.RequestConditionType(conditionType),
			Status:         corev1.ConditionTrue,
			Reason:         reason,
			Message:        message,
			LastUpdateTime: metav1.Now(),
		})
		return csr, false, nil
	} else {
		return csr, false, nil
	}
}
