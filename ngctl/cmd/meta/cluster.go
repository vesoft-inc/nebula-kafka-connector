package main

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl"
)

type clusterFlagsType struct {
	clusterName string
	replicas    int
	force       bool
	owner       string
}

var clusterFlags clusterFlagsType

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Process cluster command",
	Long:  `Execute cluster command in cli mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var createClusterCmd = &cobra.Command{
	Use:   "create",
	Short: "Create cluster in meta server.",
	Long:  `ngctl cluster create --cluster [clustername] --replica [replica] --owner [owner]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if clusterFlags.clusterName == "" {
			return metaConsoleError("Cluster name is empty", "")
		}
		if clusterFlags.replicas == 0 {
			return metaConsoleError("Replica number is invalid", "")
		}
		// donot need to specify zones
		req := meta.NewCreateClusterReq(clusterFlags.clusterName, clusterFlags.replicas,
			clusterFlags.owner, nil)
		if err := metaClient.CreateCluster(req); err != nil {
			return metaConsoleError("Create cluster failed", err.Error())
		}
		fmt.Fprintln(metaOutput, "Create cluster successfully.")
		return nil
	},
}

var alterClusterCmd = &cobra.Command{
	Use:   "alter",
	Short: "Alter cluster in meta server.",
	Long:  `ngctl cluster alter --cluster [clustername] --owner [owner]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if clusterFlags.clusterName == "" {
			return metaConsoleError("Cluster name is empty", "")
		}
		if clusterFlags.owner == "" {
			return metaConsoleError("Owner is invalid", "")
		}
		// donot need to specify zones
		req := meta.NewAlterClusterReq(clusterFlags.clusterName, clusterFlags.owner)
		if err := metaClient.AlterCluster(req); err != nil {
			return metaConsoleError("Create cluster failed", err.Error())
		}
		fmt.Fprintln(metaOutput, "Alter cluster successfully.")
		return nil
	},
}

var initClusterCmd = &cobra.Command{
	Use:   "init",
	Short: "Init cluster storage part.",
	Long:  `ngctl cluster init --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster := clusterFlags.clusterName
		if cluster == "" {
			return metaConsoleError("cluster name is empty", "")
		}
		req := meta.NewInitClusterReq(cluster)

		if err := metaClient.InitCluster(req); err != nil {
			return metaConsoleError("Init cluster failed", err.Error())
		}
		fmt.Fprintln(metaOutput, "Init cluster successfully.")
		return nil
	},
}

var showClusterCmd = &cobra.Command{
	Use:   "show",
	Short: "Show cluster, show all if no cluster name specified.",
	Long:  `nebula-meta cluster show --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster := clusterFlags.clusterName
		req := meta.NewListClustersReq(cluster)
		resp, err := metaClient.ListClusters(req)
		if err != nil {
			return metaConsoleError("Show cluster failed", err.Error())
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
		fmt.Fprintln(metaOutput, ngctl.FormatTable(header, data))
		return nil
	},
}

var dropClusterCmd = &cobra.Command{
	Use:   "drop",
	Short: "drop cluster storage part.",
	Long:  `ngctl cluster drop --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster := clusterFlags.clusterName
		force := clusterFlags.force
		if cluster == "" {
			return metaConsoleError("cluster name is empty", "")
		}

		req := meta.NewDropClusterReq(cluster, force)
		if err := metaClient.DropCluster(req); err != nil {
			return metaConsoleError("Init cluster failed", err.Error())
		}
		fmt.Fprintln(metaOutput, "Drop cluster successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.PersistentFlags().StringVarP(&clusterFlags.clusterName, "cluster", "c", "", "Cluster name")
	clusterCmd.AddCommand(createClusterCmd)
	clusterCmd.AddCommand(initClusterCmd)
	clusterCmd.AddCommand(showClusterCmd)
	createClusterCmd.Flags().IntVarP(&clusterFlags.replicas, "replica-factor", "r", 3, "replica number, default: 3")
	createClusterCmd.Flags().StringVarP(&clusterFlags.owner, "owner", "o", "", "cluster owner")

	clusterCmd.AddCommand(dropClusterCmd)
	dropClusterCmd.Flags().BoolVarP(&clusterFlags.force, "force", "f", false, "force drop cluster")

	clusterCmd.AddCommand(alterClusterCmd)
	alterClusterCmd.Flags().StringVarP(&clusterFlags.owner, "owner", "o", "", "cluster owner")
}
