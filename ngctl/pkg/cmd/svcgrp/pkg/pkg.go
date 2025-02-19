package pkg

import "github.com/spf13/cobra"

var InstallCmd = &cobra.Command{
	Use:   "install-pkg",
	Short: "Install a NebulaGraph package on a host",
	Long:  "Install a NebulaGraph package on a host",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var UninstallCmd = &cobra.Command{
	Use:   "uninstall-pkg",
	Short: "Uninstall a NebulaGraph package on a host",
	Long:  "Uninstall a NebulaGraph package on a host",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
