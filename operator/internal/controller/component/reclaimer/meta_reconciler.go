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
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/label"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/controller/component"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/async"
)

const (
	metaReconcileConcurrency = 5
)

type meta struct {
	clientSet kube.ClientSet
}

func NewMetaReconciler(clientSet kube.ClientSet) component.MetaReconcileManager {
	return &meta{clientSet}
}

func (m *meta) Reconcile(obj runtime.Object, enablePVReclaim bool) error {
	metaObj, ok := obj.(metav1.Object)
	if !ok {
		return fmt.Errorf("%+v is not a runtime.Object", obj)
	}
	namespace := metaObj.GetNamespace()
	objName := metaObj.GetName()
	selector, err := label.New().Cluster(objName).Selector()
	if err != nil {
		return err
	}
	pods, err := m.clientSet.Pod().ListPods(namespace, selector)
	if err != nil {
		return fmt.Errorf("list pods for cluster %s/%s failed: %v", namespace, objName, err)
	}

	group := async.NewGroup(context.TODO(), metaReconcileConcurrency)
	for i := range pods {
		pod := pods[i]
		if !label.Label(pod.Labels).IsNebulaComponent() {
			continue
		}

		worker := func() error {
			var hasDataPV bool
			if pod.Labels[label.ComponentLabelKey] == label.GraphdLabelVal {
				hasDataPV = false
			}

			pvcs, err := m.resolvePVCFromPod(&pod)
			if err != nil {
				if errors.IsNotFound(err) && !hasDataPV {
					return nil
				}
				return err
			}
			for i := range pvcs {
				pvc := pvcs[i]
				if err := m.clientSet.PVC().UpdateMetaInfo(pvc, &pod, enablePVReclaim); err != nil {
					return err
				}
				if pvc.Spec.VolumeName == "" {
					continue
				}
				pv, err := m.clientSet.PV().GetPersistentVolume(pvc.Spec.VolumeName)
				if err != nil {
					klog.Errorf("cluster [%s/%s] get PV %s failed: %v", namespace, objName, pvc.Spec.VolumeName, err)
					return err
				}
				if err := m.clientSet.PV().UpdateMetaInfo(obj, pv); err != nil {
					return err
				}
			}
			return nil
		}

		group.Add(func(stopCh chan interface{}) {
			stopCh <- worker()
		})
	}

	return group.Wait()
}

func (m *meta) resolvePVCFromPod(pod *corev1.Pod) ([]*corev1.PersistentVolumeClaim, error) {
	var pvcs []*corev1.PersistentVolumeClaim
	var pvcName string
	for i := range pod.Spec.Volumes {
		vol := pod.Spec.Volumes[i]
		if vol.PersistentVolumeClaim == nil {
			continue
		}
		pvcName = vol.PersistentVolumeClaim.ClaimName
		if pvcName == "" {
			continue
		}
		pvc, err := m.clientSet.PVC().GetPVC(pod.Namespace, pvcName)
		if err != nil {
			klog.Errorf("pod [%s/%s] get PVC %s failed: %v", pod.Namespace, pod.Name, pvcName, err)
			continue
		}
		pvcs = append(pvcs, pvc)
	}
	if len(pvcs) == 0 {
		err := errors.NewNotFound(corev1.Resource("pvc"), "")
		err.ErrStatus.Message = fmt.Sprintf("no PVC found for pod [%s/%s]", pod.Namespace, pod.Name)
		return nil, err
	}
	return pvcs, nil
}
