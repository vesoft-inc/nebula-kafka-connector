package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout meta server.",
	Long:  `logout`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := metaClient.Logout(); err != nil {
			return err
		}
		if err := ngctl.ClearMetaToken(); err != nil {
			return metaConsoleError("clear meta session failed", err.Error())
		}

		fmt.Fprintln(metaOutput, "Logout succeeded.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
