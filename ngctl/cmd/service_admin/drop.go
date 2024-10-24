package service_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var dropServiceCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop a service from a svcgrp",
	Long:  "Drop a service from a svcgrp",
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
		req := meta.NewDropServiceReq(serviceFlags.host, uint32(serviceFlags.port),
			serviceType, serviceFlags.svcgrpName)
		if err := common.MetaClient.DropService(req); err != nil {
			return common.NgctlError("Drop service failed", err.Error())
		}

		fmt.Fprintln(common.MetaOutput, "Drop services successfully.")
		return nil
	},
}

func init() {
	dropServiceCmd.Flags().StringVarP(&serviceFlags.serviceType, "type", "t", "", "service type")
	dropServiceCmd.Flags().StringVarP(&serviceFlags.host, "host", "H", "", "service host")
	dropServiceCmd.Flags().Int32VarP(&serviceFlags.port, "port", "P", -1, "service port")
	dropServiceCmd.Flags().StringVarP(&serviceFlags.svcgrpName, "svcgrp", "s", "", "svcgrp name")
}
