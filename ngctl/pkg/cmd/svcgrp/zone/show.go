package zone

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

var ShowCmd = &cobra.Command{
	Use:   "show-zone <svcgrp-name>",
	Short: "Show zone",
	Long:  "Show zone",
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
		req := meta.NewListZonesReq(svcgrpName)
		zones, err := common.MetaClient.ListZones(req)
		if err != nil {
			return common.NgctlError("list zones failed", err.Error())
		}
		header := []string{"Name", "Replica Factor", "Priority", "Hosts"}
		data := make([][]string, 0)
		for _, z := range zones.Zones {
			row := make([]string, 0)
			row = append(row, z.Name)
			row = append(row, fmt.Sprintf("%d", z.Replicas))
			row = append(row, fmt.Sprintf("%d", z.Priority))
			row = append(row, strings.Join(z.Hosts, ","))
			data = append(data, row)
		}
		sort.Slice(data, func(i, j int) bool {
			return data[i][0] < data[j][0]
		})
		r, err := common.Format(header, data, common.OutputFormatType(output))
		if err != nil {
			return common.NgctlError("Show hosts failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, r)
		return nil
	},
}

func init() {
	ShowCmd.Flags().StringVarP(&output, "output", "o", "table", "output format. Allowed values: table, json")
}
