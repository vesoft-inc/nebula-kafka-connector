package svcgrp_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var initsvcgrpCmd = &cobra.Command{
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
		svcgrp := svcgrpFlags.svcgrpName
		if svcgrp == "" {
			return common.NgctlError("svcgrp name is empty", "")
		}
		req := meta.NewInitServiceGroupReq(svcgrp)

		if err := common.MetaClient.InitServiceGroup(req); err != nil {
			return common.NgctlError("Init svcgrp failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Init svcgrp successfully.")
		return nil
	},
}
