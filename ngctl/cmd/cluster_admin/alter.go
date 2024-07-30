package cluster_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var alterClusterCmd = &cobra.Command{
	Use:   "alter",
	Short: "Alter the owner of a cluster.",
	Long:  `A cluster has an owner registered in the supercluster. Users can alter it.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if clusterFlags.clusterName == "" {
			return common.NgctlError("Cluster name is empty", "")
		}
		if clusterFlags.owner == "" {
			return common.NgctlError("Owner is invalid", "")
		}
		// donot need to specify zones
		req := meta.NewAlterClusterReq(clusterFlags.clusterName, clusterFlags.owner)
		if err := common.MetaClient.AlterCluster(req); err != nil {
			return common.NgctlError("Alter cluster failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Alter cluster successfully.")
		return nil
	},
}
