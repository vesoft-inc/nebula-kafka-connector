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

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/label"
)

const (
	GraphdComponentType   ComponentType = "graphd"
	GraphdPortNameGRPC                  = "grpc"
	defaultGraphdPortGRPC               = 9669
	GraphdPortNameHTTP                  = "http"
	defaultGraphdPortHTTP               = 19669
	defaultGraphdImage                  = "vesoft/nebula-graphd"
)

var _ NebulaComponent = &graphdComponent{}

// +k8s:deepcopy-gen=false
func newGraphdComponent(nc *NebulaCluster) *graphdComponent {
	return &graphdComponent{
		cluster: cluster{
			nc:  nc,
			typ: GraphdComponentType,
		},
	}
}

type graphdComponent struct {
	cluster
}

func (c *graphdComponent) GetUpdateRevision() string {
	if c.nc.Status.Graphd.Workload == nil {
		return ""
	}
	return c.nc.Status.Graphd.Workload.UpdateRevision
}

func (c *graphdComponent) GetConfig() map[string]string {
	return c.nc.Spec.Graphd.Config
}

func (c *graphdComponent) GetConfigMapKey() string {
	return getCmKey(c.ComponentType().String())
}

func (c *graphdComponent) GetLogStorageClass() *string {
	if c.nc.Spec.Graphd.LogVolumeClaim == nil {
		return nil
	}
	scName := c.nc.Spec.Graphd.LogVolumeClaim.StorageClassName
	if scName == nil || *scName == "" {
		return nil
	}
	return scName
}

func (c *graphdComponent) GetLogStorageResources() *corev1.ResourceRequirements {
	if c.nc.Spec.Graphd.LogVolumeClaim == nil {
		return nil
	}
	return c.nc.Spec.Graphd.LogVolumeClaim.Resources.DeepCopy()
}

func (c *graphdComponent) GetDataStorageResources() (*corev1.ResourceRequirements, error) {
	return nil, nil
}

func (c *graphdComponent) GetServiceSpec() *ServiceSpec {
	if c.nc.Spec.Graphd.Service == nil {
		return nil
	}
	return c.nc.Spec.Graphd.Service.ServiceSpec.DeepCopy()
}

func (c *graphdComponent) GetHeadlessServiceName() string {
	return getServiceName(c.GetName(), true)
}

func (c *graphdComponent) GetServiceFQDN() string {
	return getServiceFQDN(c.GetHeadlessServiceName(), c.GetNamespace())
}

func (c *graphdComponent) GetPodFQDN(ordinal int32) string {
	return getPodFQDN(c.GetPodName(ordinal), c.GetServiceFQDN(), true)
}

func (c *graphdComponent) GetGrpcPort() int32 {
	return getPort(c.GenerateContainerPorts(), GraphdPortNameGRPC)
}

func (c *graphdComponent) GetPort(portName string) int32 {
	return getPort(c.GenerateContainerPorts(), portName)
}

func (c *graphdComponent) GetConnAddress(portName string) string {
	return joinHostPort(c.GetServiceFQDN(), c.GetPort(portName))
}

func (c *graphdComponent) GetEndpoints(portName string) []string {
	return getConnAddresses(
		c.GetConnAddress(portName),
		c.GetName(),
		c.ComponentSpec().Replicas())
}

func (c *graphdComponent) IsReady() bool {
	if c.nc.Status.Graphd.Workload == nil {
		return false
	}
	return *c.nc.Spec.Graphd.Replicas == c.nc.Status.Graphd.Workload.ReadyReplicas &&
		rollingUpdateDone(c.nc.Status.Graphd.Workload)
}

func (c *graphdComponent) GenerateLabels() map[string]string {
	return label.New().Cluster(c.GetClusterName()).Graphd()
}

func (c *graphdComponent) GenerateContainerPorts() []corev1.ContainerPort {
	grpcPort, err := parseCustomPort(defaultGraphdPortGRPC, c.GetConfig()["port"])
	if err != nil {
		return nil
	}

	httpPort, err := parseCustomPort(defaultGraphdPortHTTP, c.GetConfig()["http_port"])
	if err != nil {
		return nil
	}

	return []corev1.ContainerPort{
		{
			Name:          GraphdPortNameGRPC,
			ContainerPort: grpcPort,
		},
		{
			Name:          GraphdPortNameHTTP,
			ContainerPort: httpPort,
		},
	}
}

func (c *graphdComponent) GenerateVolumeMounts() []corev1.VolumeMount {
	componentType := c.ComponentType().String()
	mounts := make([]corev1.VolumeMount, 0)

	if c.nc.Spec.Graphd.LogVolumeClaim != nil {
		mounts = append(mounts, corev1.VolumeMount{
			Name:      logVolume(componentType),
			MountPath: "/usr/local/nebula/logs",
			SubPath:   "logs",
		})
	}

	return mounts
}

func (c *graphdComponent) GenerateVolumes() []corev1.Volume {
	componentType := c.ComponentType().String()
	volumes := make([]corev1.Volume, 0)

	if c.nc.Spec.Graphd.LogVolumeClaim != nil {
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

func (c *graphdComponent) GenerateVolumeClaim() ([]corev1.PersistentVolumeClaim, error) {
	if c.nc.Spec.Graphd.LogVolumeClaim == nil {
		return nil, nil
	}

	componentType := c.ComponentType().String()
	logSC, logRes := c.GetLogStorageClass(), c.GetLogStorageResources()
	storageRequest, err := parseStorageRequest(logRes.Requests)
	if err != nil {
		return nil, fmt.Errorf("cannot parse storage request for %s, error: %v", componentType, err)
	}

	claims := []corev1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: logVolume(componentType),
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources:        storageRequest,
				StorageClassName: logSC,
			},
		},
	}
	return claims, nil
}

func (c *graphdComponent) GenerateWorkload(cm *corev1.ConfigMap, metadEndpoints []string) (*appsv1.StatefulSet, error) {
	return generateWorkload(c, metadEndpoints, cm)
}

func (c *graphdComponent) GenerateService() *corev1.Service {
	return generateService(c, false)
}

func (c *graphdComponent) GenerateHeadlessService() *corev1.Service {
	return generateService(c, true)
}

func (c *graphdComponent) GenerateConfigMap() *corev1.ConfigMap {
	cm := generateConfigMap(c)
	configKey := getCmKey(c.ComponentType().String())
	cm.Data[configKey] = GraphdConfigTemplate
	return cm
}

func (c *graphdComponent) UpdateComponentStatus(status *ComponentStatus) {
	c.nc.Status.Graphd.ComponentStatus = *status
}

func (c *graphdComponent) SetWorkloadStatus(status *WorkloadStatus) {
	c.nc.Status.Graphd.Workload = status
}

func (c *graphdComponent) SetVolumeStatus(status *VolumeStatus) {
	c.nc.Status.Graphd.Volume = status
}
