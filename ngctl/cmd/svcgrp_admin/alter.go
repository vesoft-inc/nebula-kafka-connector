package svcgrp_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var altersvcgrpCmd = &cobra.Command{
	Use:   "alter",
	Short: "Alter the owner of a svcgrp.",
	Long:  `A service group has an owner registered in the metad. Users can alter it.`,
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
		if svcgrpFlags.owner == "" {
			return common.NgctlError("Owner is invalid", "")
		}
		// donot need to specify zones
		req := meta.NewAlterServiceGroupReq(svcgrpFlags.svcgrpName, svcgrpFlags.owner)
		if err := common.MetaClient.AlterServiceGroup(req); err != nil {
			return common.NgctlError("Alter svcgrp failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Alter svcgrp successfully.")
		return nil
	},
}
