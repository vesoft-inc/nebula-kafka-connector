package cluster_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var createClusterCmd = &cobra.Command{
	Use:   "create",
	Short: "Create cluster in meta server.",
	Long:  `ngctl cluster create --cluster [clustername] --replica-factor [replica-factor] --owner [owner]`,
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
		if clusterFlags.replicas == 0 {
			return common.NgctlError("Replica number is invalid", "")
		}
		// donot need to specify zones
		req := meta.NewCreateClusterReq(clusterFlags.clusterName, clusterFlags.replicas,
			clusterFlags.owner, nil)
		if err := common.MetaClient.CreateCluster(req); err != nil {
			return common.NgctlError("Create cluster failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Create cluster successfully.")
		return nil
	},
}
