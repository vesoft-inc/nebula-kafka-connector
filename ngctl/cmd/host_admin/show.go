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
	Short: "Show hosts in cluster.",
	Long:  `Show all hosts currently in a cluster.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster := hostFlags.clusterName
		if cluster == "" {
			return common.NgctlError("cluster name is empty", "")
		}
		req := meta.NewShowHostsReq(cluster)
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
		fmt.Fprintln(common.MetaOutput, common.FormatTable(header, data))
		return nil
	},
}

func init() {
	// show all the hosts of a cluster
	showHostsCmd.Flags().StringVarP(&hostFlags.clusterName, "cluster", "c", "", "cluster name")
}
