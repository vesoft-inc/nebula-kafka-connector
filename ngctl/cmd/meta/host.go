package main

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl"
)

type hostFlagsType struct {
	host        string
	clusterName string
	agentPort   uint32
}

var hostFlags hostFlagsType

var hostCmd = &cobra.Command{
	Use:   "host",
	Short: `Process host command`,
	Long:  `Execute host command in cli mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var addHostCmd = &cobra.Command{
	Use:   "add",
	Short: "Add host into cluster.",
	Long:  `ngctl host add --host [host] --cluster [clustername] --agent_port [agentPort]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := hostFlags
		if flags.host == "" {
			return metaConsoleError("host is empty", "")
		}
		if flags.clusterName == "" {
			return metaConsoleError("cluster name is empty", "")
		}
		if flags.agentPort == 0 {
			return metaConsoleError("agent port must be set and non-zero", "")
		}
		req := meta.NewAddHostReq(flags.host, flags.clusterName, flags.agentPort)
		if err := metaClient.AddHost(req); err != nil {
			return metaConsoleError("Add host failed", err.Error())
		}
		fmt.Fprintln(metaOutput, "Add host successfully.")
		return nil
	},
}

var dropHostCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop host from cluster.",
	Long:  `ngctl host drop --host [host] --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := hostFlags
		if flags.host == "" {
			return metaConsoleError("host is empty", "")
		}
		if flags.clusterName == "" {
			return metaConsoleError("cluster name is empty", "")
		}
		req := meta.NewDropHostReq(flags.host, flags.clusterName)
		if err := metaClient.DropHost(req); err != nil {
			return metaConsoleError("Drop host failed", err.Error())
		}
		fmt.Fprintln(metaOutput, "Drop host successfully.")
		return nil
	},
}

var listHostsCmd = &cobra.Command{
	Use:   "list",
	Short: "List hosts in cluster.",
	Long:  `ngctl host list --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster := hostFlags.clusterName
		if cluster == "" {
			return metaConsoleError("cluster name is empty", "")
		}
		req := meta.NewListHostsReq(cluster)
		resp, err := metaClient.ListHosts(req)
		if err != nil {
			return metaConsoleError("List hosts failed", err.Error())
		}
		header := []string{"Host", "Agent Port"}
		data := make([][]string, 0)
		for _, h := range resp.HostInfoList {
			row := make([]string, 0)
			row = append(row, h.HostName)
			row = append(row, fmt.Sprintf("%d", h.AgentPort))
			data = append(data, row)
		}
		// order by host names
		sort.Slice(data, func(i, j int) bool {
			return data[i][0] < data[j][0]
		})
		fmt.Fprintln(metaOutput, ngctl.FormatTable(header, data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(hostCmd)
	hostCmd.AddCommand(addHostCmd)
	hostCmd.AddCommand(dropHostCmd)
	hostCmd.AddCommand(listHostsCmd)
	addHostCmd.Flags().StringVarP(&hostFlags.host, "host", "H", "", "the host to be added to a cluster")
	addHostCmd.Flags().StringVarP(&hostFlags.clusterName, "cluster", "c", "", "cluster name")
	addHostCmd.Flags().Uint32VarP(&hostFlags.agentPort, "agent_port", "a", 0, "port of the agent on the host")
	dropHostCmd.Flags().StringVarP(&hostFlags.host, "host", "H", "", "the host to be dropped from a cluster")
	dropHostCmd.Flags().StringVarP(&hostFlags.clusterName, "cluster", "c", "", "cluster name")
	listHostsCmd.Flags().StringVarP(&hostFlags.clusterName, "cluster", "c", "", "cluster name")
}
