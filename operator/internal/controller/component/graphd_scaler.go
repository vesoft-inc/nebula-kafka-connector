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

package component

import (
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	utilerrors "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
)

type graphScaler struct {
	clientSet kube.ClientSet
}

func NewGraphScaler(clientSet kube.ClientSet) ScaleManager {
	return &graphScaler{clientSet: clientSet}
}

func (s *graphScaler) Scale(nc *v2alpha1.NebulaCluster, oldSts, newSts *appsv1.StatefulSet) error {
	oldReplicas := pointer.Int32Deref(oldSts.Spec.Replicas, 0)
	newReplicas := pointer.Int32Deref(newSts.Spec.Replicas, 0)

	if newReplicas < oldReplicas || nc.Status.Graphd.Phase == v2alpha1.ScaleInPhase {
		return s.ScaleIn(nc, oldReplicas, newReplicas)
	}

	if newReplicas > oldReplicas || nc.Status.Graphd.Phase == v2alpha1.ScaleOutPhase {
		return s.ScaleOut(nc, oldReplicas, newReplicas)
	}

	return nil
}

func (s *graphScaler) ScaleOut(nc *v2alpha1.NebulaCluster, oldReplicas, newReplicas int32) error {
	namespace := nc.GetNamespace()
	componentName := nc.GraphdComponent().GetName()
	metad, err := s.clientSet.NebulaMetad().GetNebulaMetad(nc.GetMetadNamespace(), nc.Spec.MetadRef.Name)
	if err != nil {
		return err
	}
	metadEndpoints := metad.MetadComponent().GetEndpoints(v2alpha1.MetadPortNameGRPC)
	if err := addGraphdServices(nc, metadEndpoints, oldReplicas, newReplicas); err != nil {
		klog.Errorf("add graphd services failed: %v", err)
		return err
	}
	klog.Infof("graphd cluster [%s/%s] add services succeed", namespace, componentName)
	return nil
}

func (s *graphScaler) ScaleIn(nc *v2alpha1.NebulaCluster, oldReplicas, newReplicas int32) error {
	namespace := nc.GetNamespace()
	componentName := nc.GraphdComponent().GetName()
	nc.Status.Graphd.Phase = v2alpha1.ScaleInPhase
	if err := s.clientSet.NebulaCluster().UpdateNebulaClusterStatus(nc); err != nil {
		return err
	}

	metad, err := s.clientSet.NebulaMetad().GetNebulaMetad(nc.GetMetadNamespace(), nc.Spec.MetadRef.Name)
	if err != nil {
		return err
	}
	metadEndpoints := metad.MetadComponent().GetEndpoints(v2alpha1.MetadPortNameGRPC)

	metaClient, err := meta.NewMetaClient(strings.Join(metadEndpoints, ","))
	if err != nil {
		return err
	}
	defer func() {
		metaClient.Close()
	}()

	port := nc.GraphdComponent().GetPort(v2alpha1.GraphdPortNameGRPC)
	if oldReplicas-newReplicas > 0 {
		for i := oldReplicas - 1; i >= newReplicas; i-- {
			host := nc.GraphdComponent().GetPodFQDN(i)
			req := meta.NewDropServiceReq(host, uint32(port), meta.ServiceTypeGraphd, nc.Name)
			resp, err := metaClient.DropService(req)
			if err != nil {
				return err
			}
			if !resp.OK {
				return fmt.Errorf("drop graphd service failed, code: %s, msg: %v", resp.GetErrorCode(), resp.GetErrorMsg())
			}
		}
	}
	klog.Infof("graphd cluster [%s/%s] drop services succeed", namespace, componentName)

	if err := PVCMark(s.clientSet.PVC(), nc.GraphdComponent(), oldReplicas, newReplicas); err != nil {
		return err
	}

	deleted := true
	pvcNames := ordinalPVCNames(nc.GraphdComponent().ComponentType(), nc.GraphdComponent().GetName(), newReplicas)
	for _, pvcName := range pvcNames {
		if _, err = s.clientSet.PVC().GetPVC(nc.GetNamespace(), pvcName); err != nil {
			if !apierrors.IsNotFound(err) {
				deleted = false
				break
			}
		}
	}
	if !deleted {
		return &utilerrors.ReconcileError{Msg: fmt.Sprintf("pvc reclaim %s still in progress",
			nc.GraphdComponent().GetName())}
	}

	klog.V(4).Infof("graphd cluster [%s/%s] used pvcs were reclaimed", namespace, componentName)
	nc.Status.Graphd.Phase = v2alpha1.RunningPhase
	return nil
}
