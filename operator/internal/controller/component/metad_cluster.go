package component

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/tools/record"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/annotation"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/nebula"
	utilerrors "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
)

type MetadReconcileManager interface {
	// Reconcile reconciles the  Metad metad desired state
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
	recorder record.EventRecorder,
	updateManager MetadUpdateManager,
) MetadReconcileManager {
	return &metadCluster{
		clientSet:     clientSet,
		eventRecorder: recorder,
		updateManager: updateManager,
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

	// TODO: suspend metad cluster
	//needSuspend, err := suspendComponent(c.clientSet.Workload(), nm.MetadComponent(), oldWorkload)
	//if err != nil {
	//	return fmt.Errorf("suspend metad cluster %s failed: %v", componentName, err)
	//}
	//if needSuspend {
	//	klog.Infof("metad cluster %s is suspended, skip reconciling", componentName)
	//	return nil
	//}

	cm, cmHash, err := c.syncMetadConfigMap(nm.DeepCopy())
	if err != nil {
		return err
	}

	newWorkload, err := nm.MetadComponent().GenerateWorkload(cm, []string{})
	if err != nil {
		klog.Errorf("generate metad cluster template failed: %v", err)
		return err
	}

	if err := setTemplateAnnotations(
		newWorkload,
		map[string]string{annotation.AnnPodConfigMapHash: cmHash}); err != nil {
		return err
	}

	if !notExist {
		timestamp, ok := oldWorkload.GetAnnotations()[annotation.AnnRestartTimestamp]
		if ok && timestamp != "" {
			if err := setTemplateAnnotations(newWorkload,
				map[string]string{annotation.AnnRestartTimestamp: timestamp}); err != nil {
				return err
			}
		}
	}

	if err = c.syncNebulaMetadStatus(nm, oldWorkload); err != nil {
		return fmt.Errorf("sync metad cluster status failed: %v", err)
	}

	if notExist {
		if err := setRestartTimestamp(newWorkload); err != nil {
			return err
		}

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
		if err := c.setVersion(nm); err != nil {
			return err
		}

		endpoints := nm.GetMetadEndpoints(v2alpha1.MetadPortNameHTTP)
		if err := updateDynamicFlags(endpoints, newWorkload.GetAnnotations()); err != nil {
			return fmt.Errorf("update metad cluster %s dynamic flags failed: %v", newWorkload.GetName(), err)
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

	// TODO: show metad hosts state with metad peers
	return syncComponentStatus(nm.MetadComponent(), &nm.Status.ComponentStatus, oldWorkload)
}

func (c *metadCluster) setVersion(nm *v2alpha1.NebulaMetad) error {
	endpoints := []string{nm.GetMetadThriftConnAddress()}
	metaClient, err := nebula.NewMetaClient(endpoints)
	if err != nil {
		return err
	}
	defer func() {
		_ = metaClient.Disconnect()
	}()

	// TODO: get metad version
	return nil
}

func (c *metadCluster) syncMetadPVC(nm *v2alpha1.NebulaMetad) error {
	return syncPVC(nm.MetadComponent(), c.clientSet.PVC())
}
