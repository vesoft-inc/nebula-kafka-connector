package service_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var addServiceCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a service into a service group",
	Long:  "Add a service info a service group",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateAddorDropFlags(); err != nil {
			return err
		}
		serviceType, err := getServiceType(serviceFlags.serviceType)
		if err != nil {
			return err
		}
		req := meta.NewAddServiceReq(serviceFlags.host, uint32(serviceFlags.port), serviceType, serviceFlags.svcgrpName)
		if err := common.MetaClient.AddService(req); err != nil {
			return common.NgctlError("Add service failed", err.Error())
		}

		fmt.Fprintln(common.MetaOutput, "Add services successfully.")
		return nil
	},
}

func init() {
	addServiceCmd.Flags().StringVarP(&serviceFlags.serviceType, "type", "t", "", "service type")
	addServiceCmd.Flags().StringVarP(&serviceFlags.host, "host", "H", "", "service host")
	addServiceCmd.Flags().Int32VarP(&serviceFlags.port, "port", "P", -1, "service port")
	addServiceCmd.Flags().StringVarP(&serviceFlags.svcgrpName, "svcgrp", "s", "", "svcgrp name")
}
