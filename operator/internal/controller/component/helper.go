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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/annotation"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/label"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/config"
	utilerrors "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/hash"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/maputil"
)

func syncComponentStatus(
	component v2alpha1.NebulaComponent,
	status *v2alpha1.ComponentStatus,
	workload *appsv1.StatefulSet,
) error {
	if workload == nil {
		return nil
	}

	err := setWorkloadStatus(workload, status)
	if err != nil {
		return err
	}

	image := getContainerImage(workload, component.ComponentType().String())
	if image != "" && strings.Contains(image, ":") {
		status.Version = strings.Split(image, ":")[1]
	}

	component.UpdateComponentStatus(status)

	return nil
}

func setWorkloadStatus(sts *appsv1.StatefulSet, status *v2alpha1.ComponentStatus) error {
	workload := &v2alpha1.WorkloadStatus{}
	data, err := json.Marshal(sts.Status)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &workload); err != nil {
		return err
	}
	status.Workload = workload
	return nil
}

func getContainerImage(
	sts *appsv1.StatefulSet,
	containerName string,
) string {
	if sts == nil {
		return ""
	}
	containers := sts.Spec.Template.Spec.Containers
	for _, ctr := range containers {
		if ctr.Name == containerName {
			return ctr.Image
		}
	}
	return ""
}

func syncService(newSvc *corev1.Service, svcClient kube.Service) error {
	oldSvcTmp, err := svcClient.GetService(newSvc.Namespace, newSvc.Name)
	if apierrors.IsNotFound(err) {
		if err := setServiceLastAppliedConfigAnnotation(newSvc); err != nil {
			return err
		}
		return svcClient.CreateService(newSvc)
	}
	if err != nil {
		return err
	}

	oldSvc := oldSvcTmp.DeepCopy()
	equal, err := serviceEqual(newSvc, oldSvc)
	if err != nil {
		return err
	}

	annoEqual := maputil.IsSubMap(newSvc.Annotations, oldSvc.Annotations)
	isOrphan := metav1.GetControllerOf(oldSvc) == nil

	if !equal || !annoEqual || isOrphan {
		svc := *oldSvc
		svc.Spec = newSvc.Spec
		if err := setServiceLastAppliedConfigAnnotation(&svc); err != nil {
			return err
		}
		if oldSvc.Spec.ClusterIP != "" {
			svc.Spec.ClusterIP = oldSvc.Spec.ClusterIP
		}
		for k, v := range newSvc.Annotations {
			svc.Annotations[k] = v
		}
		if isOrphan {
			svc.OwnerReferences = newSvc.OwnerReferences
			svc.Labels = newSvc.Labels
		}
		if err := svcClient.UpdateService(&svc); err != nil {
			return err
		}
	}

	return nil
}

func setServiceLastAppliedConfigAnnotation(svc *corev1.Service) error {
	b, err := json.Marshal(svc.Spec)
	if err != nil {
		return err
	}
	if svc.Annotations == nil {
		svc.Annotations = map[string]string{}
	}
	svc.Annotations[annotation.AnnLastAppliedConfigKey] = string(b)
	return nil
}

func serviceEqual(newSvc, oldSvc *corev1.Service) (bool, error) {
	oldSpec := corev1.ServiceSpec{}
	if lastAppliedConfig, ok := oldSvc.Annotations[annotation.AnnLastAppliedConfigKey]; ok {
		err := json.Unmarshal([]byte(lastAppliedConfig), &oldSpec)
		if err != nil {
			return false, err
		}
		return apiequality.Semantic.DeepEqual(oldSpec, newSvc.Spec), nil
	}
	return false, nil
}

func setLastConfig(oldSts, newSts *appsv1.StatefulSet) error {
	spec := &appsv1.StatefulSetSpec{}
	if lastAppliedConfig, ok := oldSts.GetAnnotations()[annotation.AnnLastAppliedConfigKey]; ok {
		if err := json.Unmarshal([]byte(lastAppliedConfig), &spec); err != nil {
			return err
		}
	}
	newSts.Spec.Template.Spec = spec.Template.Spec
	return nil
}

func podTemplateEqual(newSts *appsv1.StatefulSet, oldSts *appsv1.StatefulSet) bool {
	oldStsSpec := appsv1.StatefulSetSpec{}
	lastAppliedConfig, ok := oldSts.Annotations[annotation.AnnLastAppliedConfigKey]
	if ok {
		err := json.Unmarshal([]byte(lastAppliedConfig), &oldStsSpec)
		if err != nil {
			klog.Errorf("unmarshal PodTemplate: [%s/%s]'s applied config failed,error: %v", oldSts.GetNamespace(), oldSts.GetName(), err)
			return false
		}
		return apiequality.Semantic.DeepEqual(oldStsSpec.Template.Spec, newSts.Spec.Template.Spec)
	}
	return false
}

func setPartition(sts *appsv1.StatefulSet, upgradeOrdinal int32) {
	sts.Spec.UpdateStrategy.RollingUpdate = &appsv1.RollingUpdateStatefulSetStrategy{Partition: &upgradeOrdinal}
	klog.Infof("sts [%s/%s] partition to %d", sts.GetNamespace(), sts.GetName(), upgradeOrdinal)
}

func getNextUpdatePod(component v2alpha1.NebulaComponent, replicas int32, podClient kube.Pod) (int32, error) {
	namespace := component.GetNamespace()
	updateRevision := component.GetUpdateRevision()
	for index := replicas - 1; index >= 0; index-- {
		podName := component.GetPodName(index)
		pod, err := podClient.GetPod(namespace, podName)
		if err != nil {
			return -1, err
		}
		revision, exist := pod.Labels[appsv1.ControllerRevisionHashLabelKey]
		if !exist {
			return -1, &utilerrors.ReconcileError{Msg: fmt.Sprintf("rolling updated pod %s has no label: %s",
				podName, appsv1.ControllerRevisionHashLabelKey)}
		}
		if revision == updateRevision {
			if pod.Status.Phase != corev1.PodRunning {
				return -1, &utilerrors.ReconcileError{Msg: fmt.Sprintf("rolling updated pod %s is not running", podName)}
			}
			continue
		}

		return index, nil
	}
	return -1, nil
}

func isUpdating(component v2alpha1.NebulaComponent, podClient kube.Pod, sts *appsv1.StatefulSet) (bool, error) {
	if statefulSetIsUpdating(sts) {
		return true, nil
	}

	selector, err := label.Label(component.GenerateLabels()).Selector()
	if err != nil {
		return false, err
	}

	pods, err := podClient.ListPods(component.GetNamespace(), selector)
	if err != nil {
		return false, fmt.Errorf(
			"failed to get pods for component [%s/%s], selector %s, error: %s",
			component.GetNamespace(),
			component.GetName(),
			selector,
			err,
		)
	}
	for i := range pods {
		pod := pods[i]
		revisionHash, exist := pod.Labels[appsv1.ControllerRevisionHashLabelKey]
		if !exist {
			return false, nil
		}
		if component.GetUpdateRevision() != "" &&
			revisionHash != component.GetUpdateRevision() {
			return true, nil
		}
	}
	return false, nil
}

func syncConfigMap(
	component v2alpha1.NebulaComponent,
	cmClient kube.ConfigMap,
	template,
	cmKey string,
) (*corev1.ConfigMap, string, error) {
	cmHash := hash.Hash(template)
	cm := component.GenerateConfigMap()
	cfg := component.GetConfig()
	if cfg != nil {
		customConf := config.AppendCustomConfig(template, component.GetConfig())
		cm.Data[cmKey] = customConf
		cmHash = hash.Hash(customConf)
	}

	if err := cmClient.CreateOrUpdateConfigMap(cm); err != nil {
		return nil, "", err
	}
	return cm, cmHash, nil
}

func syncPVC(
	component v2alpha1.NebulaComponent,
	pvcClient kube.PersistentVolumeClaim) error {
	replicas := int(component.ComponentSpec().Replicas())
	volumeClaims, err := component.GenerateVolumeClaim()
	if err != nil {
		return err
	}
	for _, volumeClaim := range volumeClaims {
		for i := 0; i < replicas; i++ {
			pvcName := fmt.Sprintf("%s-%s-%d", volumeClaim.Name, component.GetName(), i)
			oldPVC, err := pvcClient.GetPVC(component.GetNamespace(), pvcName)
			if err != nil {
				if !apierrors.IsNotFound(err) {
					return err
				}
			}
			if oldPVC == nil {
				continue
			}
			if volumeClaim.Spec.Resources.Requests.Storage().Cmp(*oldPVC.Spec.Resources.Requests.Storage()) == 1 {
				klog.Infof("expand PVC %s size to %s", pvcName, volumeClaim.Spec.Resources.Requests.Storage().String())
				// only update storage
				oldPVC.Spec.Resources.Requests = volumeClaim.Spec.Resources.Requests
				if err = pvcClient.UpdatePVC(oldPVC); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func setLastAppliedConfigAnnotation(sts *appsv1.StatefulSet) error {
	b, err := json.Marshal(sts.Spec)
	if err != nil {
		return err
	}
	if sts.Annotations == nil {
		sts.Annotations = map[string]string{}
	}
	sts.Annotations[annotation.AnnLastAppliedConfigKey] = string(b)
	return nil
}

func setTemplateAnnotations(sts *appsv1.StatefulSet, ann map[string]string) error {
	annotations := sts.Spec.Template.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string, len(ann))
	}

	for k, v := range ann {
		annotations[k] = v
	}
	sts.Spec.Template.SetAnnotations(annotations)
	return nil
}

func setLastReplicasAnnotation(sts *appsv1.StatefulSet) error {
	var lastReplicas int32
	val, ok := sts.GetAnnotations()[annotation.AnnLastReplicas]
	if ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		lastReplicas = int32(v)
	}
	replicas := pointer.Int32Deref(sts.Spec.Replicas, 0)
	if replicas == lastReplicas {
		return nil
	}
	annotations := make(map[string]string)
	for k, v := range sts.GetAnnotations() {
		annotations[k] = v
	}
	annotations[annotation.AnnLastReplicas] = strconv.Itoa(int(replicas))
	sts.SetAnnotations(annotations)
	return nil
}

func updateWorkload(workloadClient kube.Workload, newSts, oldSts *appsv1.StatefulSet) error {
	isOrphan := metav1.GetControllerOf(oldSts) == nil

	if newSts.Annotations == nil {
		newSts.Annotations = map[string]string{}
	}
	if oldSts.Annotations == nil {
		oldSts.Annotations = map[string]string{}
	}

	stsEqual := statefulSetEqual(*newSts, *oldSts)
	if stsEqual && !isOrphan {
		return nil
	}

	sts := oldSts
	annotations := annotation.CopyAnnotations(newSts.GetAnnotations())
	v, ok := oldSts.GetAnnotations()[annotation.AnnLastSyncTimestampKey]
	if ok {
		annotations[annotation.AnnLastSyncTimestampKey] = v
	}
	r, ok := oldSts.GetAnnotations()[annotation.AnnLastReplicas]
	if ok {
		annotations[annotation.AnnLastReplicas] = r
	}
	sts.SetAnnotations(annotations)
	sts.SetLabels(newSts.Labels)
	sts.Spec.Replicas = newSts.Spec.Replicas
	sts.Spec.UpdateStrategy = newSts.Spec.UpdateStrategy
	sts.Spec.Template = *newSts.Spec.Template.DeepCopy()
	if isOrphan {
		sts.OwnerReferences = newSts.OwnerReferences
	}

	if err := setLastReplicasAnnotation(sts); err != nil {
		return err
	}

	return workloadClient.UpdateWorkload(sts)
}

func statefulSetIsUpdating(sts *appsv1.StatefulSet) bool {
	if sts.Status.CurrentRevision != sts.Status.UpdateRevision {
		return true
	}
	if sts.Generation > sts.Status.ObservedGeneration && *sts.Spec.Replicas == sts.Status.Replicas {
		return true
	}
	return false
}

func statefulSetEqual(newSts appsv1.StatefulSet, oldSts appsv1.StatefulSet) bool {
	annotations := map[string]string{}
	for k, v := range oldSts.Annotations {
		if k == annotation.AnnLastAppliedConfigKey ||
			k == annotation.AnnLastReplicas {
			continue
		}
		annotations[k] = v
	}
	if !apiequality.Semantic.DeepEqual(newSts.Annotations, annotations) {
		return false
	}
	oldSpec := appsv1.StatefulSetSpec{}
	if lastAppliedConfig, ok := oldSts.Annotations[annotation.AnnLastAppliedConfigKey]; ok {
		err := json.Unmarshal([]byte(lastAppliedConfig), &oldSpec)
		if err != nil {
			klog.Errorf("unmarshal failed: %v", oldSts.GetNamespace(), oldSts.GetName(), err)
			return false
		}
		tmpTemplate := oldSpec.Template.DeepCopy()
		templateEqual := apiequality.Semantic.DeepEqual(*tmpTemplate, newSts.Spec.Template)
		return apiequality.Semantic.DeepEqual(oldSpec.Replicas, newSts.Spec.Replicas) &&
			apiequality.Semantic.DeepEqual(oldSpec.UpdateStrategy, newSts.Spec.UpdateStrategy) &&
			templateEqual
	}
	return false
}
