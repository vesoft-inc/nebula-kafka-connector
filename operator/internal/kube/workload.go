/*
Copyright 2021 Vesoft Inc.

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

package kube

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Workload interface {
	GetWorkload(namespace string, name string) (*appsv1.StatefulSet, error)
	CreateWorkload(sts *appsv1.StatefulSet) error
	UpdateWorkload(sts *appsv1.StatefulSet) error
	DeleteWorkload(sts *appsv1.StatefulSet) error
}

type workloadClient struct {
	kubecli client.Client
}

func NewWorkload(kubecli client.Client) Workload {
	return &workloadClient{kubecli: kubecli}
}

func (w *workloadClient) GetWorkload(namespace, name string) (*appsv1.StatefulSet, error) {
	workload := &appsv1.StatefulSet{}
	err := w.kubecli.Get(context.TODO(), types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, workload)
	if err != nil {
		klog.V(4).ErrorS(err, "failed to get workload", "namespace", namespace, "name", name)
		return nil, err
	}
	return workload, nil
}

func (w *workloadClient) CreateWorkload(sts *appsv1.StatefulSet) error {
	if err := w.kubecli.Create(context.TODO(), sts); err != nil {
		if apierrors.IsAlreadyExists(err) {
			klog.Error(err, "workload already exists")
			return nil
		}
		return err
	}
	klog.Infof("workload %s/%s created successfully", sts.GetNamespace(), sts.GetName())
	return nil
}

func (w *workloadClient) UpdateWorkload(sts *appsv1.StatefulSet) error {
	ns := sts.GetNamespace()
	stsName := sts.GetName()
	spec := sts.Spec.DeepCopy()
	labels := sts.GetLabels()
	annotations := sts.GetAnnotations()

	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		updated, err := w.GetWorkload(ns, stsName)
		if err == nil {
			sts = updated.DeepCopy()
			sts.Spec = *spec
			sts.SetLabels(labels)
			sts.SetAnnotations(annotations)
		} else {
			utilruntime.HandleError(fmt.Errorf("get workload %s/%s failed: %v", ns, stsName, err))
			return err
		}

		updateErr := w.kubecli.Update(context.TODO(), sts)
		if updateErr == nil {
			klog.Infof("workload %s/%s updated successfully", ns, stsName)
			return nil
		}
		return updateErr
	})
}

func (w *workloadClient) DeleteWorkload(sts *appsv1.StatefulSet) error {
	preconditions := metav1.Preconditions{UID: &sts.UID, ResourceVersion: &sts.ResourceVersion}
	policy := metav1.DeletePropagationForeground
	options := &client.DeleteOptions{
		PropagationPolicy: &policy,
		Preconditions:     &preconditions,
	}
	return w.kubecli.Delete(context.TODO(), sts, options)
}
