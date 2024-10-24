package svcgrp_admin

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var showsvcgrpCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of a service group",
	Long:  "Show details of a service group",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		svcgrp := svcgrpFlags.svcgrpName
		req := meta.NewListServiceGroupsReq(svcgrp)
		resp, err := common.MetaClient.ListServiceGroups(req)
		if err != nil {
			return common.NgctlError("Show svcgrp failed", err.Error())
		}

		header := []string{"Id", "Name", "Replica", "Owner"}
		data := make([][]string, 0)
		for _, s := range resp.ServiceGroups {
			row := make([]string, 0)
			row = append(row, fmt.Sprintf("%d", s.Id))
			row = append(row, fmt.Sprintf("%s", s.Name))
			row = append(row, fmt.Sprintf("%d", s.ReplicaRefactor))
			row = append(row, fmt.Sprintf("%s", s.Owner))
			data = append(data, row)
		}

		// order by svcgrp id
		sort.Slice(data, func(i, j int) bool {
			return data[i][0] < data[j][0]
		})
		// printer.FormatTable(headers []string, data [][]string)
		r, err := common.Format(header, data, common.OutputFormatType(svcgrpFlags.output))
		if err != nil {
			return common.NgctlError("Show svcgrp failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, r)
		return nil
	},
}

func init() {
	showsvcgrpCmd.Flags().StringVarP(&svcgrpFlags.output, "output", "o", "table", "output format. Allowed values: table, json")
}
