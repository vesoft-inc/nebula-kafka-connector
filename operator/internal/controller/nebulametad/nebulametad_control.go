/*
Copyright 2023 Vesoft Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nebulametad

import (
	"context"
	"fmt"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	errorutils "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v1alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/controller/component"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/controller/component/reclaimer"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	utilerrors "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
)

type ControlInterface interface {
	UpdateNebulaMetad(metad *v1alpha1.NebulaMetad) error

	DeleteNebulaMetad(metad *v1alpha1.NebulaMetad) error
}

var _ ControlInterface = &defaultNebulaMetadControl{}

func NewMetadControl(
	client client.Client,
	metadClient kube.NebulaMetad,
	metadCluster component.MetadReconcileManager,
	metaReconciler component.MetaReconcileManager,
	pvcReclaimer reclaimer.PVCReclaimer,
	conditionUpdater MetadConditionUpdater,
) ControlInterface {
	return &defaultNebulaMetadControl{
		client:           client,
		metadClient:      metadClient,
		metadCluster:     metadCluster,
		metaReconciler:   metaReconciler,
		pvcReclaimer:     pvcReclaimer,
		conditionUpdater: conditionUpdater,
	}
}

type defaultNebulaMetadControl struct {
	client           client.Client
	metadClient      kube.NebulaMetad
	metadCluster     component.MetadReconcileManager
	metaReconciler   component.MetaReconcileManager
	pvcReclaimer     reclaimer.PVCReclaimer
	conditionUpdater MetadConditionUpdater
}

func (c *defaultNebulaMetadControl) UpdateNebulaMetad(nm *v1alpha1.NebulaMetad) error {
	var errs []error
	oldStatus := nm.Status.DeepCopy()

	if err := c.updateNebulaMetad(nm); err != nil {
		errs = append(errs, err)
	}

	c.conditionUpdater.Update(nm)
	nm.Status.ObservedGeneration = nm.Generation
	if err := c.metadClient.UpdateNebulaMetadStatus(nm.DeepCopy()); err != nil {
		errs = append(errs, err)
	}

	if apiequality.Semantic.DeepEqual(&nm.Status, oldStatus) && nm.IsReady() {
		return errorutils.NewAggregate(errs)
	}

	if !nm.IsConditionReady() {
		errs = append(errs, utilerrors.ReconcileErrorf("waiting for nebulametad ready"))
	}

	return errorutils.NewAggregate(errs)
}

func (c *defaultNebulaMetadControl) updateNebulaMetad(nm *v1alpha1.NebulaMetad) error {
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

	if err := c.metaReconciler.Reconcile(nm, nm.IsPVReclaimEnabled()); err != nil {
		klog.Errorf("reconcile metad cluster pv and pvc metadata failed: %v", err)
		return err
	}

	if err := c.pvcReclaimer.Reclaim(nm); err != nil {
		klog.Errorf("reclaim metad cluster pvc failed: %v", err)
		return err
	}

	return nil
}

func (c *defaultNebulaMetadControl) DeleteNebulaMetad(nm *v1alpha1.NebulaMetad) error {
	if nm.Status.ManagedClusters > 0 {
		return fmt.Errorf("managed clusters is %d, cannot be deleted now", nm.Status.ManagedClusters)
	}
	if err := c.metadCluster.Delete(nm); err != nil {
		return err
	}
	if err := component.PVCDeleter(c.client, nm.Namespace, nm.Name); err != nil {
		return err
	}
	return kube.UpdateFinalizer(context.TODO(), c.client, nm, kube.RemoveFinalizerOpType, finalizerKey)
}
