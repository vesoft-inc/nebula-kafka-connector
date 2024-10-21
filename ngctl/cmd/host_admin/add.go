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
	Short: "Add a host into a svcgrp.",
	Long: `Add a host into a svcgrp. A host is identified by its IP address. The port of the deployed agent is also needed.

Either provides the config file, all hosts in the config file will be added into the svcgrp.
Or provides the host info, the host will be added into the svcgrp.
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		common.MetaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateAddFlags(); err != nil {
			return err
		}
		var (
			hostResrouces *common.ResourceInfo
			err           error
		)
		if hostFlags.configFile != "" {
			hostResrouces, err = getHostsWithConfig()
			if err != nil {
				return err
			}
			*hostResrouces, err = common.ConfirmResourceList(*hostResrouces)
			if err != nil {
				return err
			}

		} else {
			hostResrouces, err = getHostDirectly()
		}
		if err != nil {
			return err
		}

		for _, host := range hostResrouces.ResourceList {
			agentPort, err := strconv.Atoi(host.AgentPort)
			if err != nil {
				return err
			}
			req := meta.NewAddHostReq(host.IP, hostFlags.svcgrpName, uint32(agentPort))
			if err := common.MetaClient.AddHost(req); err != nil {
				return common.NgctlError("Add host failed", err.Error())
			}
		}
		fmt.Fprintln(common.MetaOutput, "Add hosts successfully.")
		return nil
	},
}

func getHostDirectly() (*common.ResourceInfo, error) {
	var flags = hostFlags
	hostResrouces := common.ResourceInfo{
		ResourceType:        "hosts",
		OperationOnResource: "add",
		ResourceList:        make([]common.IPAndPort, 0),
		ClusterName:         flags.svcgrpName,
	}
	hostResrouces.ResourceList = append(
		hostResrouces.ResourceList,
		common.IPAndPort{IP: flags.host, AgentPort: strconv.Itoa(int(flags.agentPort))},
	)
	return &hostResrouces, nil
}

func getHostsWithConfig() (*common.ResourceInfo, error) {
	var flags = hostFlags
	hostResrouces := common.ResourceInfo{
		ResourceType:        "hosts",
		OperationOnResource: "add",
		ResourceList:        make([]common.IPAndPort, 0),
		ClusterName:         flags.svcgrpName,
	}
	hostList, err := common.DeriveHostList("", flags.svcgrpName, false)
	if err != nil {
		return nil, common.NgctlError("Failed to derive host list", err.Error())
	}
	for _, host := range hostList {
		hostResrouces.ResourceList = append(hostResrouces.ResourceList, host)
	}
	return &hostResrouces, nil

}

func validateAddFlags() error {
	var flags = hostFlags
	if flags.svcgrpName == "" {
		return common.NgctlError("svcgrp name is empty", "")
	}
	if flags.configFile == "" {
		if flags.host == "" {
			return common.NgctlError("must provide host info", "")
		}
	} else {
		if flags.host != "" {
			return common.NgctlError("cannot use host and config file at the same time", "")
		}
		configError := common.CheckInConfigFile(flags.configFile)
		if configError != nil {
			return common.NgctlError("Error in config file", configError.Error())
		}
	}
	return nil
}

func init() {
	// add host
	addHostCmd.Flags().StringVarP(&hostFlags.host, "host", "H", "", "the host to be added to a svcgrp")
	addHostCmd.Flags().StringVarP(&hostFlags.svcgrpName, "svcgrp", "s", "", "svcgrp name")
	addHostCmd.Flags().Uint32VarP(&hostFlags.agentPort, "agent_port", "a", 6688, "agent port")
	addHostCmd.Flags().StringVarP(&hostFlags.configFile, "config", "f", "", "config file")
}
