package srvgrp_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var initSrvgrpCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a service group in the metad.",
	Long:  `Initialize a service group in the metad.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		srvgrp := srvgrpFlags.srvgrpName
		if srvgrp == "" {
			return common.NgctlError("srvgrp name is empty", "")
		}
		req := meta.NewInitServiceGroupReq(srvgrp)

		if err := common.MetaClient.InitServiceGroup(req); err != nil {
			return common.NgctlError("Init srvgrp failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Init srvgrp successfully.")
		return nil
	},
}
