package supercluster_admin

import "github.com/spf13/cobra"

type superclusterFlagsType struct {
	host        string
	withInstall bool
	configFile  string
}

var superclusterFlags superclusterFlagsType

var SuperclusterCmd = &cobra.Command{
	Use:   "supercluster",
	Short: "Process supercluster command",
	Long:  `Execute supercluster command in cli mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	SuperclusterCmd.AddCommand(createCmd)
	SuperclusterCmd.AddCommand(loginCmd)
	SuperclusterCmd.AddCommand(logoutCmd)
	SuperclusterCmd.AddCommand(passwdCmd)
	SuperclusterCmd.AddCommand(userCmd)
	SuperclusterCmd.AddCommand(stopCmd)
}
