package cluster_admin

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var showClusterCmd = &cobra.Command{
	Use:   "show",
	Short: "Show cluster, show all if no cluster name specified.",
	Long:  `ngctl cluster show --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster := clusterFlags.clusterName
		req := meta.NewListClustersReq(cluster)
		resp, err := common.MetaClient.ListClusters(req)
		if err != nil {
			return common.NgctlError("Show cluster failed", err.Error())
		}

		header := []string{"Id", "Name", "Replica", "Owner"}
		data := make([][]string, 0)
		for _, s := range resp.Clusters {
			row := make([]string, 0)
			row = append(row, fmt.Sprintf("%d", s.Id))
			row = append(row, fmt.Sprintf("%s", s.Name))
			row = append(row, fmt.Sprintf("%d", s.ReplicaRefactor))
			row = append(row, fmt.Sprintf("%s", s.Owner))
			data = append(data, row)
		}

		// order by cluster id
		sort.Slice(data, func(i, j int) bool {
			return data[i][0] < data[j][0]
		})
		// printer.FormatTable(headers []string, data [][]string)
		fmt.Fprintln(common.MetaOutput, common.FormatTable(header, data))
		return nil
	},
}
