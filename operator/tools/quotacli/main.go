package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

type Options struct {
	QuotaUser string

	QuotaNamespace string

	ClusterName string

	ConfigPath string

	UserConfigPath string

	ResourceRequests string

	ResourceLimits string

	CertPath string
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.QuotaUser, "quota-user", "", "Username for basic authentication to the API server.")
	fs.StringVar(&o.QuotaNamespace, "quota-namespace", "", "The namespace for the nebula cluster request.")
	fs.StringVar(&o.ClusterName, "cluster-name", "", "The name of the kubeconfig cluster to use, run cmd 'kubectl config current-context' to confirm.")
	fs.StringVar(&o.ConfigPath, "config-path", "", "Path to the kubeconfig file to use for CLI requests.")
	fs.StringVar(&o.UserConfigPath, "user-config-path", "kube", "Path to the saved kubeconfig file to use for the quota user.")
	fs.StringVar(&o.ResourceRequests, "resource-requests", "", "The compute resource requests total for all Pods in the namespace for the quota user, For example, 'cpu=4,memory=8Gi'.")
	fs.StringVar(&o.ResourceLimits, "resource-limits", "", "The compute resource limits total for all Pods in the namespace for the quota user, For example, 'cpu=8,memory=16Gi'.")
	fs.StringVar(&o.CertPath, "cert-path", "certs", "The directory that the certificate signing request (CSR) and private key will be written.")
}

func main() {
	opts := &Options{}
	cmd := &cobra.Command{
		Use:                   "quota-cli",
		DisableFlagsInUseLine: true,
		Short:                 "A tool for k8s user and quota management",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(context.Background(), opts)
		},
	}

	fs := pflag.NewFlagSet("generic", pflag.ExitOnError)
	opts.AddFlags(fs)
	cmd.Flags().AddFlagSet(fs)

	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}

func Run(ctx context.Context, opts *Options) error {
	if err := validate(opts); err != nil {
		return err
	}

	keyData, csrData, err := GenerateCert(opts.QuotaUser)
	if err != nil {
		return err
	}
	if err = WriteCertsToDir(opts.CertPath, opts.QuotaUser, keyData, csrData); err != nil {
		return err
	}
	klog.Infoln("Generate certs done!")

	client, err := InitClient(opts.ConfigPath)
	if err != nil {
		return err
	}

	qc := &QuotaClient{client: client}

	if err = qc.CreateNamespace(ctx, opts.QuotaNamespace); err != nil {
		klog.Error(err)
		return err
	}
	klog.Infoln("Create namespace done!")

	if err = qc.CreateRole(ctx, opts.QuotaUser, opts.QuotaNamespace); err != nil {
		klog.Error(err)
		return err
	}
	if err = qc.CreateRoleBinding(ctx, opts.QuotaUser, opts.QuotaNamespace); err != nil {
		klog.Error(err)
		return err
	}
	klog.Infoln("Create role and rolebinding done!")

	csrFile := getCsrFileName(opts.QuotaUser)
	csrFilePath := path.Join(opts.CertPath, csrFile)
	data, err := os.ReadFile(csrFilePath)
	if err != nil {
		klog.Errorf("failed to read file: %v", err)
	}

	csrName := getCsrName(opts.QuotaUser)
	if err = qc.DeleteCSR(ctx, csrName); err != nil {
		return err
	}
	if err = qc.CreateCSR(ctx, opts.QuotaUser, data); err != nil {
		klog.Error(err)
		return err
	}
	if err = qc.ApproveCSR(ctx, csrName); err != nil {
		klog.Error(err)
		return err
	}
	klog.Infoln("Approve CSR done!")

	csr, err := qc.GetCSR(ctx, csrName)
	if err != nil {
		return err
	}
	for i := 0; ; i++ {
		if csr.Status.Certificate == nil {
			klog.Infof("csr %s status certificate is nil", csr.Name)
			csr, err = qc.GetCSR(ctx, getCsrName(opts.QuotaUser))
			if err != nil {
				return err
			}
			time.Sleep(time.Millisecond * 200)
			if i < 10 {
				continue
			}
		}
		break
	}

	keyName := getKeyFileName(opts.QuotaUser)
	keyFile := path.Join(opts.CertPath, keyName)
	keyString := encodeToString(keyFile)
	apiConfig, err := LoadConfig(opts.ConfigPath)
	if err != nil {
		klog.Error(err)
		return err
	}
	userConfig := convert(apiConfig, csr, keyString, opts.ClusterName, opts.QuotaUser, opts.QuotaNamespace)
	out, err := yaml.Marshal(userConfig)
	if err != nil {
		return err
	}
	if err = WriteConfigToFile(out, opts.UserConfigPath, getUserConfig(opts.QuotaUser)); err != nil {
		klog.Error(err)
		return err
	}
	klog.Infoln("Save user kubeconfig done!")

	resourceRequirements, err := handleResourceRequirementsV1(map[string]string{"limits": opts.ResourceLimits, "requests": opts.ResourceRequests})
	if err != nil {
		return fmt.Errorf("failed to handle resource: %v", err)
	}
	if err = qc.CreateResourceQuota(ctx, opts.QuotaNamespace, resourceRequirements); err != nil {
		klog.Error(err)
		return err
	}
	klog.Infoln("Create compute resource quota done!")

	return nil
}
