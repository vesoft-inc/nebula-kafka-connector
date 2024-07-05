package cluster_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var dropClusterCmd = &cobra.Command{
	Use:   "drop",
	Short: "drop cluster storage part.",
	Long:  `ngctl cluster drop --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster := clusterFlags.clusterName
		force := clusterFlags.force
		if cluster == "" {
			return common.NgctlError("cluster name is empty", "")
		}

		req := meta.NewDropClusterReq(cluster, force)
		if err := common.MetaClient.DropCluster(req); err != nil {
			return common.NgctlError("Init cluster failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Drop cluster successfully.")
		return nil
	},
}
