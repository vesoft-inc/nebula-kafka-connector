package host_admin

import (
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

type hostFlagsType struct {
	host       string
	svcgrpName string
	agentPort  uint32
	configFile string
	output     string
	drain      bool
}

var hostFlags hostFlagsType

var HostCmd = &cobra.Command{
	Use:   "host",
	Short: `Run commands managing hosts`,
	Long:  `Run commands managing hosts`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func validateAddorDropFlags() error {
	var flags = hostFlags
	if flags.svcgrpName == "" {
		return common.NgctlError("svcgrp name is empty", "")
	}
	if flags.host == "" {
		return common.NgctlError("must provide host info", "")
	}
	return nil
}

func validateOperationFlags() error {
	var flags = hostFlags

	if flags.host == "" {
		return common.NgctlError("must provide host info", "")
	}
	if flags.configFile == "" {
		return common.NgctlError("config file is empty", "")
	}
	return nil
}

func init() {
	HostCmd.AddCommand(addHostCmd)
	HostCmd.AddCommand(dropHostCmd)
	HostCmd.AddCommand(showHostsCmd)

	HostCmd.AddCommand(installHostCmd)
	HostCmd.AddCommand(uninstallHostCmd)

}
