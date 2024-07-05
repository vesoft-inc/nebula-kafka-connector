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
		if flags.host == "" || common.IsValidIPAddress(flags.host) == false{
			return common.NgctlError("no valid host provided", "")
		}
		if flags.clusterName == "" {
			return common.NgctlError("cluster name is empty", "")
		}
		configError := common.CheckInConfigFile(flags.configFile)
		if flags.agentPort == 0 {
			if configError == nil {
				port, err := getValidAgentPort(flags.host)
				if port == 0 || err != nil {
					return common.NgctlError("agent port is empty", err.Error())
				} else {
					flags.agentPort = port
				}
			} else {
				return common.NgctlError("cannot find a valid agent port from either the config file or the commdn line options", "")
			}
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
		// Requesting the meta to add the specified host
		req := meta.NewAddHostReq(flags.host, flags.clusterName, flags.agentPort)
		if err := common.MetaClient.AddHost(req); err != nil {
			return common.NgctlError("Add host failed", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Add host successfully.")
		return nil
	},
}

func init() {
	// add host
	addHostCmd.Flags().StringVarP(&hostFlags.host, "host", "H", "", "the host to be added to a cluster")
	addHostCmd.Flags().StringVarP(&hostFlags.clusterName, "cluster", "c", "", "cluster name")
	// an agent port needs to be set by this option, or be configed in the config file, otherwise an error would be reported.
	addHostCmd.Flags().Uint32VarP(&hostFlags.agentPort, "agent_port", "a", 0, "port of the agent on the host")
	// a config file is needed when no agent port is set in the option list or an install is required.
	addHostCmd.Flags().StringVarP(&hostFlags.configFile, "config", "f", "", "config file")
	// optinoal, install NebulaGraph on the host or not
	addHostCmd.Flags().
		BoolVarP(&hostFlags.withInstall, "with_install", "w", false, "install NebulaGraph on the host or not")
}
