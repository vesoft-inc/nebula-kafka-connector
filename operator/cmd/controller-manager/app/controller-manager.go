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

package app

import (
	"context"
	"flag"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	ms "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/vesoft-inc/nebula-ng-tools/operator/apis/apps/v2alpha1"
	"github.com/vesoft-inc/nebula-ng-tools/operator/cmd/controller-manager/app/options"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/controller/nebulacluster"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/controller/nebulametad"
	klogflag "github.com/vesoft-inc/nebula-ng-tools/operator/internal/flag/klog"
	profileflag "github.com/vesoft-inc/nebula-ng-tools/operator/internal/flag/profile"
	"github.com/vesoft-inc/nebula-ng-tools/operator/internal/version"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v2alpha1.AddToScheme(clientgoscheme.Scheme))
	utilruntime.Must(v2alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

// NewControllerManagerCommand creates a *cobra.Command object with default parameters
func NewControllerManagerCommand(ctx context.Context) *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use: "nebula-controller-manager",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(ctx, opts)
		},
	}

	nfs := cliflag.NamedFlagSets{}
	fs := nfs.FlagSet("generic")
	fs.AddGoFlagSet(flag.CommandLine)
	opts.AddFlags(fs)

	logsFlagSet := nfs.FlagSet("logs")
	klogflag.Add(logsFlagSet)

	cmd.Flags().AddFlagSet(fs)
	cmd.Flags().AddFlagSet(logsFlagSet)

	return cmd
}

// Run runs the controller-manager with options. This should never exit.
func Run(ctx context.Context, opts *options.Options) error {
	klog.Infof("nebula-controller-manager version: %s", version.Version())

	logf.SetLogger(klog.Background())

	profileflag.ListenAndServe(opts.ProfileOpts)

	if len(opts.Namespaces) == 0 {
		klog.Info("nebula-controller-manager watches all namespaces")
	} else {
		klog.Infof("nebula-controller-manager watches namespaces %v", opts.Namespaces)
	}

	cfg, err := ctrlruntime.GetConfig()
	if err != nil {
		panic(err)
	}

	defaultNamespaces := make(map[string]cache.Config)
	for _, ns := range opts.Namespaces {
		defaultNamespaces[ns] = cache.Config{}
	}
	ctrlOptions := ctrlruntime.Options{
		Scheme:                     scheme,
		Logger:                     klog.Background(),
		LeaderElection:             opts.LeaderElection.LeaderElect,
		LeaderElectionID:           opts.LeaderElection.ResourceName,
		LeaderElectionNamespace:    opts.LeaderElection.ResourceNamespace,
		LeaseDuration:              &opts.LeaderElection.LeaseDuration.Duration,
		RenewDeadline:              &opts.LeaderElection.RenewDeadline.Duration,
		RetryPeriod:                &opts.LeaderElection.RetryPeriod.Duration,
		LeaderElectionResourceLock: opts.LeaderElection.ResourceLock,
		HealthProbeBindAddress:     opts.HealthProbeBindAddress,
		Metrics: ms.Options{
			BindAddress: opts.MetricsBindAddress,
		},
		Cache: cache.Options{
			SyncPeriod:        &opts.SyncPeriod.Duration,
			DefaultNamespaces: defaultNamespaces,
		},
		Controller: config.Controller{
			GroupKindConcurrency: map[string]int{
				v2alpha1.SchemeGroupVersion.WithKind("NebulaCluster").GroupKind().String(): opts.ConcurrentNebulaClusterSyncs,
			},
			RecoverPanic: pointer.Bool(true),
		},
	}

	if opts.NebulaSelector != "" {
		parsedSelector, err := labels.Parse(opts.NebulaSelector)
		if err != nil {
			klog.Errorf("couldn't convert selector into a corresponding internal selector object: %v", err)
			return err
		}
		ctrlOptions.Cache.DefaultLabelSelector = parsedSelector
	}

	mgr, err := ctrlruntime.NewManager(cfg, ctrlOptions)
	if err != nil {
		klog.Errorf("Failed to build controller manager: %v", err)
		return err
	}

	clusterReconciler, err := nebulacluster.NewClusterReconciler(mgr)
	if err != nil {
		return err
	}
	if err := clusterReconciler.SetupWithManager(mgr); err != nil {
		klog.Errorf("failed to set up NebulaCluster controller: %v", err)
		return err
	}

	metadReconciler, err := nebulametad.NewMetadReconciler(mgr)
	if err != nil {
		return err
	}
	if err := metadReconciler.SetupWithManager(mgr); err != nil {
		klog.Errorf("failed to set up NebulaMetad controller: %v", err)
		return err
	}

	if err := mgr.AddHealthzCheck("ping", healthz.Ping); err != nil {
		klog.Errorf("failed to add health check endpoint: %v", err)
		return err
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		klog.Errorf("failed to add ready check endpoint: %v", err)
		return err
	}

	if err := mgr.Start(ctx); err != nil {
		klog.Errorf("nebula-controller-manager exits unexpectedly: %v", err)
		return err
	}

	return nil
}
