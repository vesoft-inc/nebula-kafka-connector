package zone

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

var RemoveCmd = &cobra.Command{
	Use:   "remove-zone <svcgrp-name>",
	Short: "Remove zone",
	Long:  "Remove zone",
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
		req := meta.NewDropZoneReq(svcgrpName, zoneName)
		if err := common.MetaClient.DropZone(req); err != nil {
			return common.NgctlError("drop zone failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "drop zone successfully")
		return nil
	},
}

func init() {
	RemoveCmd.Flags().StringVarP(&zoneName, "name", "n", "", "zone name")
	RemoveCmd.MarkFlagRequired("name")
}
