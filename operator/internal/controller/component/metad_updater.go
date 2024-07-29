/*
Copyright 2021 Vesoft Inc.

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
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v1alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	utilerrors "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
)

type MetadUpdateManager interface {
	// Update updates the nebula metad
	Update(nc *v1alpha1.NebulaMetad, oldSts, newSts *appsv1.StatefulSet) error
}

type metadUpdater struct {
	podClient kube.Pod
}

func NewMetadUpdater(podClient kube.Pod) MetadUpdateManager {
	return &metadUpdater{podClient: podClient}
}

func (m *metadUpdater) Update(nm *v1alpha1.NebulaMetad, oldSts, newSts *appsv1.StatefulSet) error {
	if *nm.Spec.Replicas == int32(0) {
		return nil
	}

	if !podTemplateEqual(newSts, oldSts) {
		return nil
	}

	if nm.Status.Workload.UpdateRevision == nm.Status.Workload.CurrentRevision &&
		nm.Status.Phase == v1alpha1.RunningPhase {
		return nil
	}

	setPartition(newSts, *oldSts.Spec.UpdateStrategy.RollingUpdate.Partition)
	index, err := getNextUpdatePod(nm.MetadComponent(), *oldSts.Spec.Replicas, m.podClient)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return utilerrors.ReconcileErrorf("%v", err)
		}
		return err
	}
	if index >= 0 {
		return m.updateMetadPod(index, newSts)
	}

	return nil
}

func (m *metadUpdater) updateMetadPod(ordinal int32, newSts *appsv1.StatefulSet) error {
	setPartition(newSts, ordinal)

	return nil
}
