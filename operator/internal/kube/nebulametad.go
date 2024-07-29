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

package kube

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v1alpha1"
)

type NebulaMetad interface {
	GetNebulaMetad(namespace, name string) (*v1alpha1.NebulaMetad, error)
	UpdateNebulaMetadStatus(nm *v1alpha1.NebulaMetad) error
}

type nebulaMetadClient struct {
	client client.Client
}

func NewNebulaMetad(client client.Client) NebulaMetad {
	return &nebulaMetadClient{client: client}
}

func (m *nebulaMetadClient) GetNebulaMetad(namespace, name string) (*v1alpha1.NebulaMetad, error) {
	nebulaMetad := &v1alpha1.NebulaMetad{}
	err := m.client.Get(context.TODO(), types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, nebulaMetad)
	if err != nil {
		klog.V(4).ErrorS(err, "failed to get NebulaMetad", "namespace", namespace, "name", name)
		return nil, err
	}
	return nebulaMetad, nil
}

func (m *nebulaMetadClient) UpdateNebulaMetadStatus(nm *v1alpha1.NebulaMetad) error {
	ns := nm.Namespace
	status := nm.Status.DeepCopy()

	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		nmClone, err := m.GetNebulaMetad(ns, nm.Name)
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("get NebulaMetad [%s/%s] failed: %v", ns, nm.Name, err))
			return err
		}

		if reflect.DeepEqual(*status, nmClone.Status) {
			return nil
		}

		nm = nmClone.DeepCopy()
		nm.Status = *status
		updateErr := m.client.Status().Update(context.TODO(), nm)
		if updateErr == nil {
			klog.Infof("NebulaMetad [%s/%s] status updated successfully", ns, nm.Name)
			return nil
		}
		klog.Errorf("update NebulaMetad [%s/%s] status failed: %v", ns, nm.Name, updateErr)
		return updateErr
	})
}
