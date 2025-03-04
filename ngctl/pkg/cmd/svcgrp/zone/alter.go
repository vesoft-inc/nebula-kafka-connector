package zone

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

var AlterCmd = &cobra.Command{
	Use:   "alter-zone <svcgrp-name>",
	Short: "Alter zone",
	Long:  "Alter zone",
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
		getReq := meta.NewListZonesReq(svcgrpName)
		zones, err := common.MetaClient.ListZones(getReq)
		if err != nil {
			return common.NgctlError("list zones failed", err.Error())
		}
		var originalZone *meta.ZoneInfo
		for _, z := range zones.Zones {
			if z.Name == zoneName {
				originalZone = z
				break
			}
		}
		if originalZone == nil {
			return common.NgctlError("zone not found", "")
		}
		newName := originalZone.Name
		newPriority := originalZone.Priority
		if newZoneName != "" {
			newName = newZoneName
		}
		if newZonePriority != -1 {
			newPriority = newZonePriority
		}
		req := meta.NewAlterZoneReq(svcgrpName, zoneName, newName, newPriority)

		if err := common.MetaClient.AlterZone(req); err != nil {
			return common.NgctlError("rename zone failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "rename zone successfully.")
		return nil
	},
}

func init() {
	AlterCmd.Flags().StringVarP(&zoneName, "name", "n", "", "zone name")
	AlterCmd.Flags().StringVarP(&newZoneName, "new-name", "E", "", "new zone name")
	AlterCmd.Flags().IntVarP(&newZonePriority, "new-priority", "P", -1, "new zone priority")
	AlterCmd.MarkFlagRequired("name")
}
