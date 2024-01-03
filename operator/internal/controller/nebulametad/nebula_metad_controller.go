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

package nebulametad

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/kubernetes"
	"k8s.io/kubernetes/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/controller/component"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/kube"
	errorsutil "github.com/vesoft-inc/nebula-ng-tools/operator/internal/util/errors"
)

const (
	defaultTimeout   = 5 * time.Second
	reconcileTimeOut = 10 * time.Second

	finalizerKey = "apps.nebula-graph.io/metad-cleanup"
)

// NebulaMetadReconciler reconciles a NebulaMetad object
type NebulaMetadReconciler struct {
	control ControlInterface
	client.Client
	Scheme *runtime.Scheme
}

func NewMetadReconciler(mgr ctrl.Manager) (*NebulaMetadReconciler, error) {
	clientSet, err := kube.NewClientSet(mgr.GetConfig())
	if err != nil {
		return nil, err
	}

	metadUpdater := component.NewMetadUpdater(clientSet.Pod())

	kubeClient, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		return nil, fmt.Errorf("create kubernetes client failed: %v", err)
	}
	eventBroadcaster := record.NewBroadcasterWithCorrelatorOptions(record.CorrelatorOptions{QPS: 1})
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedv1.EventSinkImpl{Interface: typedv1.New(kubeClient.CoreV1().RESTClient()).Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "nebula-metad-controller"})

	return &NebulaMetadReconciler{
		control: NewMetadControl(mgr.GetClient(), component.NewMetadManager(clientSet, recorder, metadUpdater)),
		Client:  mgr.GetClient(),
		Scheme:  mgr.GetScheme(),
	}, nil
}

//+kubebuilder:rbac:groups=apps.nebula-graph.io,resources=nebulametads,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.nebula-graph.io,resources=nebulametads/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.nebula-graph.io,resources=nebulametads/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NebulaMetad object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *NebulaMetadReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, retErr error) {
	key := req.NamespacedName.String()
	subCtx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	startTime := time.Now()
	defer func() {
		if retErr == nil {
			if res.Requeue || res.RequeueAfter > 0 {
				klog.Infof("Finished reconciling NebulaCluster [%s] (%v), result: %v", key, time.Since(startTime), res)
			} else {
				klog.Infof("Finished reconciling NebulaCluster [%s], spendTime: (%v)", key, time.Since(startTime))
			}
		} else {
			klog.Errorf("Failed to reconcile NebulaCluster [%s], spendTime: (%v)", key, time.Since(startTime))
		}
	}()

	var nebulaMetad v2alpha1.NebulaMetad
	if err := r.Get(subCtx, req.NamespacedName, &v2alpha1.NebulaMetad{}); err != nil {
		if apierrors.IsNotFound(err) {
			klog.Info("Skip because the loadBalance does not exist anymore")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	klog.Info("Start to reconcile nebulaMetad")
	if err := r.syncNebulaMetad(ctx, nebulaMetad.DeepCopy()); err != nil {
		klog.Errorf("NebulaMetad [%s] reconcile failed: %v", key, err)

		if errorsutil.IsReconcileError(err) {
			return ctrl.Result{RequeueAfter: reconcileTimeOut}, nil
		}
		return ctrl.Result{RequeueAfter: defaultTimeout}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NebulaMetadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v2alpha1.NebulaMetad{}).
		Complete(r)
}
