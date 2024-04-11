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

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
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

func (s *storagedCluster) Reconcile(metaClient meta.Client, nc *v2alpha1.NebulaCluster) error {
	if nc.Spec.Storaged == nil {
		return nil
	}
	if err := s.syncStoragedHeadlessService(nc); err != nil {
		return err
	}
	return s.syncStoragedWorkload(metaClient, nc)
}

func (s *storagedCluster) syncStoragedHeadlessService(nc *v2alpha1.NebulaCluster) error {
	newSvc := nc.StoragedComponent().GenerateHeadlessService()
	if newSvc == nil {
		return nil
	}

	return syncService(newSvc, s.clientSet.Service())
}

func (s *storagedCluster) syncStoragedWorkload(metaClient meta.Client, nc *v2alpha1.NebulaCluster) error {
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

	metad, err := s.clientSet.NebulaMetad().GetNebulaMetad(nc.GetMetadNamespace(), nc.Spec.MetadRef.Name)
	if err != nil {
		return err
	}
	metadEndpoints := metad.MetadComponent().GetEndpoints(v2alpha1.MetadPortNameGRPC)
	newSts, err := nc.StoragedComponent().GenerateWorkload(cm, metadEndpoints)
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

	if !nc.Status.Storaged.ServicesAdded {
		if err := addStorageServices(metaClient, nc, *oldSts.Spec.Replicas, *newSts.Spec.Replicas); err != nil {
			return err
		}
		klog.Infof("storaged cluster [%s/%s] add services succeed", namespace, componentName)
		nc.Status.Storaged.ServicesAdded = true
	}
	if !nc.Status.Storaged.InitPartitionDone {
		if err := s.initCluster(metaClient, nc.Name); err != nil {
			return err
		}
		klog.Infof("init storaged cluster [%s/%s] partitions succeed", namespace, componentName)
		nc.Status.Storaged.InitPartitionDone = true
	}

	if err := s.scaleManager.Scale(metaClient, nc, oldSts, newSts); err != nil {
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

func (s *storagedCluster) initCluster(metaClient meta.Client, clusterName string) error {
	req := meta.NewInitClusterReq(clusterName)
	if err := metaClient.InitCluster(req); err != nil {
		if ne, ok := err.(*nebula.NebulaError); ok {
			// TODO [NM007]: Part already inited
			if ne.Code() != "NM007" {
				klog.Errorf("init cluster failed: %v", err)
				return err
			}
		} else {
			klog.Errorf("init cluster got unkonw error: %v", err)
			return err
		}
	}
	return nil
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

func addStorageServices(metaClient meta.Client, nc *v2alpha1.NebulaCluster, oldReplicas, newReplicas int32) error {
	var start int32
	if newReplicas > oldReplicas {
		start = oldReplicas
	}

	port := nc.StoragedComponent().GetPort(v2alpha1.StoragedPortNameGRPC)
	for i := start; i < newReplicas; i++ {
		host := nc.StoragedComponent().GetPodFQDN(i)
		req := meta.NewAddServiceReq(host, uint32(port), meta.ServiceTypeStoraged, nc.Name)
		if err := metaClient.AddService(req); err != nil {
			if ne, ok := err.(*nebula.NebulaError); ok {
				// TODO [NM019]: Service with static port already exists
				if ne.Code() != "NM019" {
					klog.Errorf("add storaged service failed: %v", err)
					return err
				}
			} else {
				klog.Errorf("add service got unkonw error: %v", err)
				return err
			}
		}
	}

	return nil
}
