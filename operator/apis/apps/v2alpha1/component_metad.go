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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/label"
)

const (
	MetadComponentType   ComponentType = "metad"
	MetadPortNameGRPC                  = "grpc"
	defaultMetadPortGRPC               = 9559
	MetadPortNameHTTP                  = "http"
	defaultMetadPortHTTP               = 19559
	defaultMetadImage                  = "vesoft/nebula-metad"
)

var _ NebulaComponent = &metadComponent{}

func newMetadComponent(metad *NebulaMetad) NebulaComponent {
	return &metadComponent{nm: metad, typ: MetadComponentType}
}

type metadComponent struct {
	nm  *NebulaMetad
	typ ComponentType
}

func (m *metadComponent) GetCluster() *NebulaCluster {
	return nil
}

func (m *metadComponent) GetNamespace() string {
	return m.nm.Namespace
}

func (m *metadComponent) GetName() string {
	return getComponentName(m.nm.Name, m.ComponentType())
}

func (m *metadComponent) GetPodName(ordinal int32) string {
	return getPodName(m.GetName(), ordinal)
}

func (m *metadComponent) GenerateOwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{
		{
			APIVersion:         m.nm.APIVersion,
			Kind:               m.nm.Kind,
			Name:               m.nm.GetName(),
			UID:                m.nm.GetUID(),
			Controller:         pointer.Bool(true),
			BlockOwnerDeletion: pointer.Bool(true),
		},
	}
}

func (m *metadComponent) ImagePullSecrets() []corev1.LocalObjectReference {
	return m.nm.Spec.ImagePullSecrets
}

func (m *metadComponent) ImagePullPolicy() *corev1.PullPolicy {
	return m.nm.Spec.ImagePullPolicy
}

func (m *metadComponent) ComponentType() ComponentType {
	return m.typ
}

func (m *metadComponent) ComponentSpec() ComponentAccessor {
	return buildComponentAccessor(nil, &m.nm.Spec.ComponentSpec)
}

func (m *metadComponent) GetLogStorageClass() *string {
	if m.nm.Spec.LogVolumeClaim == nil {
		return nil
	}
	scName := m.nm.Spec.LogVolumeClaim.StorageClassName
	if scName == nil || *scName == "" {
		return nil
	}
	return scName
}

func (m *metadComponent) GetDataStorageClass() *string {
	if m.nm.Spec.DataVolumeClaim == nil {
		return nil
	}
	scName := m.nm.Spec.DataVolumeClaim.StorageClassName
	if scName == nil || *scName == "" {
		return nil
	}
	return scName
}

func (m *metadComponent) GetLogStorageResources() *corev1.ResourceRequirements {
	if m.nm.Spec.LogVolumeClaim == nil {
		return nil
	}
	return m.nm.Spec.LogVolumeClaim.Resources.DeepCopy()
}

func (m *metadComponent) GetDataStorageResources() (*corev1.ResourceRequirements, error) {
	if m.nm.Spec.DataVolumeClaim == nil {
		return nil, nil
	}
	return m.nm.Spec.DataVolumeClaim.Resources.DeepCopy(), nil
}

func (m *metadComponent) GetConfig() map[string]string {
	return m.nm.Spec.Config
}

func (m *metadComponent) GetConfigMapKey() string {
	return getCmKey(m.ComponentType().String())
}

func (m *metadComponent) GetServiceSpec() *ServiceSpec {
	if m.nm.Spec.Service == nil {
		return nil
	}
	return m.nm.Spec.Service.DeepCopy()
}

func (m *metadComponent) GetHeadlessServiceName() string {
	return getServiceName(m.GetName(), true)
}

func (m *metadComponent) GetServiceFQDN() string {
	return getServiceFQDN(m.GetHeadlessServiceName(), m.nm.GetNamespace())
}

func (m *metadComponent) GetPodFQDN(ordinal int32) string {
	return getPodFQDN(m.GetPodName(ordinal), m.GetServiceFQDN(), true)
}

func (m *metadComponent) GetGrpcPort() int32 {
	return getPort(m.GenerateContainerPorts(), MetadPortNameGRPC)
}

func (m *metadComponent) GetPort(portName string) int32 {
	return getPort(m.GenerateContainerPorts(), portName)
}

func (m *metadComponent) GetConnAddress(portName string) string {
	return joinHostPort(m.GetServiceFQDN(), m.GetPort(portName))
}

func (m *metadComponent) GetEndpoints(portName string) []string {
	return getConnAddresses(
		m.GetConnAddress(portName),
		m.GetName(),
		m.ComponentSpec().Replicas())
}

func (m *metadComponent) GenerateLabels() map[string]string {
	return label.New().Cluster(m.nm.Name).Metad()
}

func (m *metadComponent) GenerateContainerPorts() []corev1.ContainerPort {
	grpcPort, err := parseCustomPort(defaultMetadPortGRPC, m.GetConfig()["port"])
	if err != nil {
		return nil
	}

	httpPort, err := parseCustomPort(defaultMetadPortHTTP, m.GetConfig()["http_port"])
	if err != nil {
		return nil
	}

	return []corev1.ContainerPort{
		{
			Name:          MetadPortNameGRPC,
			ContainerPort: grpcPort,
		},
		{
			Name:          MetadPortNameHTTP,
			ContainerPort: httpPort,
		},
	}
}

func (m *metadComponent) GenerateVolumeMounts() []corev1.VolumeMount {
	componentType := m.ComponentType().String()
	mounts := []corev1.VolumeMount{
		{
			Name:      dataVolume(componentType),
			MountPath: "/usr/local/nebula/data",
			SubPath:   "data",
		},
	}

	if m.nm.Spec.LogVolumeClaim != nil {
		mounts = append(mounts, corev1.VolumeMount{
			Name:      logVolume(componentType),
			MountPath: "/usr/local/nebula/logs",
			SubPath:   "logs",
		})
	}

	return mounts
}

func (m *metadComponent) GenerateVolumes() []corev1.Volume {
	componentType := m.ComponentType().String()
	volumes := []corev1.Volume{
		{
			Name: dataVolume(componentType),
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: dataVolume(componentType),
				},
			},
		},
	}

	if m.nm.Spec.LogVolumeClaim != nil {
		volumes = append(volumes, corev1.Volume{
			Name: logVolume(componentType),
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: logVolume(componentType),
				},
			},
		})
	}

	return volumes
}

func (m *metadComponent) GenerateVolumeClaim() ([]corev1.PersistentVolumeClaim, error) {
	componentType := m.ComponentType().String()
	claims := make([]corev1.PersistentVolumeClaim, 0)

	dataRes, err := m.GetDataStorageResources()
	if err != nil {
		return nil, err
	}
	dataSC := m.GetDataStorageClass()
	dataReq, err := parseStorageRequest(dataRes.Requests)
	if err != nil {
		return nil, fmt.Errorf("cannot parse storage request for %s data volume, error: %v", componentType, err)
	}

	claims = append(claims, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: dataVolume(componentType),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources:        dataReq,
			StorageClassName: dataSC,
		},
	})

	if m.nm.Spec.LogVolumeClaim != nil {
		logSC, logRes := m.GetLogStorageClass(), m.GetLogStorageResources()
		logReq, err := parseStorageRequest(logRes.Requests)
		if err != nil {
			return nil, fmt.Errorf("cannot parse storage request for %s log volume, error: %v", componentType, err)
		}

		claims = append(claims, corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name: logVolume(componentType),
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources:        logReq,
				StorageClassName: logSC,
			},
		})
	}

	return claims, nil
}

func (m *metadComponent) GenerateWorkload(cm *corev1.ConfigMap, _ []string) (*appsv1.StatefulSet, error) {
	metadEndpoints := m.GetEndpoints(MetadPortNameGRPC)
	return generateWorkload(m, metadEndpoints, cm)
}

func (m *metadComponent) GenerateService() *corev1.Service {
	return nil
}

func (m *metadComponent) GenerateHeadlessService() *corev1.Service {
	return generateService(m, true)
}

func (m *metadComponent) GenerateConfigMap() *corev1.ConfigMap {
	cm := generateConfigMap(m)
	configKey := getCmKey(m.ComponentType().String())
	cm.Data[configKey] = MetadhConfigTemplate
	return cm
}

func (m *metadComponent) IsReady() bool {
	if m.nm.Status.Workload == nil {
		return false
	}
	return *m.nm.Spec.Replicas == m.nm.Status.Workload.ReadyReplicas &&
		rollingUpdateDone(m.nm.Status.Workload)
}

func (m *metadComponent) GetUpdateRevision() string {
	if m.nm.Status.Workload == nil {
		return ""
	}
	return m.nm.Status.Workload.UpdateRevision
}

func (m *metadComponent) UpdateComponentStatus(status *ComponentStatus) {
	m.nm.Status.ComponentStatus = *status
}

func (m *metadComponent) SetVolumeStatus(status *VolumeStatus) {
	m.nm.Status.Volume = status
}
