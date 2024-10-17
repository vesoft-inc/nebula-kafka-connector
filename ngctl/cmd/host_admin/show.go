package host_admin

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var showHostsCmd = &cobra.Command{
	Use:   "show",
	Short: "Show hosts in srvgrp.",
	Long:  `Show all hosts currently in a srvgrp.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		srvgrp := hostFlags.srvgrpName
		if srvgrp == "" {
			return common.NgctlError("srvgrp name is empty", "")
		}
		req := meta.NewShowHostsReq(srvgrp)
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
	// show all the hosts of a srvgrp
	showHostsCmd.Flags().StringVarP(&hostFlags.srvgrpName, "srvgrp", "s", "", "srvgrp name")
}
