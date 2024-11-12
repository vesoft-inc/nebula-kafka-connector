package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

var configPath string
var opts CreateOptions

func main() {
	cmd := &cobra.Command{
		Use:                   "quota-cli",
		DisableFlagsInUseLine: true,
		Short:                 "A tool for k8s user and quota management",
		Long:                  `A tool for k8s user and quota management. Use 'quota-cli -h' to see usage.`,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVarP(&configPath, "config-path", "c", "", "Path to the kubeconfig file to use for CLI requests.")
	createUserQuotaCmd.PersistentFlags().StringVar(&opts.QuotaUser, "quota-user", "", "Username for basic authentication to the API server.")
	createUserQuotaCmd.PersistentFlags().StringVar(&opts.QuotaNamespace, "quota-namespace", "", "The namespace for the nebula cluster request.")
	createUserQuotaCmd.PersistentFlags().StringVar(&opts.ClusterName, "cluster-name", "", "The name of the kubeconfig cluster to use, run cmd 'kubectl config current-context' to confirm.")
	createUserQuotaCmd.PersistentFlags().StringVar(&opts.UserConfigPath, "user-config-path", "kube", "Path to the saved kubeconfig file to use for the quota user.")
	createUserQuotaCmd.PersistentFlags().StringVar(&opts.ResourceRequests, "resource-requests", "", "The compute resource requests total for all Pods in the namespace for the quota user, For example, 'cpu=4,memory=8Gi'.")
	createUserQuotaCmd.PersistentFlags().StringVar(&opts.ResourceLimits, "resource-limits", "", "The compute resource limits total for all Pods in the namespace for the quota user, For example, 'cpu=8,memory=16Gi'.")
	createUserQuotaCmd.PersistentFlags().StringVar(&opts.CertPath, "cert-path", "certs", "The directory that the certificate signing request (CSR) and private key will be written.")

	cmd.AddCommand(createUserQuotaCmd)
	cmd.AddCommand(showUserQuotaCmd)

	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}

var createUserQuotaCmd = &cobra.Command{
	Use:   "create",
	Short: "create user and resource quota for nebula graph",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := opts.Validate(); err != nil {
			return err
		}
		return createUserQuota(context.TODO(), &opts)
	},
}

var showUserQuotaCmd = &cobra.Command{
	Use:   "list",
	Short: "list all user and resource quotas used for nebula graph",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listUserQuotas(context.TODO())
	},
}

func createUserQuota(ctx context.Context, opts *CreateOptions) error {
	keyData, csrData, err := GenerateCert(opts.QuotaUser)
	if err != nil {
		return err
	}
	if err = WriteCertsToDir(opts.CertPath, opts.QuotaUser, keyData, csrData); err != nil {
		return err
	}
	klog.Infoln("Generate certs done!")

	client, err := InitClient(configPath)
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
	apiConfig, err := LoadConfig(configPath)
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
	if err = qc.CreateResourceQuota(ctx, opts.QuotaNamespace, opts.QuotaUser, resourceRequirements); err != nil {
		klog.Error(err)
		return err
	}
	klog.Infoln("Create compute resource quota done!")

	return nil
}

func listUserQuotas(ctx context.Context) error {
	client, err := InitClient(configPath)
	if err != nil {
		return err
	}
	qc := &QuotaClient{client: client}
	resourceQuotas, err := qc.ListResourceQuotas(ctx)
	if err != nil {
		klog.Error(err)
		return err
	}

	first := true
	for i := range resourceQuotas {
		rq := &resourceQuotas[i]
		s, err := describeQuota(rq)
		if err != nil {
			return err
		}
		if first {
			first = false
			fmt.Fprint(os.Stdout, s)
		} else {
			fmt.Fprintf(os.Stdout, "\n\n%s", s)
		}
	}

	if len(resourceQuotas) == 0 {
		fmt.Fprintln(os.Stderr, "No resources found")
	}

	return nil
}
