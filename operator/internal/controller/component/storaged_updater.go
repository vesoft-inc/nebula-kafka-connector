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
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	utilerrors "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
)

type storageUpdater struct {
	podClient kube.Pod
}

func NewStorageUpdater(podClient kube.Pod) UpdateManager {
	return &storageUpdater{podClient: podClient}
}

func (g *storageUpdater) Update(nc *v2alpha1.NebulaCluster, oldSts, newSts *appsv1.StatefulSet) error {
	if *nc.Spec.Storaged.Replicas == int32(0) {
		return nil
	}

	// TODO metad phase
	if nc.Status.Storaged.Phase == v2alpha1.ScaleInPhase ||
		nc.Status.Storaged.Phase == v2alpha1.ScaleOutPhase {
		return setLastConfig(oldSts, newSts)
	}

	// template had been changed
	if !podTemplateEqual(newSts, oldSts) {
		return nil
	}

	if nc.Status.Storaged.Workload.UpdateRevision == nc.Status.Storaged.Workload.CurrentRevision &&
		nc.Status.Storaged.Phase == v2alpha1.RunningPhase {
		return nil
	}

	setPartition(newSts, *oldSts.Spec.UpdateStrategy.RollingUpdate.Partition)
	index, err := getNextUpdatePod(nc.StoragedComponent(), *oldSts.Spec.Replicas, g.podClient)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return utilerrors.ReconcileErrorf("%v", err)
		}
		return err
	}
	if index >= 0 {
		return g.updateStoragedPod(index, newSts)
	}

	return nil
}

func (g *storageUpdater) RestartPod(nc *v2alpha1.NebulaCluster, ordinal int32) error {
	return nil
}

func (g *storageUpdater) updateStoragedPod(ordinal int32, newSts *appsv1.StatefulSet) error {
	setPartition(newSts, ordinal)
	return nil
}
