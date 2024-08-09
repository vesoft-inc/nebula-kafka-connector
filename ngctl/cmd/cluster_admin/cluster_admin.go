package cluster_admin

import (
	"github.com/spf13/cobra"
)

type clusterFlagsType struct {
	clusterName string
	replicas    int
	force       bool
	owner       string
	output      string
}

var clusterFlags clusterFlagsType

var ClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Run commands managing a cluster",
	Long:  `A cluster may contain multiple graphd and storaged services. Users could run commands to create, initialize, alternate, show, backup or drop a cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ClusterCmd.PersistentFlags().StringVarP(&clusterFlags.clusterName, "cluster", "c", "", "Cluster name")
	ClusterCmd.AddCommand(createClusterCmd)
	ClusterCmd.AddCommand(initClusterCmd)
	ClusterCmd.AddCommand(showClusterCmd)
	showClusterCmd.Flags().StringVarP(&clusterFlags.output, "output", "o", "table", "output format. Allowed values: table, json")

	ClusterCmd.AddCommand(brCmd)
	createClusterCmd.Flags().IntVarP(&clusterFlags.replicas, "replica-factor", "r", 3, "replica number, default: 3")
	createClusterCmd.Flags().StringVarP(&clusterFlags.owner, "owner", "o", "", "cluster owner")

	ClusterCmd.AddCommand(dropClusterCmd)
	dropClusterCmd.Flags().BoolVarP(&clusterFlags.force, "force", "f", false, "force drop cluster")

	ClusterCmd.AddCommand(alterClusterCmd)
	alterClusterCmd.Flags().StringVarP(&clusterFlags.owner, "owner", "o", "", "cluster owner")
}
