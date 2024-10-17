package srvgrp_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var dropSrvgrpCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop a service group from the metad.",
	// what is the service group is not empty?
	Long: `Drop a service group from the metad.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		srvgrp := srvgrpFlags.srvgrpName
		force := srvgrpFlags.force
		if srvgrp == "" {
			return common.NgctlError("srvgrp name is empty", "")
		}

		req := meta.NewDropServiceGroupReq(srvgrp, force)
		if err := common.MetaClient.DropServiceGroup(req); err != nil {
			return common.NgctlError("Drop srvgrp failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Drop srvgrp successfully.")
		return nil
	},
}
