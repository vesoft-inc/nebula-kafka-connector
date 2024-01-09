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

package reclaimer

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/annotation"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/label"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
)

type PVCReclaimer interface {
	Reclaim(obj runtime.Object) error
}

type pvcReclaimer struct {
	clientSet kube.ClientSet
}

func NewPVCReclaimer(clientSet kube.ClientSet) PVCReclaimer {
	return &pvcReclaimer{clientSet: clientSet}
}

func (p *pvcReclaimer) Reclaim(obj runtime.Object) error {
	return p.reclaimPV(obj)
}

func (p *pvcReclaimer) reclaimPV(obj runtime.Object) error {
	metaObj, ok := obj.(metav1.Object)
	if !ok {
		return fmt.Errorf("%+v is not a runtime.Object", obj)
	}
	namespace := metaObj.GetNamespace()
	objName := metaObj.GetName()

	pvcs, err := p.listPVCs(namespace, objName)
	if err != nil {
		return err
	}

	for i := range pvcs {
		pvc := pvcs[i]
		pvcName := pvc.GetName()
		if !label.Label(pvc.Labels).IsNebulaComponent() {
			klog.V(4).Infof("skip reclaim for PVC %s is not associate with nebula component", pvcName)
			continue
		}

		if pvc.Status.Phase != corev1.ClaimBound {
			klog.V(4).Infof("skip reclaim for PVC %s %s status is not bound", pvcName)
			continue
		}

		if pvc.DeletionTimestamp != nil {
			klog.V(4).Infof("skip reclaim for PVC %s has been deleted", pvcName)
			continue
		}

		if pvc.Annotations[annotation.AnnPVCDeferDeletingKey] == "" {
			klog.V(4).Infof("skip reclaim for PVC %s has not been marked as defer deleting pvc", pvcName)
			continue
		}

		podName, exist := pvc.Annotations[annotation.AnnPodNameKey]
		if !exist {
			klog.V(4).Infof("skip reclaim for PVC %s has no pod name annotation", pvcName)
			continue
		}

		_, err := p.clientSet.Pod().GetPod(namespace, podName)
		if err == nil {
			klog.V(4).Infof("skip reclaim for PVC %s is still referenced by a pod", pvcName)
			continue
		}
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("cluster [%s/%s] get PVC %s pod %s from cache failed: %v", namespace, objName, pvcName, podName, err)
		}

		pvName := pvc.Spec.VolumeName
		pv, err := p.clientSet.PV().GetPersistentVolume(pvName)
		if err != nil {
			if apierrors.IsNotFound(err) {
				continue
			}
			return fmt.Errorf("cluster [%s/%s] get PVC %s PV %s failed: %v", namespace, objName, pvcName, pvName, err)
		}

		if pv.Spec.PersistentVolumeReclaimPolicy != corev1.PersistentVolumeReclaimDelete {
			if err := p.clientSet.PV().PatchPVReclaimPolicy(pv, corev1.PersistentVolumeReclaimDelete); err != nil {
				return fmt.Errorf("cluster [%s/%s] patch PV %s to %s failed: %v", namespace, objName, pvName,
					corev1.PersistentVolumeReclaimDelete, err)
			}
			klog.Infof("patch PV %s policy to Delete successfully", pvName)
		}

		if err := p.clientSet.PVC().DeletePVC(pvc.Namespace, pvcName); err != nil {
			if !apierrors.IsNotFound(err) {
				klog.Errorf("cluster [%s/%s] delete PVC %s failed: %v", namespace, objName, pvcName, err)
				return err
			}
		}
		klog.Infof("cluster [%s/%s] reclaim PV %s successfully", namespace, objName, pvName)
	}
	return nil
}

func (p *pvcReclaimer) listPVCs(namespace, objName string) ([]corev1.PersistentVolumeClaim, error) {
	selector, err := label.New().Cluster(objName).Selector()
	if err != nil {
		return nil, fmt.Errorf("get cluster [%s/%s] label selector failed: %v", namespace, objName, err)
	}

	pvcs, err := p.clientSet.PVC().ListPVCs(namespace, selector)
	if err != nil {
		return nil, fmt.Errorf("cluster [%s/%s] list PVC failed: %v", namespace, objName, err)
	}
	return pvcs, nil
}
