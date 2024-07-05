package host_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

var dropHostCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop host from cluster.",
	Long:  `ngctl host drop --host [host] --cluster [clustername]`,
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
		fmt.Println("Dropping host", flags.host, "from cluster", flags.clusterName)
		req := meta.NewDropHostReq(flags.host, flags.clusterName)
		if err := common.MetaClient.DropHost(req); err != nil {
			return common.NgctlError("Drop host failed", err.Error())
		}
		// Uninstall
		if flags.withInstall {
			if err := common.CheckInConfigFile(flags.configFile); err != nil {
				return common.NgctlError("Failed to get a valid config file", err.Error())
			}
			if err := UninstallOnHost(); err != nil {
				return common.NgctlError("Uninstall on host failed", err.Error())
			}
		}
		fmt.Fprintln(common.MetaOutput, "Drop host successfully.")
		return nil
	},
}

func init() {
	// drop host
	// the option list is similar to the above
	dropHostCmd.Flags().StringVarP(&hostFlags.host, "host", "H", "", "the host to be dropped from a cluster")
	dropHostCmd.Flags().StringVarP(&hostFlags.clusterName, "cluster", "c", "", "cluster name")
	dropHostCmd.Flags().Uint32VarP(&hostFlags.agentPort, "agent_port", "a", 0, "port of the agent on the host")
	dropHostCmd.Flags().StringVarP(&hostFlags.configFile, "config", "f", "", "config file")
	dropHostCmd.Flags().
		BoolVarP(&hostFlags.withUninstall, "with_uninstall", "u", false, "uninstall NebulaGraph on the host or not")
}
