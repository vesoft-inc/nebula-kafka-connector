package zone

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

var RenameCmd = &cobra.Command{
	Use:   "rename-zone",
	Short: "Rename zone",
	Long:  "Rename zone",
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
		req := meta.NewRenameZoneReq(svcgrpName, zoneName, newZoneName)

		if err := common.MetaClient.RenameZone(req); err != nil {
			return common.NgctlError("rename zone failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "rename zone successfully.")
		return nil
	},
}

func init() {
	RenameCmd.Flags().StringVarP(&zoneName, "name", "n", "", "zone name")
	RenameCmd.Flags().StringVarP(&newZoneName, "new-name", "E", "", "new zone name")
	RenameCmd.MarkFlagRequired("name")
	RenameCmd.SetUsageTemplate(common.GetUsageTemplate("ngctl svcgrp rename-zone <svcgrp-name> [flags]"))
}
