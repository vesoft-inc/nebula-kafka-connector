package zone

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

var AddCmd = &cobra.Command{
	Use:   "add-zone",
	Short: "Add zone",
	Long:  "Add zone",
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

		req := meta.NewAddZoneReq(svcgrpName, zoneName, replicas, priority)
		if err := common.MetaClient.AddZone(req); err != nil {
			return common.NgctlError("add zone failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "add zone successfully.")
		return nil
	},
}

func init() {
	AddCmd.Flags().StringVarP(&zoneName, "name", "n", "", "zone name")
	AddCmd.Flags().IntVarP(&replicas, "replicas", "r", 0, "zone replicas")
	AddCmd.Flags().IntVarP(&priority, "priority", "p", 0, "zone priority")
	AddCmd.MarkFlagRequired("name")
	AddCmd.SetUsageTemplate(common.GetUsageTemplate("ngctl svcgrp add-zone <svcgrp-name> [flags]"))
}
