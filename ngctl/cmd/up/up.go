package up

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/job"
)

type upFlagsType struct {
	configFile      string
	defaultPassword string
}

var upFlags upFlagsType

var UpCmd = &cobra.Command{
	Use:   "up",
	Short: "Up a cluster, including metad and service groups",
	Long:  "Up a cluster, including metad and service groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := config.NewConfigFromFile(upFlags.configFile)
		if err != nil {
			return err
		}
		if err := job.NgadmJob.UpCluster(c, upFlags.defaultPassword); err != nil {
			return err
		}

		fmt.Fprintf(common.MetaOutput, "Up cluster successfully.\n")
		return nil
	},
}

func init() {
	UpCmd.Flags().StringVarP(&upFlags.configFile, "config", "f", "", "config file")
	UpCmd.Flags().StringVarP(&upFlags.defaultPassword, "password", "p", "", "default metad password after install")
}
