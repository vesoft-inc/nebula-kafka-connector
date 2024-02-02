package cmd

import "github.com/spf13/cobra"

func RegisteUninstallCmd(rootCmd *cobra.Command) {
	verifyCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "uninstall nebula cluster",
		Long:  `uninstall nebula cluster by config file which is default to .ngadmin.yaml, check if agent is running`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			return nil
		},
	}
	rootCmd.AddCommand(verifyCmd)
}
