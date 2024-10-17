package srvgrp_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var alterSrvgrpCmd = &cobra.Command{
	Use:   "alter",
	Short: "Alter the owner of a srvgrp.",
	Long:  `A service group has an owner registered in the metad. Users can alter it.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if srvgrpFlags.srvgrpName == "" {
			return common.NgctlError("ServiceGroup name is empty", "")
		}
		if srvgrpFlags.owner == "" {
			return common.NgctlError("Owner is invalid", "")
		}
		// donot need to specify zones
		req := meta.NewAlterServiceGroupReq(srvgrpFlags.srvgrpName, srvgrpFlags.owner)
		if err := common.MetaClient.AlterServiceGroup(req); err != nil {
			return common.NgctlError("Alter srvgrp failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Alter srvgrp successfully.")
		return nil
	},
}
