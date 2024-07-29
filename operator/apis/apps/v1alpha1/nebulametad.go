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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func (nm *NebulaMetad) MetadComponent() NebulaComponent {
	return newMetadComponent(nm)
}

func (nm *NebulaMetad) GetMetadEndpoints(portName string) []string {
	return nm.MetadComponent().GetEndpoints(portName)
}

func (nm *NebulaMetad) GetMetadThriftConnAddress() string {
	return nm.MetadComponent().GetConnAddress(MetadPortNameGRPC)
}

func (nm *NebulaMetad) IsSuspendEnabled() bool {
	return false
}

func (nm *NebulaMetad) IsPVReclaimEnabled() bool {
	return pointer.BoolDeref(nm.Spec.EnablePVReclaim, false)
}

func (nm *NebulaMetad) IsReady() bool {
	return nm.Status.ObservedGeneration == nm.Generation && nm.IsConditionReady()
}

func (nm *NebulaMetad) IsConditionReady() bool {
	for _, condition := range nm.Status.Conditions {
		if condition.Type == NebulaMetadReady {
			return condition.Status == metav1.ConditionTrue
		}
	}
	return false
}
