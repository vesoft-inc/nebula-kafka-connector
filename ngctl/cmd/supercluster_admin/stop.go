package supercluster_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/service_admin"
)

// By default, we stop all the metad services in a supercluster.
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a supercluster",
	Long:  `Stop a supercluster without deleting its data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := common.CheckInConfigFile(superclusterFlags.configFile)
		if err != nil {
			return common.NgctlError("Failed to check in config file", err.Error())
		}
		// TODO(Xuntao): check the status of graphd/storaged services in the supercluster
		for _, metad_host := range common.ConfigSpec.Spec.Metad.Hosts {
			agent_host := metad_host.Agent.Host
			path := common.ConfigSpec.InstallPath
			err = service_admin.ServiceOperation(types.Agent{Host: agent_host}, "metad", path, "stop")
			if err != nil {
				return common.NgctlError(fmt.Sprintf("Failed to stop the metad service at %s:%s", metad_host.IP, metad_host.Port), err.Error())
			}
		}

		return nil
	},
}

func init() {
	stopCmd.Flags().StringVarP(&superclusterFlags.configFile, "config", "f", "", "The config file used to create the supercluster")
}
