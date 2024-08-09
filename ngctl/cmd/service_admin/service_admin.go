package service_admin

import (
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"

	"github.com/spf13/cobra"
)

type ServiceFlagsType struct {
	serviceType       string
	host              string
	port              uint32
	clusterName       string
	configFile        string
	serviceConfigFile string
	agent             types.Agent
	output            string
}

var ServiceFlags ServiceFlagsType

var ServiceAdminCmd = &cobra.Command{
	Use:   "service",
	Short: "Run commands managing services.",
	Long:  `Run commands managing services.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	ServiceAdminCmd.AddCommand(addServiceCmd)
	ServiceAdminCmd.AddCommand(showServiceCmd)
	showServiceCmd.Flags().StringVarP(&ServiceFlags.output, "output", "o", "table", "output format. Allowed values: table, json")

	ServiceAdminCmd.AddCommand(dropServiceCmd)
	ServiceAdminCmd.AddCommand(startServiceCmd)
	ServiceAdminCmd.AddCommand(stopServiceCmd)
}
