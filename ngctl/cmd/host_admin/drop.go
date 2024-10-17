package host_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var dropHostCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop a host from a srvgrp.",
	Long:  `Drop a host from a srvgrp.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := hostFlags
		if flags.srvgrpName == "" {
			return common.NgctlError("srvgrp name is empty", "")
		}
		if err := common.CheckInConfigFile(flags.configFile); err != nil {
			return common.NgctlError("Failed to get a valid config file", err.Error())
		}
		hostList, err := common.DeriveHostList(flags.host, flags.srvgrpName, false)
		if err != nil {
			return common.NgctlError("Failed to derive host list", err.Error())
		}
		// Prepare the resource request
		hostResrouces := common.ResourceInfo{
			ResourceType:        "hosts",
			OperationOnResource: "drop",
			ResourceList:        make([]common.IPAndPort, 0),
			ClusterName:         flags.srvgrpName,
		}
		for _, host := range hostList {
			hostResrouces.ResourceList = append(hostResrouces.ResourceList, host)
		}
		if flags.host == "" {
			hostResrouces, err = common.ConfirmResourceList(hostResrouces)
			if err != nil {
				return common.NgctlError("Failed to confirm the host list", err.Error())
			}
		}
		for _, host := range hostResrouces.ResourceList {
			req := meta.NewDropHostReq(host.IP, flags.srvgrpName)
			if err := common.MetaClient.DropHost(req); err != nil {
				fmt.Fprintln(common.MetaOutput, fmt.Sprintf("Drop host %s failed: %s", host.IP, err.Error()))
			} else {
				fmt.Fprintln(common.MetaOutput, fmt.Sprintf("Drop host %s successfully.", host.IP))
			}
		}
		return nil
	},
}

func init() {
	// drop host
	// the option list is similar to the above
	dropHostCmd.Flags().StringVarP(&hostFlags.host, "host", "H", "", "the host to be dropped from a srvgrp")
	dropHostCmd.Flags().StringVarP(&hostFlags.srvgrpName, "srvgrp", "s", "", "srvgrp name")
	dropHostCmd.Flags().Uint32VarP(&hostFlags.agentPort, "agent_port", "a", 6688, "port of the agent on the host")
	dropHostCmd.Flags().StringVarP(&hostFlags.configFile, "config", "f", "", "config file")
}
