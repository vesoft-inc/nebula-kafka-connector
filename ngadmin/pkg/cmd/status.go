package cmd

import "github.com/spf13/cobra"

func RegisteStatusCmd(rootCmd *cobra.Command) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "status",
		Short: "view nebulagraph cluster status",
	}
	rootCmd.AddCommand(cmd)
	return cmd
}
