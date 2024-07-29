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
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v1alpha1"
)

const (
	WorkloadReady = "Ready"
	// WorkloadNotUpToDate is added when one of workloads is not up-to-date.
	WorkloadNotUpToDate = "WorkloadNotUpToDate"
	// StoragedUnhealthy is added when one of storaged pods is unhealthy.
	StoragedUnhealthy = "StoragedUnhealthy"
	// GraphdUnhealthy is added when one of graphd pods is unhealthy.
	GraphdUnhealthy = "GraphdUnhealthy"
)

type ClusterConditionUpdater interface {
	Update(cluster *v1alpha1.NebulaCluster)
}

var _ ClusterConditionUpdater = &nebulaClusterConditionUpdater{}

func NewClusterConditionUpdater() ClusterConditionUpdater {
	return &nebulaClusterConditionUpdater{}
}

type nebulaClusterConditionUpdater struct{}

func (u *nebulaClusterConditionUpdater) Update(nc *v1alpha1.NebulaCluster) {
	u.updateReadyCondition(nc)
}

func allWorkloadsAreUpToDate(nc *v1alpha1.NebulaCluster) bool {
	isUpToDate := func(status *v1alpha1.WorkloadStatus, requireExist bool) bool {
		if status == nil {
			return !requireExist
		}
		return status.CurrentRevision == status.UpdateRevision
	}

	updated := (isUpToDate(nc.Status.Storaged.Workload, false)) &&
		(isUpToDate(nc.Status.Graphd.Workload, false))

	return updated
}

func (u *nebulaClusterConditionUpdater) updateReadyCondition(nc *v1alpha1.NebulaCluster) {
	status := metav1.ConditionFalse
	var reason string
	var message string

	switch {
	case !allWorkloadsAreUpToDate(nc):
		reason = WorkloadNotUpToDate
		message = "Workload is in progress"
	case !nc.StoragedComponent().IsReady():
		reason = StoragedUnhealthy
		message = "Storaged is not healthy"
	case !nc.GraphdComponent().IsReady():
		reason = GraphdUnhealthy
		message = "Graphd is not healthy"
	default:
		status = metav1.ConditionTrue
		reason = WorkloadReady
		message = "Nebula cluster is running"
	}

	cond := metav1.Condition{
		Type:    v1alpha1.NebulaClusterReady,
		Status:  status,
		Reason:  reason,
		Message: message,
	}
	meta.SetStatusCondition(&nc.Status.Conditions, cond)
}
