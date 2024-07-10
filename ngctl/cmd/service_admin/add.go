package service_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var addServiceCmd = &cobra.Command{
	Use:   "add",
	Short: `Add service into assigned cluster.`,
	Long:  `ngctl service add --type [graphd|storaged] --host [host] --port [port] --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := ServiceFlags
		var serviceType meta.ServiceType
		if flags.serviceType == "graphd" {
			serviceType = meta.ServiceTypeGraphd
		} else if flags.serviceType == "storaged" {
			serviceType = meta.ServiceTypeStoraged
		} else {
			return common.NgctlError("service type is not correct, valid value is graphd or storaged", "")
		}
		if flags.host == "" || common.IsValidIPAddress(flags.host) == false {
			return common.NgctlError("no valid host provided", "")
		}
		if flags.clusterName == "" {
			return common.NgctlError("cluster name is empty", "")
		}
		req := meta.NewAddServiceReq(flags.host, flags.port, serviceType, flags.clusterName)
		if err := common.MetaClient.AddService(req); err != nil {
			return common.NgctlError("Add service failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Add service successfully.")
		return nil
	},
}

func init() {
	addServiceCmd.Flags().StringVarP(&ServiceFlags.serviceType, "type", "t", "", "service type")
	addServiceCmd.Flags().StringVarP(&ServiceFlags.host, "host", "H", "", "service host")
	addServiceCmd.Flags().Uint32VarP(&ServiceFlags.port, "port", "P", 0, "service port")
	addServiceCmd.Flags().StringVarP(&ServiceFlags.clusterName, "cluster", "c", "", "cluster name")
	addServiceCmd.MarkFlagsRequiredTogether("type", "host", "port", "cluster")
}
