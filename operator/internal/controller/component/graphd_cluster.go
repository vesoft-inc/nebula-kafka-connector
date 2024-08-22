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
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v1alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/annotation"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	utilerrors "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
)

type graphdCluster struct {
	clientSet     kube.ClientSet
	scaleManager  ScaleManager
	updateManager UpdateManager
	eventRecorder record.EventRecorder
}

func NewGraphdCluster(clientSet kube.ClientSet, sm ScaleManager, um UpdateManager, recorder record.EventRecorder) ReconcileManager {
	return &graphdCluster{
		clientSet:     clientSet,
		scaleManager:  sm,
		updateManager: um,
		eventRecorder: recorder}
}

func (g *graphdCluster) Reconcile(metaClient meta.Client, nc *v1alpha1.NebulaCluster) error {
	if nc.Spec.Graphd == nil {
		return nil
	}
	if err := g.syncGraphdHeadlessService(nc); err != nil {
		return err
	}
	if err := g.syncGraphdService(nc); err != nil {
		return err
	}
	return g.syncGraphdWorkload(metaClient, nc)
}

func (g *graphdCluster) syncGraphdService(nc *v1alpha1.NebulaCluster) error {
	newSvc := nc.GraphdComponent().GenerateService()
	if newSvc == nil {
		return nil
	}

	return syncService(newSvc, g.clientSet.Service())
}

func (g *graphdCluster) syncGraphdHeadlessService(nc *v1alpha1.NebulaCluster) error {
	newSvc := nc.GraphdComponent().GenerateHeadlessService()
	if newSvc == nil {
		return nil
	}

	return syncService(newSvc, g.clientSet.Service())
}

func (g *graphdCluster) syncGraphdWorkload(metaClient meta.Client, nc *v1alpha1.NebulaCluster) error {
	namespace := nc.GetNamespace()
	componentName := nc.GraphdComponent().GetName()

	oldWorkloadTemp, err := g.clientSet.Workload().GetWorkload(namespace, componentName)
	if err != nil && !apierrors.IsNotFound(err) {
		klog.Errorf("get graphd cluster failed: %v", err)
		return err
	}

	notExist := apierrors.IsNotFound(err)
	oldSts := oldWorkloadTemp.DeepCopy()

	needSuspend, err := suspendComponent(g.clientSet.Workload(), nc.GraphdComponent(), oldSts)
	if err != nil {
		return fmt.Errorf("suspend graphd cluster %s failed: %v", componentName, err)
	}
	if needSuspend {
		klog.Infof("graphd cluster %s is suspended, skip reconciling", componentName)
		return nil
	}

	cm, cmHash, err := g.syncGraphdConfigMap(nc.DeepCopy())
	if err != nil {
		return err
	}

	metad, err := g.clientSet.NebulaMetad().GetNebulaMetad(nc.GetMetadNamespace(), nc.Spec.MetadRef.Name)
	if err != nil {
		return err
	}
	metadEndpoints := metad.MetadComponent().GetEndpoints(v1alpha1.MetadPortNameGRPC)
	newSts, err := nc.GraphdComponent().GenerateWorkload(cm, metadEndpoints)
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
		nc.Status.Graphd.Workload = &v1alpha1.WorkloadStatus{}
		return utilerrors.ReconcileErrorf("waiting for graphd cluster %s running", newSts.GetName())
	}

	if !nc.Status.Graphd.ServicesAdded {
		if err := addGraphdServices(metaClient, nc, *oldSts.Spec.Replicas, *newSts.Spec.Replicas); err != nil {
			return err
		}
		klog.Infof("graphd cluster [%s/%s] add services succeed", namespace, componentName)
		nc.Status.Graphd.ServicesAdded = true
	}

	if err := g.scaleManager.Scale(metaClient, nc, oldSts, newSts); err != nil {
		klog.Errorf("scale graphd cluster [%s/%s] failed: %v", namespace, componentName, err)
		return err
	}

	equal := podTemplateEqual(newSts, oldSts)
	if !equal || nc.Status.Graphd.Phase == v1alpha1.UpdatePhase {
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

func (g *graphdCluster) syncNebulaClusterStatus(nc *v1alpha1.NebulaCluster, newSts, oldSts *appsv1.StatefulSet) error {
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
		nc.Status.Graphd.Phase = v1alpha1.UpdatePhase
	} else if newReplicas < oldReplicas || (ok && newReplicas < lastReplicas) {
		nc.Status.Graphd.Phase = v1alpha1.ScaleInPhase
		if nc.Spec.Graphd.LogVolumeClaim != nil {
			if err := PVCMark(g.clientSet.PVC(), nc.GraphdComponent(), oldReplicas, newReplicas); err != nil {
				return err
			}
		}
	} else if newReplicas > oldReplicas || (ok && newReplicas > lastReplicas) {
		nc.Status.Graphd.Phase = v1alpha1.ScaleOutPhase
	} else {
		nc.Status.Graphd.Phase = v1alpha1.RunningPhase
	}

	return syncComponentStatus(nc.GraphdComponent(), &nc.Status.Graphd.ComponentStatus, oldSts)
}

func (g *graphdCluster) syncGraphdConfigMap(nc *v1alpha1.NebulaCluster) (*corev1.ConfigMap, string, error) {
	return syncConfigMap(
		nc.GraphdComponent(),
		g.clientSet.ConfigMap(),
		v1alpha1.GraphdConfigTemplate,
		nc.GraphdComponent().GetConfigMapKey())
}

func (g *graphdCluster) syncGraphdPVC(nc *v1alpha1.NebulaCluster) error {
	volumeStatus, err := syncPVC(nc.GraphdComponent(), g.clientSet.StorageClass(), g.clientSet.PVC())
	if err != nil {
		return err
	}
	nc.GraphdComponent().SetVolumeStatus(volumeStatus)
	return nil
}

func (g *graphdCluster) Delete(nc *v1alpha1.NebulaCluster) error {
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

func addGraphdServices(metaClient meta.Client, nc *v1alpha1.NebulaCluster, oldReplicas, newReplicas int32) error {
	var start int32
	if newReplicas > oldReplicas {
		start = oldReplicas
	}

	port := nc.GraphdComponent().GetPort(v1alpha1.GraphdPortNameGRPC)
	for i := start; i < newReplicas; i++ {
		host := nc.GraphdComponent().GetPodFQDN(i)
		hostReq := meta.NewAddHostReq(host, nc.Name, v1alpha1.AgentPort)
		if err := metaClient.AddHost(hostReq); err != nil {
			if ne, ok := err.(*nebula.NebulaError); ok {
				if ne.Code() != "NI104" {
					klog.Errorf("add host failed: %v", err)
					return err
				}
			} else {
				klog.Errorf("add host got unkonw error: %v", err)
				return err
			}
		}
		svcReq := meta.NewAddServiceReq(host, uint32(port), meta.ServiceTypeGraphd, nc.Name)
		if err := metaClient.AddService(svcReq); err != nil {
			if ne, ok := err.(*nebula.NebulaError); ok {
				if ne.Code() != "NI107" {
					klog.Errorf("add graphd service failed: %v", err)
					return err
				}
			} else {
				klog.Errorf("add graphd service got unkonw error: %v", err)
				return err
			}
		}
	}

	return nil
}
