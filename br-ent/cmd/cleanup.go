package cmd

import (
	"context"
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/cleanup"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/log"

	"github.com/spf13/cobra"
)

func NewCleanupCmd() *cobra.Command {
	cleanupCmd := &cobra.Command{
		Use:          "cleanup",
		Short:        "Cleanup backup files in external storage and cluster self",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := log.SetLog(cmd.Flags())
			if err != nil {
				return fmt.Errorf("init logger failed: %w", err)
			}

			cfg := &config.CleanupConfig{}
			err = cfg.ParseFlags(cmd.Flags())
			if err != nil {
				return fmt.Errorf("parse flags failed")
			}

			c, err := cleanup.NewCleanup(context.TODO(), cfg)
			if err != nil {
				return err
			}

			err = c.Clean()
			if err != nil {
				return err
			}

			return nil
		},
	}

	config.AddCommonFlags(cleanupCmd.PersistentFlags())
	config.AddCleanupFlags(cleanupCmd.PersistentFlags())
	return cleanupCmd
}
