package main

import (
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/cmd"
)

var brCmd = &cobra.Command{
	Use:   "br",
	Short: "Process br command",
	Long:  `Execute br command in cli mode.`,
}

func init() {
	rootCmd.AddCommand(brCmd)
	brCmd.AddCommand(cmd.NewBackupCmd(), cmd.NewVersionCmd(), cmd.NewRestoreCmd(), cmd.NewCleanupCmd(), cmd.NewShowCmd())
}
