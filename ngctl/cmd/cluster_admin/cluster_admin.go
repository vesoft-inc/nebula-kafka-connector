package cluster_admin

import (
	"github.com/spf13/cobra"
)

type clusterFlagsType struct {
	clusterName string
	replicas    int
	force       bool
	owner       string
}

var clusterFlags clusterFlagsType

var ClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Process cluster command",
	Long:  `Execute cluster command in cli mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ClusterCmd.PersistentFlags().StringVarP(&clusterFlags.clusterName, "cluster", "c", "", "Cluster name")
	ClusterCmd.AddCommand(createClusterCmd)
	ClusterCmd.AddCommand(initClusterCmd)
	ClusterCmd.AddCommand(showClusterCmd)
	ClusterCmd.AddCommand(brCmd)
	createClusterCmd.Flags().IntVarP(&clusterFlags.replicas, "replica-factor", "r", 3, "replica number, default: 3")
	createClusterCmd.Flags().StringVarP(&clusterFlags.owner, "owner", "o", "", "cluster owner")

	ClusterCmd.AddCommand(dropClusterCmd)
	dropClusterCmd.Flags().BoolVarP(&clusterFlags.force, "force", "f", false, "force drop cluster")

	ClusterCmd.AddCommand(alterClusterCmd)
	alterClusterCmd.Flags().StringVarP(&clusterFlags.owner, "owner", "o", "", "cluster owner")
}
