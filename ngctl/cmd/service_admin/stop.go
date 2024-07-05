package service_admin

import (
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var stopServiceCmd = &cobra.Command{
	Use:   "stop",
	Short: `Stop a service on a host`,
	Long:  `Stop a service on a host`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := ServiceFlags
		if flags.serviceType == "" || (flags.serviceType != "graphd" && flags.serviceType != "storaged") {
			return common.NgctlError("service type is wrong", "")
		}
		err := common.CheckInConfigFile(flags.configFile)
		if err != nil {
			return common.NgctlError("Failed to parse config file for the install path", err.Error())
		}
		if flags.host == "" || common.IsValidIPAddress(flags.host) == false{
			return common.NgctlError("no valid host provided", "")
		}
		agent, err := common.GetAgentForHost(flags.host)
		if err != nil {
			return common.NgctlError("Failed to get the agent for the service", err.Error())
		}
		err = ServiceOperation(agent, flags.serviceType, common.ConfigSpec.InstallPath, "stop")
		if err != nil {
			return common.NgctlError("Failed to start service", err.Error())
		}
		return nil
	},
}

func init() {
	stopServiceCmd.Flags().StringVarP(&ServiceFlags.serviceType, "type", "t", "", "service type")
	stopServiceCmd.Flags().StringVarP(&ServiceFlags.host, "host", "H", "", "host")
	stopServiceCmd.Flags().Uint32VarP(&ServiceFlags.port, "port", "P", 0, "port")
	stopServiceCmd.Flags().StringVarP(&ServiceFlags.configFile, "config", "f", "", "config file path")
	stopServiceCmd.MarkFlagsRequiredTogether("type", "host", "port", "config")
}
