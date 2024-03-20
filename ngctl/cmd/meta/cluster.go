package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngql/printer"
)

type clusterFlagsType struct {
	clusterName string
	replicas    int
	zones       []string
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
	Long:  `ngctl cluster create --cluster [clustername] --replica [replica] --zones [zone1,zone2,...]`,
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
		req := meta.NewCreateClusterReq(clusterFlags.clusterName, clusterFlags.replicas, clusterFlags.zones)
		resp, err := metaClient.CreateCluster(req)
		if err != nil {
			return metaConsoleError("Create cluster failed", err.Error())
		}
		if !resp.IsSucceeded() {
			return metaConsoleError("Create cluster failed", resp.GetErrorMsg())
		}
		fmt.Fprintln(metaOutput, "Create cluster successfully.")
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

		resp, err := metaClient.InitCluster(req)
		if err != nil {
			return metaConsoleError("Init cluster failed", err.Error())
		}
		if !resp.IsSucceeded() {
			return metaConsoleError("Init cluster failed", resp.GetErrorMsg())
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
		req := meta.NewShowClusterReq(cluster)
		resp, err := metaClient.ShowCluster(req)
		if err != nil {
			return metaConsoleError("Show cluster failed", err.Error())
		}
		if !resp.IsSucceeded() {
			return metaConsoleError("Show cluster failed", resp.GetErrorMsg())
		}
		header := []string{"cluster id", "cluster name", "replica", "zones"}
		data := make([][]string, 0)
		for _, s := range resp.Clusters {
			row := make([]string, 0)
			row = append(row, fmt.Sprintf("%d", s.ClusterId))
			row = append(row, fmt.Sprintf("%s", s.ClusterName))
			row = append(row, fmt.Sprintf("%d", s.Replica))
			row = append(row, fmt.Sprintf("%s", strings.Join(s.Zones, ",")))
			data = append(data, row)
		}

		// order by cluster id
		sort.Slice(data, func(i, j int) bool {
			return data[i][0] < data[j][0]
		})
		// printer.FormatTable(headers []string, data [][]string)
		fmt.Fprintln(metaOutput, printer.FormatTable(header, data))
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
	createClusterCmd.Flags().StringArrayVarP(&clusterFlags.zones, "zones", "z", []string{}, "zones")
}
