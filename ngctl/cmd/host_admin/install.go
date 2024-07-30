package host_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var installHostCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a NebulaGraph package on a host.",
	Long:  `Install a NebulaGraph package on a host.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := hostFlags
		configError := common.CheckInConfigFile(flags.configFile)
		if configError != nil {
			return common.NgctlError("config file error", configError.Error())
		}
		// install on metad is managed by supercluster_admin, not managed here.
		hostList, err := common.DeriveHostList(flags.host, flags.clusterName, false)
		if err != nil {
			return common.NgctlError("Failed to derive host list", err.Error())
		}
		hostResrouces := common.ResourceInfo{
			ResourceType:        "hosts",
			OperationOnResource: "add",
			ResourceList:        make([]common.IPAndPort, 0),
			ClusterName:         flags.clusterName,
		}
		for _, host := range hostList {
			hostResrouces.ResourceList = append(hostResrouces.ResourceList, host)
		}
		if flags.host == "" {
			hostResrouces, err = common.ConfirmResourceList(hostResrouces)
			if err != nil {
				return common.NgctlError("Failed to confirm host list", err.Error())
			}
		}
		if err = InstallOnHost(hostResrouces.ResourceList, false); err != nil {
			return common.NgctlError("Failed to install NebulaGraph on the host", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Install NebulaGraph on hosts successfully.")
		return nil
	},
}

var uninstallHostCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall an NebulaGraph package on a host.",
	Long:  `Uninstall an NebulaGraph package on a host.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := hostFlags
		configError := common.CheckInConfigFile(flags.configFile)
		if configError != nil {
			return common.NgctlError("config file error", configError.Error())
		}
		// install on metad is managed by supercluster_admin, not managed here.
		hostList, err := common.DeriveHostList(flags.host, flags.clusterName, false)
		if err != nil {
			return common.NgctlError("Failed to derive host list", err.Error())
		}
		hostResrouces := common.ResourceInfo{
			ResourceType:        "hosts",
			OperationOnResource: "add",
			ResourceList:        make([]common.IPAndPort, 0),
			ClusterName:         flags.clusterName,
		}
		for _, host := range hostList {
			hostResrouces.ResourceList = append(hostResrouces.ResourceList, host)
		}
		if flags.host == "" {
			hostResrouces, err = common.ConfirmResourceList(hostResrouces)
			if err != nil {
				return common.NgctlError("Failed to confirm host list", err.Error())
			}
		}
		if err = UninstallOnHost(hostResrouces.ResourceList); err != nil {
			return common.NgctlError("Failed to uninstall NebulaGraph on the host", err.Error())
		}
		fmt.Fprintln(common.MetaOutput, "Uninstall NebulaGraph on hosts successfully.")
		return nil
	},
}

func init() {
	installHostCmd.Flags().StringVarP(&hostFlags.host, "host", "H", "", "on which host to install the NebulaGraph package")
	installHostCmd.Flags().StringVarP(&hostFlags.clusterName, "cluster", "c", "", "cluster name")
	installHostCmd.Flags().StringVarP(&hostFlags.configFile, "config", "f", "", "config file")

	uninstallHostCmd.Flags().StringVarP(&hostFlags.host, "host", "H", "", "on which host to uninstall the NebulaGraph package")
	uninstallHostCmd.Flags().StringVarP(&hostFlags.clusterName, "cluster", "c", "", "cluster name")
	uninstallHostCmd.Flags().StringVarP(&hostFlags.configFile, "config", "f", "", "config file")
}
