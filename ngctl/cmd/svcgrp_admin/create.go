package svcgrp_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var createsvcgrpCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a service group in the metad",
	Long:  "Create a service group and register in the metad. The service group will be empty after creation. Users need to add hosts and services into a newly created svcgrp",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if svcgrpFlags.svcgrpName == "" {
			return common.NgctlError("ServiceGroup name is empty", "")
		}
		if svcgrpFlags.replicas == 0 {
			return common.NgctlError("Replica number is invalid", "")
		}
		// donot need to specify zones
		req := meta.NewCreateServiceGroupReq(svcgrpFlags.svcgrpName, svcgrpFlags.replicas,
			svcgrpFlags.owner, nil)
		if err := common.MetaClient.CreateServiceGroup(req); err != nil {
			return common.NgctlError("Create svcgrp failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Create svcgrp successfully.")
		return nil
	},
}

func init() {
	createsvcgrpCmd.Flags().IntVarP(&svcgrpFlags.replicas, "replica-factor", "r", 3, "replica number, default: 3")
	createsvcgrpCmd.Flags().StringVarP(&svcgrpFlags.owner, "owner", "o", "", "svcgrp owner")
}
