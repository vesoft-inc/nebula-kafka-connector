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

package nebulacluster

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/controller/component"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/controller/component/reclaimer"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/discovery"
	utilerrors "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
)

const (
	defaultTimeout   = 5 * time.Second
	reconcileTimeOut = 10 * time.Second
)

// ClusterReconciler reconciles a NebulaCluster object
type ClusterReconciler struct {
	control ControlInterface
	client  client.Client
}

func NewClusterReconciler(mgr ctrl.Manager) (*ClusterReconciler, error) {
	clientSet, err := kube.NewClientSet(mgr.GetConfig())
	if err != nil {
		return nil, err
	}

	storageScaler := component.NewStorageScaler(clientSet)
	graphScaler := component.NewGraphScaler(clientSet)
	graphdUpdater := component.NewGraphdUpdater(clientSet.Pod())
	storagedUpdater := component.NewStorageUpdater(clientSet.Pod())

	dm, err := discovery.New(mgr.GetConfig())
	if err != nil {
		return nil, fmt.Errorf("create discovery client failed: %v", err)
	}
	info, err := dm.GetServerVersion()
	if err != nil {
		return nil, fmt.Errorf("create apiserver info failed: %v", err)
	}

	valid, err := kube.ValidVersion(info)
	if err != nil {
		return nil, fmt.Errorf("get server version failed: %v", err)
	}
	if !valid {
		return nil, fmt.Errorf("server version not supported")
	}

	kubeClient, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		return nil, fmt.Errorf("create kubernetes client failed: %v", err)
	}
	eventBroadcaster := record.NewBroadcasterWithCorrelatorOptions(record.CorrelatorOptions{QPS: 1})
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedv1.EventSinkImpl{Interface: typedv1.New(kubeClient.CoreV1().RESTClient()).Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "nebula-cluster-controller"})

	return &ClusterReconciler{
		control: NewDefaultNebulaClusterControl(
			mgr.GetClient(),
			clientSet.NebulaCluster(),
			component.NewGraphdCluster(
				clientSet,
				graphScaler,
				graphdUpdater,
				recorder),
			component.NewStoragedCluster(
				clientSet,
				storageScaler,
				storagedUpdater,
				recorder),
			reclaimer.NewMetaReconciler(clientSet),
			reclaimer.NewPVCReclaimer(clientSet),
			NewClusterConditionUpdater(),
		),
		client: mgr.GetClient(),
	}, nil
}

// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch
// +kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=persistentvolumes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch;list
// +kubebuilder:rbac:groups=apps.nebula-graph.io,resources=nebulaclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.nebula-graph.io,resources=nebulaclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.nebula-graph.io,resources=nebulaclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *ClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res reconcile.Result, retErr error) {
	key := req.NamespacedName.String()
	subCtx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	startTime := time.Now()
	defer func() {
		if retErr == nil {
			if res.Requeue || res.RequeueAfter > 0 {
				klog.V(4).Infof("Finished reconciling NebulaCluster [%s] (%v), result: %v", key, time.Since(startTime), res)
			} else {
				klog.V(4).Infof("Finished reconciling NebulaCluster [%s], spendTime: (%v)", key, time.Since(startTime))
			}
		} else {
			klog.Errorf("Failed to reconcile NebulaCluster [%s], spendTime: (%v)", key, time.Since(startTime))
		}
	}()

	var nebulaCluster v2alpha1.NebulaCluster
	if err := r.client.Get(subCtx, req.NamespacedName, &nebulaCluster); err != nil {
		if apierrors.IsNotFound(err) {
			klog.Infof("Skipping because NebulaCluster [%s] has been deleted", key)
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	klog.Info("Start to reconcile NebulaCluster")

	if err := r.syncNebulaCluster(nebulaCluster.DeepCopy()); err != nil {
		if utilerrors.IsReconcileError(err) {
			klog.Infof("NebulaCluster [%s] reconcile details: %v", key, err)
			return ctrl.Result{RequeueAfter: reconcileTimeOut}, nil
		}
		klog.Errorf("NebulaCluster [%s] reconcile failed: %v", key, err)
		return ctrl.Result{RequeueAfter: defaultTimeout}, nil
	}
	return ctrl.Result{}, nil
}

func (r *ClusterReconciler) syncNebulaCluster(nc *v2alpha1.NebulaCluster) error {
	if nc.DeletionTimestamp != nil {
		return r.control.DeleteCluster(nc)
	}
	return r.control.UpdateNebulaCluster(nc)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v2alpha1.NebulaCluster{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}
