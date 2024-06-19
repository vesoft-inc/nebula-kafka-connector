package main

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl"
)

type serviceFlagsType struct {
	serviceType string
	host        string
	port        uint32
	clusterName string
}

var serviceFlags serviceFlagsType

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Process service command",
	Long:  `Execute service command in cli mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var addServiceCmd = &cobra.Command{
	Use:   "add",
	Short: `Add service into assigned cluster.`,
	Long:  `ngctl service add --type [graphd|storaged] --host [host] --port [port] --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := serviceFlags
		var serviceType meta.ServiceType
		if flags.serviceType == "graphd" {
			serviceType = meta.ServiceTypeGraphd
		} else if flags.serviceType == "storaged" {
			serviceType = meta.ServiceTypeStoraged
		} else {
			return metaConsoleError("service type is not correct, valid value is graphd or storaged", "")
		}
		if flags.host == "" {
			return metaConsoleError("service host is empty", "")
		}

		if flags.clusterName == "" {
			return metaConsoleError("cluster name is empty", "")
		}
		req := meta.NewAddServiceReq(flags.host, flags.port, serviceType, flags.clusterName)
		if err := metaClient.AddService(req); err != nil {
			return metaConsoleError("Add service failed", err.Error())
		}
		fmt.Fprintln(metaOutput, "Add service successfully.")
		return nil
	},
}

var dropServiceCmd = &cobra.Command{
	Use:   "drop",
	Short: `Drop service from assigned cluster.`,
	Long:  `ngctl service drop --type [graphd|storaged] --host [host] --port [port] --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := serviceFlags
		var serviceType meta.ServiceType
		if flags.serviceType == "graphd" {
			serviceType = meta.ServiceTypeGraphd
		} else if flags.serviceType == "storaged" {
			serviceType = meta.ServiceTypeStoraged
		} else {
			return metaConsoleError("service type is not correct, valid value is graphd or storaged", "")
		}
		if flags.host == "" {
			return metaConsoleError("service host is empty", "")
		}

		if flags.clusterName == "" {
			return metaConsoleError("cluster name is empty", "")
		}
		req := meta.NewDropServiceReq(flags.host, flags.port, serviceType, flags.clusterName)
		if err := metaClient.DropService(req); err != nil {
			return metaConsoleError("Drop service failed", err.Error())
		}
		fmt.Fprintln(metaOutput, "Drop service successfully.")
		return nil
	},
}

var showServiceCmd = &cobra.Command{
	Use:   "show",
	Short: "Show service in cluster.",
	Long:  `ngctl service show --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster := serviceFlags.clusterName
		if cluster == "" {
			return metaConsoleError("cluster name is empty", "")
		}
		req := meta.NewListServicesReq(cluster)

		resp, err := metaClient.ListServices(req)
		if err != nil {
			return metaConsoleError("Show service failed", err.Error())
		}

		header := []string{"Id", "Type", "Host", "Port"}
		data := make([][]string, 0)
		for _, s := range resp.Services {
			row := make([]string, 0)
			row = append(row, fmt.Sprintf("%d", s.Id))
			if s.Type == meta.ServiceTypeGraphd {
				row = append(row, "graphd")
			} else {
				row = append(row, "storaged")
			}
			row = append(row, s.Host)
			row = append(row, fmt.Sprintf("%d", s.Port))
			data = append(data, row)
		}

		// order by service id
		sort.Slice(data, func(i, j int) bool {
			return data[i][0] < data[j][0]
		})
		fmt.Fprintln(metaOutput, ngctl.FormatTable(header, data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.PersistentFlags().StringVarP(&serviceFlags.clusterName, "cluster", "c", "", "Cluster name")
	serviceCmd.AddCommand(addServiceCmd)
	serviceCmd.AddCommand(showServiceCmd)
	serviceCmd.AddCommand(dropServiceCmd)

	addServiceCmd.Flags().StringVarP(&serviceFlags.serviceType, "type", "t", "", "service type")
	addServiceCmd.Flags().StringVarP(&serviceFlags.host, "host", "H", "", "service host")
	addServiceCmd.Flags().Uint32VarP(&serviceFlags.port, "port", "P", 0, "service port")
	addServiceCmd.Flags().StringVarP(&serviceFlags.clusterName, "cluster", "c", "", "cluster name")
	addServiceCmd.MarkFlagsRequiredTogether("type", "host", "port", "cluster")

	dropServiceCmd.Flags().StringVarP(&serviceFlags.serviceType, "type", "t", "", "service type")
	dropServiceCmd.Flags().StringVarP(&serviceFlags.host, "host", "H", "", "service host")
	dropServiceCmd.Flags().Uint32VarP(&serviceFlags.port, "port", "P", 0, "service port")
	dropServiceCmd.Flags().StringVarP(&serviceFlags.clusterName, "cluster", "c", "", "cluster name")
	dropServiceCmd.MarkFlagsRequiredTogether("type", "host", "port", "cluster")
}
