package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	certv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

// Each level has 2 spaces for PrefixWriter
const (
	LEVEL_0 = iota
	LEVEL_1
	LEVEL_2
	LEVEL_3
	LEVEL_4
)

const (
	NameLabelKey  string = "app.kubernetes.io/name"
	OwnerLabelKey string = "app.kubernetes.io/owner"

	NameLabelValue string = "nebula-graph"
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

func (r *QuotaClient) CreateResourceQuota(ctx context.Context, namespace, user string, resourceRequirements corev1.ResourceRequirements) error {
	quota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "compute-resource",
			Namespace: namespace,
			Labels: map[string]string{
				NameLabelKey:  NameLabelValue,
				OwnerLabelKey: user,
			},
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

func (r *QuotaClient) ListResourceQuotas(ctx context.Context) ([]corev1.ResourceQuota, error) {
	selector := fmt.Sprintf("%s=%s", NameLabelKey, NameLabelValue)
	list, err := r.client.CoreV1().ResourceQuotas(corev1.NamespaceAll).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
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

func describeQuota(resourceQuota *corev1.ResourceQuota) (string, error) {
	return tabbedString(func(out io.Writer) error {
		w := NewPrefixWriter(out)
		w.Write(LEVEL_0, "User:\t%s\n", resourceQuota.Labels[OwnerLabelKey])
		w.Write(LEVEL_0, "Name:\t%s\n", resourceQuota.Name)
		w.Write(LEVEL_0, "Namespace:\t%s\n", resourceQuota.Namespace)
		if len(resourceQuota.Spec.Scopes) > 0 {
			scopes := make([]string, 0, len(resourceQuota.Spec.Scopes))
			for _, scope := range resourceQuota.Spec.Scopes {
				scopes = append(scopes, string(scope))
			}
			sort.Strings(scopes)
			w.Write(LEVEL_0, "Scopes:\t%s\n", strings.Join(scopes, ", "))
			for _, scope := range scopes {
				helpText := helpTextForResourceQuotaScope(corev1.ResourceQuotaScope(scope))
				if len(helpText) > 0 {
					w.Write(LEVEL_0, " * %s\n", helpText)
				}
			}
		}
		w.Write(LEVEL_0, "Resource\tUsed\tHard\n")
		w.Write(LEVEL_0, "--------\t----\t----\n")

		resources := make([]corev1.ResourceName, 0, len(resourceQuota.Status.Hard))
		for resource := range resourceQuota.Status.Hard {
			resources = append(resources, resource)
		}
		sort.Sort(SortableResourceNames(resources))

		msg := "%v\t%v\t%v\n"
		for i := range resources {
			resourceName := resources[i]
			hardQuantity := resourceQuota.Status.Hard[resourceName]
			usedQuantity := resourceQuota.Status.Used[resourceName]
			if hardQuantity.Format != usedQuantity.Format {
				usedQuantity = *resource.NewQuantity(usedQuantity.Value(), hardQuantity.Format)
			}
			w.Write(LEVEL_0, msg, resourceName, usedQuantity.String(), hardQuantity.String())
		}
		return nil
	})
}

func tabbedString(f func(io.Writer) error) (string, error) {
	out := new(tabwriter.Writer)
	buf := &bytes.Buffer{}
	out.Init(buf, 0, 8, 2, ' ', 0)

	err := f(out)
	if err != nil {
		return "", err
	}

	out.Flush()
	return buf.String(), nil
}

func helpTextForResourceQuotaScope(scope corev1.ResourceQuotaScope) string {
	switch scope {
	case corev1.ResourceQuotaScopeTerminating:
		return "Matches all pods that have an active deadline. These pods have a limited lifespan on a node before being actively terminated by the system."
	case corev1.ResourceQuotaScopeNotTerminating:
		return "Matches all pods that do not have an active deadline. These pods usually include long running pods whose container command is not expected to terminate."
	case corev1.ResourceQuotaScopeBestEffort:
		return "Matches all pods that do not have resource requirements set. These pods have a best effort quality of service."
	case corev1.ResourceQuotaScopeNotBestEffort:
		return "Matches all pods that have at least one resource requirement set. These pods have a burstable or guaranteed quality of service."
	default:
		return ""
	}
}
