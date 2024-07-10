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
}

var ServiceFlags ServiceFlagsType

var ServiceAdminCmd = &cobra.Command{
	Use:   "service",
	Short: "Process service command",
	Long:  `Execute service command in cli mode.`,
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
	ServiceAdminCmd.AddCommand(dropServiceCmd)
	ServiceAdminCmd.AddCommand(startServiceCmd)
	ServiceAdminCmd.AddCommand(stopServiceCmd)
}
