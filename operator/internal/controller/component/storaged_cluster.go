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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/annotation"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	utilerrors "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
)

type storagedCluster struct {
	clientSet     kube.ClientSet
	scaleManager  ScaleManager
	updateManager UpdateManager
	eventRecorder record.EventRecorder
}

func NewStoragedCluster(clientSet kube.ClientSet, sm ScaleManager, um UpdateManager, recorder record.EventRecorder) ReconcileManager {
	return &storagedCluster{
		clientSet:     clientSet,
		scaleManager:  sm,
		updateManager: um,
		eventRecorder: recorder,
	}
}

func (s *storagedCluster) Reconcile(nc *v2alpha1.NebulaCluster) error {
	if nc.Spec.Storaged == nil {
		return nil
	}
	if err := s.syncStoragedHeadlessService(nc); err != nil {
		return err
	}
	return s.syncStoragedWorkload(nc)
}

func (s *storagedCluster) syncStoragedHeadlessService(nc *v2alpha1.NebulaCluster) error {
	newSvc := nc.StoragedComponent().GenerateHeadlessService()
	if newSvc == nil {
		return nil
	}

	return syncService(newSvc, s.clientSet.Service())
}

func (s *storagedCluster) syncStoragedWorkload(nc *v2alpha1.NebulaCluster) error {
	namespace := nc.GetNamespace()
	componentName := nc.StoragedComponent().GetName()

	oldWorkloadTemp, err := s.clientSet.Workload().GetWorkload(namespace, componentName)
	if err != nil && !apierrors.IsNotFound(err) {
		klog.Errorf("get storaged cluster failed: %v", err)
		return err
	}

	notExist := apierrors.IsNotFound(err)
	oldSts := oldWorkloadTemp.DeepCopy()

	cm, cmHash, err := s.syncStoragedConfigMap(nc.DeepCopy())
	if err != nil {
		return err
	}

	// TODO metad endpoints
	newSts, err := nc.StoragedComponent().GenerateWorkload(cm, nil)
	if err != nil {
		klog.Errorf("generate storaged cluster template failed: %v", err)
		return err
	}

	if err := setTemplateAnnotations(newSts, map[string]string{annotation.AnnPodConfigMapHash: cmHash}); err != nil {
		return err
	}

	if err := s.syncNebulaClusterStatus(nc, oldSts); err != nil {
		return fmt.Errorf("sync storaged cluster status failed: %v", err)
	}

	if notExist {
		if err := setLastAppliedConfigAnnotation(newSts); err != nil {
			return err
		}
		if err := s.clientSet.Workload().CreateWorkload(newSts); err != nil {
			return err
		}
		nc.Status.Storaged.Workload = &v2alpha1.WorkloadStatus{}
		return utilerrors.ReconcileErrorf("waiting for storaged cluster %s running", newSts.GetName())
	}

	if err := s.scaleManager.Scale(nc, oldSts, newSts); err != nil {
		klog.Errorf("scale storaged cluster [%s/%s] failed: %v", namespace, componentName, err)
		return err
	}

	equal := podTemplateEqual(newSts, oldSts)
	if !equal || nc.Status.Storaged.Phase == v2alpha1.UpdatePhase {
		if err := s.updateManager.Update(nc, oldSts, newSts); err != nil {
			return err
		}
	}

	if err := s.syncStoragedPVC(nc); err != nil {
		return err
	}

	if equal && nc.StoragedComponent().IsReady() {
		if err := setLastReplicasAnnotation(oldSts); err != nil {
			return err
		}
	}

	return updateWorkload(s.clientSet.Workload(), newSts, oldSts)
}

func (s *storagedCluster) syncNebulaClusterStatus(nc *v2alpha1.NebulaCluster, oldSts *appsv1.StatefulSet) error {
	if oldSts == nil {
		return nil
	}

	if nc.Status.Storaged.Phase == "" {
		nc.Status.Storaged.Phase = v2alpha1.RunningPhase
	}

	updating, err := isUpdating(nc.StoragedComponent(), s.clientSet.Pod(), oldSts)
	if err != nil {
		return err
	}

	// TODO metad phase
	if updating {
		nc.Status.Storaged.Phase = v2alpha1.UpdatePhase
	}

	return syncComponentStatus(nc.StoragedComponent(), &nc.Status.Storaged.ComponentStatus, oldSts)
}

func (s *storagedCluster) syncStoragedConfigMap(nc *v2alpha1.NebulaCluster) (*corev1.ConfigMap, string, error) {
	return syncConfigMap(
		nc.StoragedComponent(),
		s.clientSet.ConfigMap(),
		v2alpha1.StoragedConfigTemplate,
		nc.StoragedComponent().GetConfigMapKey())
}

func (s *storagedCluster) syncStoragedPVC(nc *v2alpha1.NebulaCluster) error {
	return syncPVC(nc.StoragedComponent(), s.clientSet.PVC())
}

func (s *storagedCluster) Delete(nc *v2alpha1.NebulaCluster) error {
	if nc.Spec.Storaged == nil {
		return nil
	}
	namespace := nc.GetNamespace()
	componentName := nc.StoragedComponent().GetName()
	workload, err := s.clientSet.Workload().GetWorkload(namespace, componentName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("get storaged cluster failed: %v", err)
		return err
	}
	return s.clientSet.Workload().DeleteWorkload(workload)
}
