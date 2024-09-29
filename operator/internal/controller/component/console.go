/*
Copyright 2024 Vesoft Inc.

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
	"strconv"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v1alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/pkg/label"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	utilerrors "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
)

const (
	defaultConsoleImage = "vesoft/ngql"
)

type nebulaConsole struct {
	clientSet kube.ClientSet
}

func NewNebulaConsole(clientSet kube.ClientSet) ReconcileManager {
	return &nebulaConsole{clientSet: clientSet}
}

func (c *nebulaConsole) Reconcile(_ meta.Client, nc *v1alpha1.NebulaCluster) error {
	if nc.Spec.Console == nil {
		return nil
	}
	if !nc.GraphdComponent().IsReady() {
		return utilerrors.ReconcileErrorf("waiting for graphd cluster [%s/%s] ready", nc.Namespace, nc.GraphdComponent().GetName())
	}
	return c.syncConsolePod(nc)
}

func (c *nebulaConsole) Delete(nc *v1alpha1.NebulaCluster) error {
	if nc.Spec.Console == nil {
		return nil
	}
	podName := getConsolePodName(nc.Name)
	_, err := c.clientSet.Pod().GetPod(nc.Namespace, podName)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return c.clientSet.Pod().DeletePod(nc.Namespace, podName, true)
}

func (c *nebulaConsole) syncConsolePod(nc *v1alpha1.NebulaCluster) error {
	newPod, err := c.generatePod(nc)
	if err != nil {
		return err
	}
	oldPod, err := c.clientSet.Pod().GetPod(newPod.Namespace, newPod.Name)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	notExist := apierrors.IsNotFound(err)
	if notExist {
		if err := setPodLastAppliedConfigAnnotation(newPod); err != nil {
			return err
		}
		return c.clientSet.Pod().CreatePod(newPod)
	}
	return updateSinglePod(c.clientSet, newPod, oldPod)
}

func (c *nebulaConsole) generatePod(nc *v1alpha1.NebulaCluster) (*corev1.Pod, error) {
	mounts := make([]corev1.VolumeMount, 0)
	cmd := []string{
		"ngql",
		"-H",
		nc.GetGraphdServiceName(),
		"-P",
		strconv.Itoa(int(nc.GraphdComponent().GetPort(v1alpha1.GraphdPortNameGRPC))),
	}

	username, password, err := kube.GetCredential(c.clientSet, nc.Namespace, nc.Spec.CredentialSecret)
	if err != nil {
		return nil, err
	}
	cmd = append(cmd, "-u", username, "-p", password)

	container := corev1.Container{
		Name:            "console",
		Image:           nc.GetConsoleImage(),
		ImagePullPolicy: corev1.PullAlways,
		Command:         cmd,
		VolumeMounts:    mounts,
		Stdin:           true,
		StdinOnce:       true,
		TTY:             true,
	}

	volumes := make([]corev1.Volume, 0)
	podName := getConsolePodName(nc.Name)
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:            podName,
			Namespace:       nc.GetNamespace(),
			Labels:          c.getConsoleLabels(nc),
			OwnerReferences: nc.GenerateOwnerReferences(),
		},
		Spec: corev1.PodSpec{
			SchedulerName:      nc.Spec.SchedulerName,
			NodeSelector:       c.getNodeSelector(nc),
			Containers:         []corev1.Container{container},
			ImagePullSecrets:   nc.Spec.ImagePullSecrets,
			ServiceAccountName: v1alpha1.NebulaServiceAccountName,
			Volumes:            volumes,
		},
	}, nil
}

func (c *nebulaConsole) getConsoleLabels(nc *v1alpha1.NebulaCluster) map[string]string {
	selector := label.New().Cluster(nc.GetName()).Console()
	labels := selector.Copy().Labels()
	return labels
}

func (c *nebulaConsole) getNodeSelector(nc *v1alpha1.NebulaCluster) map[string]string {
	selector := map[string]string{}
	for k, v := range nc.Spec.NodeSelector {
		selector[k] = v
	}
	consoleSelector := nc.Spec.Console.NodeSelector
	if consoleSelector != nil {
		for k, v := range consoleSelector {
			selector[k] = v
		}
	}
	return selector
}

func getConsolePodName(clusterName string) string {
	return clusterName + "-console"
}
