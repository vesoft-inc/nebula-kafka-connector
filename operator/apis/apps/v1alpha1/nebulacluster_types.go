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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// NebulaClusterReady indicates that the nebula cluster is ready or not.
	// This is defined as:
	// - All workloads are up-to-date (currentRevision == updateRevision).
	// - All nebula component pods are healthy.
	NebulaClusterReady = "Ready"
)

// NebulaClusterSpec defines the desired state of NebulaCluster
type NebulaClusterSpec struct {
	Graphd *GraphdSpec `json:"graphd"`

	Storaged *StoragedSpec `json:"storaged"`

	MetadRef *NamespacedObjectReference `json:"metadRef"`

	// +kubebuilder:default=1
	ReplicaFactor int32 `json:"replicaFactor,omitempty"`

	// tenant cluster owner
	// +kubebuilder:default=root
	Owner string `json:"owner"`

	Zones []string `json:"zones,omitempty"`

	// CredentialSecret contains user login info.
	CredentialSecret string `json:"credentialSecret,omitempty"`

	// +kubebuilder:default=default-scheduler
	// +optional
	SchedulerName string `json:"schedulerName"`

	// TopologySpreadConstraints specifies how to spread matching pods among the given topology.
	TopologySpreadConstraints []TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`

	// +optional
	Suspend *bool `json:"suspend,omitempty"`

	// Flag to enable/disable PV reclaim while the nebula cluster deleted, default false
	// +optional
	EnablePVReclaim *bool `json:"enablePVReclaim,omitempty"`

	// +kubebuilder:default=Always
	ImagePullPolicy *corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// +optional
	Console *ConsoleSpec `json:"console,omitempty"`

	// +optional
	Agent *AgentContainerSpec `json:"agent,omitempty"`

	// +optional
	AlpineImage *string `json:"alpineImage,omitempty"`

	// +optional
	Exporter *ExporterSpec `json:"exporter,omitempty"`
}

// ExporterSpec defines the desired state of Exporter
type ExporterSpec struct {
	ComponentSpec `json:",inline"`

	// Maximum number of parallel scrape requests
	// +kubebuilder:default=40
	// +optional
	MaxRequests int32 `json:"maxRequests,omitempty"`

	// CollectRegex means the regex to filter metrics
	// +optional
	CollectRegex string `json:"collectRegex,omitempty"`

	// IgnoreRegex means the regex to ignore metrics
	// +optional
	IgnoreRegex string `json:"ignoreRegex,omitempty"`

	// +optional
	// +kubebuilder:default=9100
	HTTPPort int32 `json:"httpPort,omitempty"`
}

type ConsoleSpec struct {
	// Container image.
	// +optional
	Image string `json:"image,omitempty"`

	// Version tag for container image.
	// +optional
	Version string `json:"version,omitempty"`

	// K8S nodeSelector.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

// AgentContainerSpec defines the desired state of Agent
type AgentContainerSpec struct {
	// Container image.
	// +optional
	Image string `json:"image,omitempty"`

	// Version tag for container image.
	// +kubebuilder:default=latest
	// +optional
	Version string `json:"version,omitempty"`

	// K8S resources settings.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Container environment variables.
	// +optional
	EnvVars []corev1.EnvVar `json:"env,omitempty"`

	// +optional
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
}

type GraphdSpec struct {
	ComponentSpec `json:",inline"`

	// GraphProjection defines a main-memory volatile projection
	// of a durable data-source graph from the storage.
	// +optional
	GraphProjection *GraphProjection `json:"graphProjection,omitempty"`

	// Config defines a graphd configuration load into ConfigMap
	Config map[string]string `json:"config,omitempty"`

	// Service defines a k8s service of Graphd cluster.
	// +optional
	Service *GraphdServiceSpec `json:"service,omitempty"`

	// K8S persistent volume claim for Graphd log volume.
	// +optional
	LogVolumeClaim *StorageClaim `json:"logVolumeClaim,omitempty"`
}

type StoragedSpec struct {
	ComponentSpec `json:",inline"`

	// Config defines a storaged configuration load into ConfigMap
	Config map[string]string `json:"config,omitempty"`

	// Service defines a Kubernetes service of Storaged cluster.
	// +optional
	Service *ServiceSpec `json:"service,omitempty"`

	// K8S persistent volume claim for Storaged log volume.
	// +optional
	LogVolumeClaim *StorageClaim `json:"logVolumeClaim,omitempty"`

	// K8S persistent volume claim for Storaged data volume.
	// +optional
	DataVolumeClaims []StorageClaim `json:"dataVolumeClaims,omitempty"`

	// Flag to enable/disable auto balance data and leader while the nebula storaged scale out, default false
	// +optional
	EnableAutoBalance *bool `json:"enableAutoBalance,omitempty"`
}

// StorageClaim contains details of storage
type StorageClaim struct {
	// Resources represents the minimum resources the volume should have.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Name of the StorageClass required by the claim.
	// +optional
	StorageClassName *string `json:"storageClassName,omitempty"`
}

// GraphdServiceSpec is the service spec of graphd
type GraphdServiceSpec struct {
	ServiceSpec `json:",inline"`

	// LoadBalancerIP is the loadBalancerIP of service
	// +optional
	LoadBalancerIP *string `json:"loadBalancerIP,omitempty"`

	// ExternalTrafficPolicy of the service
	// +optional
	ExternalTrafficPolicy *corev1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`
}

// NebulaClusterStatus defines the observed state of NebulaCluster
type NebulaClusterStatus struct {
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
	Graphd             GraphdStatus       `json:"graphd,omitempty"`
	Storaged           StoragedStatus     `json:"storaged,omitempty"`
	CreatedDone        bool               `json:"createdDone,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
	Version            string             `json:"version,omitempty"`
}

// StoragedStatus describes the status and version of nebula storaged.
type StoragedStatus struct {
	ComponentStatus   `json:",inline"`
	ServicesAdded     bool          `json:"servicesAdded,omitempty"`
	InitPartitionDone bool          `json:"initPartitionDone,omitempty"`
	Volume            *VolumeStatus `json:"volume,omitempty"`
}

// GraphdStatus describes the status and version of nebula graphd.
type GraphdStatus struct {
	ComponentStatus `json:",inline"`
	ServicesAdded   bool          `json:"servicesAdded,omitempty"`
	Volume          *VolumeStatus `json:"volume,omitempty"`
}

// ComponentStatus is the status and version of a nebula component.
type ComponentStatus struct {
	Version  string          `json:"version,omitempty"`
	Phase    ComponentPhase  `json:"phase,omitempty"`
	Workload *WorkloadStatus `json:"workload,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=ngnc
// +kubebuilder:printcolumn:name="READY",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="GRAPHD-DESIRED",type="string",JSONPath=".spec.graphd.replicas",description="The desired number of graphd pods."
// +kubebuilder:printcolumn:name="GRAPHD-READY",type="string",JSONPath=".status.graphd.workload.readyReplicas",description="The number of graphd pods ready."
// +kubebuilder:printcolumn:name="STORAGED-DESIRED",type="string",JSONPath=".spec.storaged.replicas",description="The desired number of storaged pods."
// +kubebuilder:printcolumn:name="STORAGED-READY",type="string",JSONPath=".status.storaged.workload.readyReplicas",description="The number of storaged pods ready."
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp",description="CreationTimestamp is a timestamp representing the server time when this object was created. It is represented in RFC3339 form and is in UTC."

// NebulaCluster is the Schema for the nebulaclusters API
type NebulaCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NebulaClusterSpec   `json:"spec,omitempty"`
	Status NebulaClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NebulaClusterList contains a list of NebulaCluster
type NebulaClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NebulaCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NebulaCluster{}, &NebulaClusterList{})
}
