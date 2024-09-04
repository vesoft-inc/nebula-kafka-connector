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

package nebulacluster

import (
	"context"
	"fmt"
	"strings"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	errorutils "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v1alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/controller/component"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/controller/component/reclaimer"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	utilerrors "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
)

const (
	finalizer = "apps.nebula-graph-ng.io/cluster-cleanup"
)

type ControlInterface interface {
	UpdateNebulaCluster(cluster *v1alpha1.NebulaCluster) error

	DeleteCluster(cluster *v1alpha1.NebulaCluster) error
}

var _ ControlInterface = &defaultNebulaClusterControl{}

func NewDefaultNebulaClusterControl(
	client client.Client,
	clientSet kube.ClientSet,
	graphdCluster component.ReconcileManager,
	storagedCluster component.ReconcileManager,
	exporter component.ReconcileManager,
	console component.ReconcileManager,
	metaReconciler component.MetaReconcileManager,
	pvcReclaimer reclaimer.PVCReclaimer,
	conditionUpdater ClusterConditionUpdater,
) ControlInterface {
	return &defaultNebulaClusterControl{
		client:           client,
		clientSet:        clientSet,
		graphdCluster:    graphdCluster,
		storagedCluster:  storagedCluster,
		exporter:         exporter,
		console:          console,
		metaReconciler:   metaReconciler,
		pvcReclaimer:     pvcReclaimer,
		conditionUpdater: conditionUpdater,
	}
}

type defaultNebulaClusterControl struct {
	client           client.Client
	clientSet        kube.ClientSet
	graphdCluster    component.ReconcileManager
	storagedCluster  component.ReconcileManager
	exporter         component.ReconcileManager
	console          component.ReconcileManager
	metaReconciler   component.MetaReconcileManager
	pvcReclaimer     reclaimer.PVCReclaimer
	conditionUpdater ClusterConditionUpdater
}

func (c *defaultNebulaClusterControl) UpdateNebulaCluster(nc *v1alpha1.NebulaCluster) error {
	var errs []error
	oldStatus := nc.Status.DeepCopy()

	if err := c.updateNebulaCluster(nc); err != nil {
		errs = append(errs, err)
	}

	c.conditionUpdater.Update(nc)
	nc.Status.ObservedGeneration = nc.Generation
	if err := c.clientSet.NebulaCluster().UpdateNebulaClusterStatus(nc.DeepCopy()); err != nil {
		errs = append(errs, err)
	}

	if apiequality.Semantic.DeepEqual(&nc.Status, oldStatus) && nc.IsReady() {
		return errorutils.NewAggregate(errs)
	}

	if !nc.IsConditionReady() {
		errs = append(errs, utilerrors.ReconcileErrorf("waiting for nebulacluster ready"))
	}

	return errorutils.NewAggregate(errs)
}

func (c *defaultNebulaClusterControl) DeleteCluster(nc *v1alpha1.NebulaCluster) error {
	if err := c.graphdCluster.Delete(nc); err != nil {
		return err
	}
	if err := c.storagedCluster.Delete(nc); err != nil {
		return err
	}
	if err := component.PVCDeleter(c.client, nc.Namespace, nc.Name); err != nil {
		return err
	}

	if nc.Spec.MetadRef == nil {
		return component.ErrorMetadReferenceIsNil
	}

	metad, err := c.clientSet.NebulaMetad().GetNebulaMetad(nc.GetMetadNamespace(), nc.Spec.MetadRef.Name)
	if err != nil {
		klog.Errorf("failed to get metad cluster: %v", err)
		return err
	}
	if !metad.MetadComponent().IsReady() {
		return fmt.Errorf("metad [%s/%s] is not ready", metad.Namespace, metad.Name)
	}

	username, password, err := kube.GetCredential(c.clientSet, nc.Namespace, nc.Spec.CredentialSecret)
	if err != nil {
		return err
	}
	metadEndpoints := metad.MetadComponent().GetEndpoints(v1alpha1.MetadPortNameGRPC)
	metaClient, err := meta.NewMetaClient(strings.Join(metadEndpoints, ","), meta.WithUserPassword(username, password))
	if err != nil {
		return err
	}
	if _, err := metaClient.Login(); err != nil {
		klog.Errorf("login metad failed: %v", err)
		return err
	}
	defer func() {
		metaClient.Close()
	}()

	req := meta.NewDropClusterReq(nc.Name, true)
	if err := metaClient.DropCluster(req); err != nil {
		if ne, ok := err.(*nebula.NebulaError); ok {
			if ne.Code() != nebula.ERROR_META_CLUSTER_NOT_FOUND {
				klog.Errorf("drop cluster failed: %v", err)
				return err
			}
		} else {
			klog.Errorf("drop cluster got unkonw error: %v", err)
			return err
		}
	}

	return kube.UpdateFinalizer(context.TODO(), c.client, nc, kube.RemoveFinalizerOpType, finalizer)
}

func (c *defaultNebulaClusterControl) updateNebulaCluster(nc *v1alpha1.NebulaCluster) error {
	if !kube.HasFinalizer(nc, finalizer) {
		if err := kube.UpdateFinalizer(context.TODO(), c.client, nc, kube.AddFinalizerOpType, finalizer); err != nil {
			return err
		}
	}

	if err := kube.CheckRBAC(context.TODO(), c.client, nc.Namespace); err != nil {
		return err
	}

	if nc.Spec.MetadRef == nil {
		return component.ErrorMetadReferenceIsNil
	}

	metad, err := c.clientSet.NebulaMetad().GetNebulaMetad(nc.GetMetadNamespace(), nc.Spec.MetadRef.Name)
	if err != nil {
		klog.Errorf("failed to get metad cluster: %v", err)
		return err
	}
	if !metad.MetadComponent().IsReady() {
		return fmt.Errorf("metad [%s/%s] is not ready", metad.Namespace, metad.Name)
	}

	username, password, err := kube.GetCredential(c.clientSet, nc.Namespace, nc.Spec.CredentialSecret)
	if err != nil {
		return err
	}
	metadEndpoints := metad.MetadComponent().GetEndpoints(v1alpha1.MetadPortNameGRPC)
	metaClient, err := meta.NewMetaClient(strings.Join(metadEndpoints, ","), meta.WithUserPassword(username, password))
	if err != nil {
		return err
	}
	if _, err := metaClient.Login(); err != nil {
		klog.Errorf("login metad failed: %v", err)
		return err
	}
	defer func() {
		metaClient.Close()
	}()

	if !nc.Status.CreatedDone {
		req := meta.NewCreateClusterReq(nc.Name, int(nc.Spec.ReplicaFactor), "root", nc.Spec.Zones)
		if err = metaClient.CreateCluster(req); err != nil {
			if ne, ok := err.(*nebula.NebulaError); ok {
				if ne.Code() != nebula.ERROR_META_CLUSTER_ALREADY_EXISTS {
					klog.Errorf("create cluster failed: %v", err)
					return err
				}
			} else {
				klog.Errorf("create cluster got unkonw error: %v", err)
				return err
			}
		}
		nc.Status.CreatedDone = true
	}

	nc.SetOwnerReferences(metad.GenerateOwnerReferences())
	if err := c.clientSet.NebulaCluster().UpdateNebulaCluster(nc); err != nil {
		klog.Errorf("set cluster owner references failed: %v", err)
		return err
	}

	if err := c.storagedCluster.Reconcile(metaClient, nc); err != nil {
		klog.Errorf("reconcile storaged cluster failed: %v", err)
		return err
	}

	if err := c.graphdCluster.Reconcile(metaClient, nc); err != nil {
		klog.Errorf("reconcile graphd cluster failed: %v", err)
		return err
	}

	if err := c.console.Reconcile(metaClient, nc); err != nil {
		klog.Errorf("reconcile console failed: %v", err)
		return err
	}

	if err := c.exporter.Reconcile(metaClient, nc); err != nil {
		klog.Errorf("reconcile exporter failed: %v", err)
		return err
	}

	if err := c.metaReconciler.Reconcile(nc, nc.IsPVReclaimEnabled()); err != nil {
		klog.Errorf("reconcile pv and pvc metadata failed: %v", err)
		return err
	}

	if err := c.pvcReclaimer.Reclaim(nc); err != nil {
		klog.Errorf("reclaim pvc failed: %v", err)
		return err
	}

	return nil
}
