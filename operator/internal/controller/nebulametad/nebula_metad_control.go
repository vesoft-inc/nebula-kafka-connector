package nebulametad

import (
	"context"

	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/controller/component"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
)

type ControlInterface interface {
	UpdateNebulaMetad(metad *v2alpha1.NebulaMetad) error

	DeleteNebulaMetad(metad *v2alpha1.NebulaMetad) error
}

var _ ControlInterface = &defaultNebulaMetadControl{}

func NewMetadControl(
	client client.Client,
	metadCluster component.MetadReconcileManager,
) ControlInterface {
	return &defaultNebulaMetadControl{
		client:       client,
		metadCluster: metadCluster,
	}
}

type defaultNebulaMetadControl struct {
	client       client.Client
	metadCluster component.MetadReconcileManager
}

func (c *defaultNebulaMetadControl) UpdateNebulaMetad(nm *v2alpha1.NebulaMetad) error {
	return nil
}

func (c *defaultNebulaMetadControl) updateNebulaMetad(nm *v2alpha1.NebulaMetad) error {
	if !kube.HasFinalizer(nm, finalizerKey) {
		if err := kube.UpdateFinalizer(context.TODO(), c.client, nm, kube.AddFinalizerOpType, finalizerKey); err != nil {
			return err
		}
	}

	if err := kube.CheckRBAC(context.TODO(), c.client, nm.Namespace); err != nil {
		return err
	}

	if err := c.metadCluster.Reconcile(nm); err != nil {
		klog.Errorf("reconcile metad cluster failed: %v", err)
		return err
	}

	return nil
}

func (c *defaultNebulaMetadControl) DeleteNebulaMetad(nm *v2alpha1.NebulaMetad) error {
	if err := c.metadCluster.Delete(nm); err != nil {
		return err
	}

	if err := component.MetadPVCDeleter(c.client, nm.Namespace, nm.Name); err != nil {
		return err
	}

	return kube.UpdateFinalizer(context.TODO(), c.client, nm, kube.RemoveFinalizerOpType, finalizerKey)
}

func (r *NebulaMetadReconciler) syncNebulaMetad(ctx context.Context, nm *v2alpha1.NebulaMetad) error {
	if nm.DeletionTimestamp != nil {
		return r.control.DeleteNebulaMetad(nm)
	}
	return r.control.UpdateNebulaMetad(nm)
}
