package service

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

var DropServiceCmd = &cobra.Command{
	Use:   "drop-service <svcgrp-name>",
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
		req := meta.NewDropServiceReq(serviceFlags.host, uint32(serviceFlags.port),
			serviceType, svcgrpName)
		if err := common.MetaClient.DropService(req); err != nil {
			return common.NgctlError("Drop service failed", err.Error())
		}

		fmt.Fprintln(common.MetaOutput, "Drop services successfully.")
		return nil
	},
}

func init() {
	DropServiceCmd.Flags().StringVarP(&serviceFlags.serviceType, "type", "t", "", "service type")
	DropServiceCmd.Flags().StringVarP(&serviceFlags.host, "host", "H", "", "service host")
	DropServiceCmd.Flags().Int32VarP(&serviceFlags.port, "port", "P", -1, "service port")
}
