package metad_admin

import "github.com/spf13/cobra"

type metadFlagsType struct {
	host              string
	withInstall       bool
	configFile        string
	serviceConfigFile string
}

var metadFlags metadFlagsType

var SupermetadCmd = &cobra.Command{
	Use:   "metad",
	Short: "Run commands managing a metad.",
	Long:  `Run commands managing a metad.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	SupermetadCmd.AddCommand(createCmd)
	SupermetadCmd.AddCommand(loginCmd)
	SupermetadCmd.AddCommand(logoutCmd)
	SupermetadCmd.AddCommand(passwdCmd)
	SupermetadCmd.AddCommand(userCmd)
	SupermetadCmd.AddCommand(stopCmd)
}
