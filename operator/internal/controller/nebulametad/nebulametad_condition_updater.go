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
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
)

const (
	WorkloadReady = "Ready"
	// MetadNotUpToDate is added when metad is not up-to-date.
	MetadNotUpToDate = "MetadNotUpToDate"
	// MetaddUnhealthy is added when one of metad pods is unhealthy.
	MetaddUnhealthy = "MetaddUnhealthy"
)

type MetadConditionUpdater interface {
	Update(metad *v2alpha1.NebulaMetad)
}

var _ MetadConditionUpdater = &nebulaMetadConditionUpdater{}

func NewClusterConditionUpdater() MetadConditionUpdater {
	return &nebulaMetadConditionUpdater{}
}

type nebulaMetadConditionUpdater struct{}

func (u *nebulaMetadConditionUpdater) Update(nm *v2alpha1.NebulaMetad) {
	u.updateReadyCondition(nm)
}

func workloadIsUpToDate(nm *v2alpha1.NebulaMetad) bool {
	isUpToDate := func(status *v2alpha1.WorkloadStatus, requireExist bool) bool {
		if status == nil {
			return !requireExist
		}
		return status.CurrentRevision == status.UpdateRevision
	}
	return isUpToDate(nm.Status.Workload, false)
}

func (u *nebulaMetadConditionUpdater) updateReadyCondition(nm *v2alpha1.NebulaMetad) {
	status := metav1.ConditionFalse
	var reason string
	var message string

	switch {
	case !workloadIsUpToDate(nm):
		reason = MetadNotUpToDate
		message = "Workload is in progress"
	case !nm.MetadComponent().IsReady():
		reason = MetaddUnhealthy
		message = "Metad is not healthy"
	default:
		status = metav1.ConditionTrue
		reason = WorkloadReady
		message = "Nebula metad is running"
	}

	cond := metav1.Condition{
		Type:    v2alpha1.NebulaMetadReady,
		Status:  status,
		Reason:  reason,
		Message: message,
	}
	meta.SetStatusCondition(&nm.Status.Conditions, cond)
}
