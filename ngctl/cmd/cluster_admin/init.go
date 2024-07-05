package cluster_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var initClusterCmd = &cobra.Command{
	Use:   "init",
	Short: "Init cluster storage part.",
	Long:  `ngctl cluster init --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster := clusterFlags.clusterName
		if cluster == "" {
			return common.NgctlError("cluster name is empty", "")
		}
		req := meta.NewInitClusterReq(cluster)

		if err := common.MetaClient.InitCluster(req); err != nil {
			return common.NgctlError("Init cluster failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Init cluster successfully.")
		return nil
	},
}
