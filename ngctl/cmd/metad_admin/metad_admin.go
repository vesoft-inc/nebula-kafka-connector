package metad_admin

import (
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

type metadFlagsType struct {
	host       string
	configFile string
	output     string
}

var metadFlags metadFlagsType

var MetadCmd = &cobra.Command{
	Use:   "metad",
	Short: "Run commands managing a metad",
	Long:  "Run commands managing a metad",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func validateOperationFlags() error {
	var flags = metadFlags

	if flags.configFile == "" {
		return common.NgctlError("config file is empty", "")
	}
	return nil
}

func init() {
	MetadCmd.AddCommand(startMetadCmd)
	MetadCmd.AddCommand(stopMetadCmd)
	MetadCmd.AddCommand(configMetadCmd)
	MetadCmd.AddCommand(showMetadCmd)
}
