package down

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/job"
)

type downFlagsType struct {
	configFile string
	drain      bool
}

var downFlags downFlagsType

var DownCmd = &cobra.Command{
	Use:   "down",
	Short: "Down a cluster",
	Long:  "Down a cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := config.NewConfigFromFile(downFlags.configFile)
		if err != nil {
			return err
		}
		if err := job.NgadmJob.DownCluster(c, downFlags.drain); err != nil {
			return err
		}
		fmt.Fprintf(common.MetaOutput, "Down cluster successfully.\n")
		return nil
	},
}

func init() {
	DownCmd.Flags().StringVarP(&downFlags.configFile, "config", "f", "", "config file")
	DownCmd.Flags().BoolVarP(&downFlags.drain, "drain", "d", false, "whether to delete the data, default is false")
}
