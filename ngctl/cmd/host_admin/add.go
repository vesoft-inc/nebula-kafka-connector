package host_admin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

func getValidAgentPort(hostIP string) (uint32, error) {
	// try to get the agent port from the config file
	agent, err := common.GetAgentForHost(hostIP)
	if err != nil {
		return 0, common.NgctlError("cannot find a valid agent port from either the config file or the commdn line options", "")
	}
	var u64 uint64
	u64, err = strconv.ParseUint(strings.Split(agent.Host, ":")[1], 10, 32)
	if err != nil {
		return 0, common.NgctlError("failed to acquire a valid agent port", "")
	}
	return uint32(u64), nil
}

var addHostCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a host into a cluster.",
	Long:  `Add a host into a cluster, with the NebulaGraph pkg installed`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := hostFlags
		if flags.clusterName == "" {
			return common.NgctlError("cluster name is empty", "")
		}
		configError := common.CheckInConfigFile(flags.configFile)
		hostList, err := common.DeriveHostList(flags.host, flags.clusterName)
		if err != nil {
			return common.NgctlError("Failed to derive host list", err.Error())
		}
		// Requesting the meta to add the specified host
		for _, host := range hostList {
			agentPort, err := strconv.Atoi(host.AgentPort)
			if err != nil {
				return common.NgctlError("failed to get the agent port for host "+host.IP, err.Error())
			}
			req := meta.NewAddHostReq(host.IP, flags.clusterName, uint32(agentPort))
			if err := common.MetaClient.AddHost(req); err != nil {
				return common.NgctlError("Add host failed", err.Error())
			}
			fmt.Fprintln(common.MetaOutput, "Add host successfully.")
		}
		// Install
		if flags.withInstall {
			if configError != nil {
				return common.NgctlError("Failed to get a valid config file needed to install on the host", configError.Error())
			}
			if err := InstallOnHost(false); err != nil {
				return common.NgctlError("Failed to install on the host", err.Error())
			}
		}
		return nil
	},
}

func init() {
	// add host
	addHostCmd.Flags().StringVarP(&hostFlags.host, "host", "H", "", "the host to be added to a cluster")
	addHostCmd.Flags().StringVarP(&hostFlags.clusterName, "cluster", "c", "", "cluster name")
	addHostCmd.Flags().StringVarP(&hostFlags.configFile, "config", "f", "", "config file")
	// optinoal, install NebulaGraph on the host or not
	addHostCmd.Flags().
		BoolVarP(&hostFlags.withInstall, "with_install", "w", false, "install NebulaGraph on the host or not")

	addHostCmd.MarkFlagRequired("cluster")
	addHostCmd.MarkFlagRequired("config")
}
