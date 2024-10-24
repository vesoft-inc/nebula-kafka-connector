package metad_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/job"
)

var configMetadCmd = &cobra.Command{
	Use:   "config",
	Short: "Config metad",
	Long:  "Config metad",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateOperationFlags(); err != nil {
			return err
		}
		c, err := config.NewConfigFromFile(metadFlags.configFile)
		if err != nil {
			return common.NgctlError("Failed to check in config file", err.Error())
		}
		instance, err := c.GetServiceInstances("", "metad", metadFlags.host)
		if err != nil {
			return err
		}
		hosts := make([]string, 0)
		for _, inst := range instance {
			hosts = append(hosts, inst.Host)
		}
		if err := job.NgadmJob.UpdateConfig(c, "", hosts, "metad"); err != nil {
			return err
		}
		fmt.Fprintf(common.MetaOutput, "Config metad successfully.\n")
		return nil
	},
}

func init() {
	configMetadCmd.Flags().StringVarP(&metadFlags.host, "host", "H", "", "host")
	configMetadCmd.Flags().StringVarP(&metadFlags.configFile, "config", "f", "", "config file for ngctl")
}
