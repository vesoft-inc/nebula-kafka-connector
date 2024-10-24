package host_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/job"
)

var installHostCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a NebulaGraph package on a host",
	Long:  "Install a NebulaGraph package on a host",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateOperationFlags(); err != nil {
			return err
		}
		c, err := config.NewConfigFromFile(hostFlags.configFile)
		if err != nil {
			return common.NgctlError("Failed to check in config file", err.Error())
		}
		instance, err := c.GetInstanceFromHost(hostFlags.host)
		if err != nil {
			return common.NgctlError("Must config the host in the config file", err.Error())
		}
		if err := job.NgadmJob.InstallHost(c, instance.Host); err != nil {
			return err
		}
		fmt.Fprintf(common.MetaOutput, "Install NebulaGraph on host successfully.\n")
		return nil
	},
}

func init() {
	installHostCmd.Flags().StringVarP(&hostFlags.host, "host", "H", "", "on which host to install the NebulaGraph package")
	installHostCmd.Flags().StringVarP(&hostFlags.configFile, "config", "f", "", "config file")
}
