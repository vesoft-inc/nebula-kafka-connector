/*
Copyright 2023.

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

package component

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v1alpha1"
)

type ReconcileManager interface {
	// Reconcile reconciles the cluster to desired state
	Reconcile(metaClient meta.Client, cluster *v1alpha1.NebulaCluster) error

	// Delete deletes the cluster
	Delete(cluster *v1alpha1.NebulaCluster) error
}

type ScaleManager interface {
	// Scale scales the cluster
	Scale(metaClient meta.Client, nc *v1alpha1.NebulaCluster, oldSts, newSts *appsv1.StatefulSet) error
	// ScaleIn scales in the cluster
	ScaleIn(metaClient meta.Client, nc *v1alpha1.NebulaCluster, oldReplicas, newReplicas int32) error
	// ScaleOut scales out the cluster
	ScaleOut(metaClient meta.Client, nc *v1alpha1.NebulaCluster, oldReplicas, newReplicas int32) error
}

type UpdateManager interface {
	// Update updates the cluster
	Update(nc *v1alpha1.NebulaCluster, oldSts, newSts *appsv1.StatefulSet) error

	// RestartPod restart the specified Pod
	RestartPod(nc *v1alpha1.NebulaCluster, ordinal int32) error
}

type MetaReconcileManager interface {
	// Reconcile reconciles the PV object meta to desired state
	Reconcile(object runtime.Object, enablePVReclaim bool) error
}
