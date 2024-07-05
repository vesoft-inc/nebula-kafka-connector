package main

import (
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/cluster_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/host_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/service_admin"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/supercluster_admin"
)

var rootCmd = &cobra.Command{
	Use:   "ngctl",
	Short: "Execute ngctl commands in cli mode.",
	Long:  `Execute ngctl commands in cli mode. Use 'ngctl -h' to see usage.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	// cmds communicating with the metad to manage the supercluster
	rootCmd.AddCommand(supercluster_admin.SuperclusterCmd)
	// cmds communicating with the metad to manage a specific cluster
	rootCmd.AddCommand(cluster_admin.ClusterCmd)
	// cmds communicating with a agent to mange hosts
	rootCmd.AddCommand(host_admin.HostCmd)
	// cmds coomunicating with a agent to manage services on a host
	rootCmd.AddCommand(service_admin.ServiceAdminCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
