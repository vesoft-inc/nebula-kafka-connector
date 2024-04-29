package main

import (
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/cmd"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "br-ent",
		Short: "br-ent is a backup and restore tool",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	rootCmd.AddCommand(cmd.NewBackupCmd(), cmd.NewVersionCmd(), cmd.NewRestoreCmd(), cmd.NewCleanupCmd(), cmd.NewShowCmd())
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
