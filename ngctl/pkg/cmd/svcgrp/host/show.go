package host

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

var ShowHostsCmd = &cobra.Command{
	Use:   "show-host",
	Short: "Show hosts in svcgrp",
	Long:  "Show all hosts currently in a svcgrp",
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

		req := meta.NewShowHostsReq(svcgrpName)
		resp, err := common.MetaClient.ListHosts(req)
		if err != nil {
			return common.NgctlError("Show hosts failed", err.Error())
		}
		header := []string{"Host", "Agent Port"}
		data := make([][]string, 0)
		for _, h := range resp.HostInfoList {
			row := make([]string, 0)
			row = append(row, h.HostName)
			row = append(row, fmt.Sprintf("%d", h.AgentPort))
			data = append(data, row)
		}
		// order by host names
		sort.Slice(data, func(i, j int) bool {
			return data[i][0] < data[j][0]
		})
		r, err := common.Format(header, data, common.OutputFormatType(hostFlags.output))
		if err != nil {
			return common.NgctlError("Show hosts failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, r)
		return nil
	},
}

func init() {
	// show all the hosts of a svcgrp
	ShowHostsCmd.Flags().StringVarP(&hostFlags.output, "output", "o", "table", "output format. Allowed values: table, json")

}
