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

package v2alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// NebulaMetadReady indicates that the nebula metad is ready or not.
	// This is defined as:
	// - All nebula metad pods are up-to-date (currentRevision == updateRevision).
	// - All nebula metad pods are healthy.
	NebulaMetadReady = "Ready"
)

// NebulaMetadSpec defines the desired state of NebulaMetad
type NebulaMetadSpec struct {
	ComponentSpec `json:",inline"`

	// +kubebuilder:default=Always
	ImagePullPolicy *corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// Flag to enable/disable PV reclaim while the nebula cluster deleted, default false
	// +optional
	EnablePVReclaim *bool `json:"enablePVReclaim,omitempty"`

	// Config defines a metad configuration load into ConfigMap
	Config map[string]string `json:"config,omitempty"`

	// Service defines a Kubernetes service of Metad cluster.
	// +optional
	Service *ServiceSpec `json:"service,omitempty"`

	// K8S persistent volume claim for Metad log volume.
	// +optional
	LogVolumeClaim *StorageClaim `json:"logVolumeClaim,omitempty"`

	// K8S persistent volume claim for Metad data volume.
	// +optional
	DataVolumeClaim *StorageClaim `json:"dataVolumeClaim,omitempty"`
}

// NebulaMetadStatus defines the observed state of NebulaMetad
type NebulaMetadStatus struct {
	ComponentStatus    `json:",inline"`
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
	ManagedClusters    int32              `json:"managedClusters,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=nm
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="METAD-DESIRED",type="string",JSONPath=".spec.replicas",description="The desired number of metad pods."
// +kubebuilder:printcolumn:name="METAD-READY",type="string",JSONPath=".status.workload.readyReplicas",description="The number of metad pods ready."
// +kubebuilder:printcolumn:name="MANAGED-CLUSTERS",type="string",JSONPath=".status.managedClusters",description="The managed clusters of this metad."
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp",description="CreationTimestamp is a timestamp representing the server time when this object was created. It is represented in RFC3339 form and is in UTC."

// NebulaMetad is the Schema for the nebulametads API
type NebulaMetad struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NebulaMetadSpec   `json:"spec,omitempty"`
	Status NebulaMetadStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NebulaMetadList contains a list of NebulaMetad
type NebulaMetadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NebulaMetad `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NebulaMetad{}, &NebulaMetadList{})
}
