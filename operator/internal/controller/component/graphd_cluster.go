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

package component

import (
	"fmt"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/annotation"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	utilerrors "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
)

type graphdCluster struct {
	clientSet     kube.ClientSet
	updateManager UpdateManager
	eventRecorder record.EventRecorder
}

func NewGraphdCluster(clientSet kube.ClientSet, um UpdateManager, recorder record.EventRecorder) ReconcileManager {
	return &graphdCluster{clientSet: clientSet, updateManager: um, eventRecorder: recorder}
}

func (g *graphdCluster) Reconcile(nc *v2alpha1.NebulaCluster) error {
	if nc.Spec.Graphd == nil {
		return nil
	}
	if err := g.syncGraphdHeadlessService(nc); err != nil {
		return err
	}
	if err := g.syncGraphdService(nc); err != nil {
		return err
	}
	return g.syncGraphdWorkload(nc)
}

func (g *graphdCluster) syncGraphdService(nc *v2alpha1.NebulaCluster) error {
	newSvc := nc.GraphdComponent().GenerateService()
	if newSvc == nil {
		return nil
	}

	return syncService(newSvc, g.clientSet.Service())
}

func (g *graphdCluster) syncGraphdHeadlessService(nc *v2alpha1.NebulaCluster) error {
	newSvc := nc.GraphdComponent().GenerateHeadlessService()
	if newSvc == nil {
		return nil
	}

	return syncService(newSvc, g.clientSet.Service())
}

func (g *graphdCluster) syncGraphdWorkload(nc *v2alpha1.NebulaCluster) error {
	namespace := nc.GetNamespace()
	componentName := nc.GraphdComponent().GetName()

	oldWorkloadTemp, err := g.clientSet.Workload().GetWorkload(namespace, componentName)
	if err != nil && !apierrors.IsNotFound(err) {
		klog.Errorf("get graphd cluster failed: %v", err)
		return err
	}

	notExist := apierrors.IsNotFound(err)
	oldSts := oldWorkloadTemp.DeepCopy()

	cm, cmHash, err := g.syncGraphdConfigMap(nc.DeepCopy())
	if err != nil {
		return err
	}

	// TODO metad endpoints
	newSts, err := nc.GraphdComponent().GenerateWorkload(cm, nil)
	if err != nil {
		klog.Errorf("generate graphd cluster template failed: %v", err)
		return err
	}

	if err := setTemplateAnnotations(newSts, map[string]string{annotation.AnnPodConfigMapHash: cmHash}); err != nil {
		return err
	}

	if err := g.syncNebulaClusterStatus(nc, newSts, oldSts); err != nil {
		return fmt.Errorf("sync graphd cluster status failed: %v", err)
	}

	if notExist {
		if err := setLastAppliedConfigAnnotation(newSts); err != nil {
			return err
		}
		if err := g.clientSet.Workload().CreateWorkload(newSts); err != nil {
			return err
		}
		nc.Status.Graphd.Workload = &v2alpha1.WorkloadStatus{}
		return utilerrors.ReconcileErrorf("waiting for graphd cluster %s running", newSts.GetName())
	}

	equal := podTemplateEqual(newSts, oldSts)
	if !equal || nc.Status.Graphd.Phase == v2alpha1.UpdatePhase {
		if err := g.updateManager.Update(nc, oldSts, newSts); err != nil {
			return err
		}
	}

	if err := g.syncGraphdPVC(nc); err != nil {
		return err
	}

	if equal && nc.GraphdComponent().IsReady() {
		if err := setLastReplicasAnnotation(oldSts); err != nil {
			return err
		}
	}

	return updateWorkload(g.clientSet.Workload(), newSts, oldSts)
}

func (g *graphdCluster) syncNebulaClusterStatus(nc *v2alpha1.NebulaCluster, newSts, oldSts *appsv1.StatefulSet) error {
	if oldSts == nil {
		return nil
	}

	updating, err := isUpdating(nc.GraphdComponent(), g.clientSet.Pod(), oldSts)
	if err != nil {
		return err
	}

	var lastReplicas int32
	val, ok := oldSts.GetAnnotations()[annotation.AnnLastReplicas]
	if ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		lastReplicas = int32(v)
	}

	newReplicas := pointer.Int32Deref(newSts.Spec.Replicas, 0)
	oldReplicas := pointer.Int32Deref(oldSts.Spec.Replicas, 0)
	// TODO metad phase
	if updating {
		nc.Status.Graphd.Phase = v2alpha1.UpdatePhase
	} else if newReplicas < oldReplicas || (ok && newReplicas < lastReplicas) {
		nc.Status.Graphd.Phase = v2alpha1.ScaleInPhase
		if nc.Spec.Graphd.LogVolumeClaim != nil {
			if err := PVCMark(g.clientSet.PVC(), nc.GraphdComponent(), oldReplicas, newReplicas); err != nil {
				return err
			}
		}
	} else if newReplicas > oldReplicas || (ok && newReplicas > lastReplicas) {
		nc.Status.Graphd.Phase = v2alpha1.ScaleOutPhase
	} else {
		nc.Status.Graphd.Phase = v2alpha1.RunningPhase
	}

	return syncComponentStatus(nc.GraphdComponent(), &nc.Status.Graphd.ComponentStatus, oldSts)
}

func (g *graphdCluster) syncGraphdConfigMap(nc *v2alpha1.NebulaCluster) (*corev1.ConfigMap, string, error) {
	return syncConfigMap(
		nc.GraphdComponent(),
		g.clientSet.ConfigMap(),
		v2alpha1.GraphdConfigTemplate,
		nc.GraphdComponent().GetConfigMapKey())
}

func (g *graphdCluster) syncGraphdPVC(nc *v2alpha1.NebulaCluster) error {
	return syncPVC(nc.GraphdComponent(), g.clientSet.PVC())
}

func (g *graphdCluster) Delete(nc *v2alpha1.NebulaCluster) error {
	if nc.Spec.Graphd == nil {
		return nil
	}
	namespace := nc.GetNamespace()
	componentName := nc.GraphdComponent().GetName()
	workload, err := g.clientSet.Workload().GetWorkload(namespace, componentName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("get graphd cluster failed: %v", err)
		return err
	}
	return g.clientSet.Workload().DeleteWorkload(workload)
}
