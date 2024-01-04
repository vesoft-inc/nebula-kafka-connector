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
	"k8s.io/utils/pointer"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
)

type storageScaler struct {
	clientSet kube.ClientSet
}

func NewStorageScaler(clientSet kube.ClientSet) ScaleManager {
	return &storageScaler{clientSet: clientSet}
}

func (s *storageScaler) Scale(nc *v2alpha1.NebulaCluster, oldSts, newSts *appsv1.StatefulSet) error {
	oldReplicas := pointer.Int32Deref(oldSts.Spec.Replicas, 0)
	newReplicas := pointer.Int32Deref(newSts.Spec.Replicas, 0)

	if newReplicas < oldReplicas || nc.Status.Storaged.Phase == v2alpha1.ScaleInPhase {
		return s.ScaleIn(nc, oldReplicas, newReplicas)
	}

	if newReplicas > oldReplicas || nc.Status.Storaged.Phase == v2alpha1.ScaleOutPhase {
		return s.ScaleOut(nc)
	}

	return nil
}

func (s *storageScaler) ScaleOut(nc *v2alpha1.NebulaCluster) error {
	return nil
}

func (s *storageScaler) ScaleIn(nc *v2alpha1.NebulaCluster, oldReplicas, newReplicas int32) error {
	return nil
}
