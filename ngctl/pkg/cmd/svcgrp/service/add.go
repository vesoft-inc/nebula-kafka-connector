package service

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

var AddServiceCmd = &cobra.Command{
	Use:   "add-service",
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
		svcgrpName, err := common.GetResourceName(args)
		if err != nil {
			return err
		}
		if err := validateAddOrDropFlags(); err != nil {
			return err
		}
		serviceType, err := getServiceType(serviceFlags.serviceType)
		if err != nil {
			return err
		}
		req := meta.NewAddServiceReq(serviceFlags.host, uint32(serviceFlags.port), serviceType, svcgrpName)
		if err := common.MetaClient.AddService(req); err != nil {
			return common.NgctlError("Add service failed", err.Error())
		}

		fmt.Fprintln(common.MetaOutput, "Add services successfully.")
		return nil
	},
}

func init() {
	AddServiceCmd.Flags().StringVarP(&serviceFlags.serviceType, "type", "t", "", "service type")
	AddServiceCmd.Flags().StringVarP(&serviceFlags.host, "host", "H", "", "service host")
	AddServiceCmd.Flags().Int32VarP(&serviceFlags.port, "port", "P", -1, "service port")
	AddServiceCmd.SetUsageTemplate(common.GetUsageTemplate("ngctl svcgrp add-service <svcgrp-name> [flags]"))
}
