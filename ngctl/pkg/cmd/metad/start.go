package metad

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/job"
)

var startMetadCmd = &cobra.Command{
	Use:   "start",
	Short: "Start metad",
	Long:  "Start metad",
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
		if err := job.NgadmJob.ServiceOperation(c, hosts, "metad", "start"); err != nil {
			return err
		}
		fmt.Fprintf(common.MetaOutput, "Start metad successfully.\n")
		return nil
	},
}

func init() {
	startMetadCmd.Flags().StringVarP(&metadFlags.host, "host", "H", "", "host")
	startMetadCmd.Flags().StringVarP(&metadFlags.configFile, "config", "f", "", "config file for ngctl")
}
