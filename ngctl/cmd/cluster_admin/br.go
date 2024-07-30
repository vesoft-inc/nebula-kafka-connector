package cluster_admin

import (
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/cmd"
)

var brCmd = &cobra.Command{
	Use:   "br",
	Short: "Backup the data of a cluster.",
	Long:  `Users can bakcup the data of a cluster and restore a cluster from backup files. Users can also cleanup or show the backup files.`,
}

func init() {
	brCmd.AddCommand(cmd.NewBackupCmd(), cmd.NewVersionCmd(), cmd.NewRestoreCmd(), cmd.NewCleanupCmd(), cmd.NewShowCmd())
}
