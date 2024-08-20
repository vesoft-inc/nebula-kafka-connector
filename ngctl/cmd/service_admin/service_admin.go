package service_admin

import (
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"

	"github.com/spf13/cobra"
)

type ServiceFlagsType struct {
	serviceType       string
	host              string
	port              int32
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

// validate flags for service commands, i.e. add, drop
// if configFile is provided, then host, port, serviceType should not be provided
func validateServiceFlags() error {
	var flags = ServiceFlags
	if flags.clusterName == "" {
		return common.NgctlError("cluster name is empty", "")
	}
	if flags.configFile == "" {
		if flags.host == "" || flags.port < 0 || flags.serviceType == "" {
			return common.NgctlError("must provide service info [host, port, type]", "")
		}
	} else {
		if flags.host != "" || flags.port >= 0 || flags.serviceType != "" {
			return common.NgctlError("cannot use service info and config file at the same time", "")
		}
		configError := common.CheckInConfigFile(flags.configFile)
		if configError != nil {
			return common.NgctlError("Error in config file", configError.Error())
		}
	}
	return nil
}

// validate flags for service operation commands, i.e. start, stop, restart
// must provides configFile, if user also provides host, port, serviceType,
// just filter the provided information.
func validateServiceOperationFlags() error {
	var flags = ServiceFlags
	if flags.configFile == "" {
		return common.NgctlError("must provide configFile", "")
	}
	return nil
}
