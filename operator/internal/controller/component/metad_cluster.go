/*
Copyright 2024.

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
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/annotation"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	utilerrors "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
)

type MetadReconcileManager interface {
	// Reconcile reconciles the metad desired state
	Reconcile(metad *v2alpha1.NebulaMetad) error

	// Delete deletes the metad
	Delete(metad *v2alpha1.NebulaMetad) error
}

type metadCluster struct {
	clientSet     kube.ClientSet
	updateManager MetadUpdateManager
	eventRecorder record.EventRecorder
}

func NewMetadManager(
	clientSet kube.ClientSet,
	updateManager MetadUpdateManager,
	recorder record.EventRecorder,
) MetadReconcileManager {
	return &metadCluster{
		clientSet:     clientSet,
		updateManager: updateManager,
		eventRecorder: recorder,
	}
}

func (c *metadCluster) Reconcile(nm *v2alpha1.NebulaMetad) error {
	if err := c.syncMetadHeadlessService(nm); err != nil {
		return err
	}

	return c.syncMetadWorkload(nm)
}

func (c *metadCluster) Delete(nm *v2alpha1.NebulaMetad) error {
	namespace := nm.GetNamespace()
	componentName := nm.MetadComponent().GetName()

	workload, err := c.clientSet.Workload().GetWorkload(namespace, componentName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		klog.Errorf("get metad cluster failed: %v", err)
		return err
	}
	return c.clientSet.Workload().DeleteWorkload(workload)
}

func (c *metadCluster) syncMetadHeadlessService(nm *v2alpha1.NebulaMetad) error {
	newSvc := nm.MetadComponent().GenerateHeadlessService()
	if newSvc == nil {
		return nil
	}

	return syncService(newSvc, c.clientSet.Service())
}

func (c *metadCluster) syncMetadWorkload(nm *v2alpha1.NebulaMetad) error {
	namespace := nm.GetNamespace()
	componentName := nm.MetadComponent().GetName()

	oldWorkloadTemp, err := c.clientSet.Workload().GetWorkload(namespace, componentName)
	if err != nil && !apierrors.IsNotFound(err) {
		klog.Errorf("get metad cluster failed: %v", err)
		return err
	}

	notExist := apierrors.IsNotFound(err)
	oldWorkload := oldWorkloadTemp.DeepCopy()

	cm, cmHash, err := c.syncMetadConfigMap(nm.DeepCopy())
	if err != nil {
		return err
	}

	newWorkload, err := nm.MetadComponent().GenerateWorkload(cm, nil)
	if err != nil {
		klog.Errorf("generate metad cluster template failed: %v", err)
		return err
	}

	if err := setTemplateAnnotations(
		newWorkload,
		map[string]string{annotation.AnnPodConfigMapHash: cmHash}); err != nil {
		return err
	}

	if err = c.syncNebulaMetadStatus(nm, oldWorkload); err != nil {
		return fmt.Errorf("sync metad cluster status failed: %v", err)
	}

	if notExist {
		if err := setLastAppliedConfigAnnotation(newWorkload); err != nil {
			return err
		}
		if err := c.clientSet.Workload().CreateWorkload(newWorkload); err != nil {
			return err
		}
		nm.Status.Workload = &v2alpha1.WorkloadStatus{}
		return utilerrors.ReconcileErrorf("waiting for metad cluster %s running", newWorkload.GetName())
	}

	equal := podTemplateEqual(newWorkload, oldWorkload)
	if !equal || nm.Status.Phase == v2alpha1.UpdatePhase {
		if err := c.updateManager.Update(nm, oldWorkload, newWorkload); err != nil {
			return err
		}
	}

	if err := c.syncMetadPVC(nm); err != nil {
		return err
	}

	if equal && nm.MetadComponent().IsReady() {
		username, password, err := kube.GetCredential(c.clientSet, nm.Namespace, nm.Spec.CredentialSecret)
		if err != nil {
			return err
		}
		endpoints := []string{nm.GetMetadThriftConnAddress()}
		metaClient, err := meta.NewMetaClient(strings.Join(endpoints, ","), meta.WithUserPassword(username, password))
		if err != nil {
			return err
		}
		if _, err := metaClient.Login(); err != nil {
			klog.Errorf("login metad failed: %v", err)
			return err
		}
		defer func() {
			metaClient.Close()
		}()

		if err := c.syncManagedClusters(metaClient, nm); err != nil {
			return err
		}
		if err := c.setVersion(metaClient, nm); err != nil {
			return err
		}
	}

	return updateWorkload(c.clientSet.Workload(), newWorkload, oldWorkload)
}

func (c *metadCluster) syncMetadConfigMap(nm *v2alpha1.NebulaMetad) (*corev1.ConfigMap, string, error) {
	return syncConfigMap(
		nm.MetadComponent(),
		c.clientSet.ConfigMap(),
		v2alpha1.MetadhConfigTemplate,
		nm.MetadComponent().GetConfigMapKey())
}

func (c *metadCluster) syncNebulaMetadStatus(nm *v2alpha1.NebulaMetad, oldWorkload *appsv1.StatefulSet) error {
	if oldWorkload == nil {
		return nil
	}
	updating, err := isUpdating(nm.MetadComponent(), c.clientSet.Pod(), oldWorkload)
	if err != nil {
		return err
	}

	if updating {
		nm.Status.Phase = v2alpha1.UpdatePhase
	} else {
		nm.Status.Phase = v2alpha1.RunningPhase
	}

	return syncComponentStatus(nm.MetadComponent(), &nm.Status.ComponentStatus, oldWorkload)
}

func (c *metadCluster) syncManagedClusters(mc meta.Client, nm *v2alpha1.NebulaMetad) error {
	req := meta.NewListClustersReq("")
	resp, err := mc.ListClusters(req)
	if err != nil {
		return err
	}
	nm.Status.ManagedClusters = int32(len(resp.Clusters))
	return nil
}

// TODO set server version
func (c *metadCluster) setVersion(mc meta.Client, nm *v2alpha1.NebulaMetad) error {
	return nil
}

func (c *metadCluster) syncMetadPVC(nm *v2alpha1.NebulaMetad) error {
	volumeStatus, err := syncPVC(nm.MetadComponent(), c.clientSet.StorageClass(), c.clientSet.PVC())
	if err != nil {
		return nil
	}
	nm.MetadComponent().SetVolumeStatus(volumeStatus)
	return nil
}
