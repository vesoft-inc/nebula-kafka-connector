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

package v2alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func (nc *NebulaCluster) GraphdComponent() NebulaClusterComponent {
	return newGraphdComponent(nc)
}

func (nc *NebulaCluster) StoragedComponent() NebulaClusterComponent {
	return newStoragedComponent(nc)
}

func (nc *NebulaCluster) ComponentByType(typ ComponentType) (NebulaClusterComponent, error) {
	switch typ {
	case GraphdComponentType:
		return nc.GraphdComponent(), nil
	case StoragedComponentType:
		return nc.StoragedComponent(), nil
	}

	return nil, fmt.Errorf("unsupport component %s", typ)
}

func (nc *NebulaCluster) GetStoragedEndpoints(portName string) []string {
	return nc.StoragedComponent().GetEndpoints(portName)
}

func (nc *NebulaCluster) GetGraphdEndpoints(portName string) []string {
	return nc.GraphdComponent().GetEndpoints(portName)
}

func (nc *NebulaCluster) GetGraphdServiceName() string {
	return getServiceName(nc.GraphdComponent().GetName(), false)
}

func (nc *NebulaCluster) GetClusterName() string {
	return nc.Name
}

func (nc *NebulaCluster) GenerateOwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion:         nc.APIVersion,
			Kind:               nc.Kind,
			Name:               nc.GetName(),
			UID:                nc.GetUID(),
			Controller:         pointer.Bool(true),
			BlockOwnerDeletion: pointer.Bool(true),
		},
	}
}

func (nc *NebulaCluster) IsPVReclaimEnabled() bool {
	return pointer.BoolDeref(nc.Spec.EnablePVReclaim, false)
}

func (nc *NebulaCluster) IsAutoBalanceEnabled() bool {
	return pointer.BoolDeref(nc.Spec.Storaged.EnableAutoBalance, false)
}

func (nc *NebulaCluster) IsZoneEnabled() bool {
	return len(nc.Spec.Zones) > 0
}

func (nc *NebulaCluster) IsReady() bool {
	return nc.Status.ObservedGeneration == nc.Generation && nc.IsConditionReady()
}

func (nc *NebulaCluster) IsConditionReady() bool {
	for _, condition := range nc.Status.Conditions {
		if condition.Type == NebulaClusterReady {
			return condition.Status == metav1.ConditionTrue
		}
	}
	return false
}
