package cmd

import (
	"context"
	"fmt"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/cleanup"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/backup"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/log"

	log2 "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewBackupCmd() *cobra.Command {
	backupCmd := &cobra.Command{
		Use:          "backup",
		Short:        "Backup data files to external storage for restore",
		SilenceUsage: true,
	}

	config.AddCommonFlags(backupCmd.PersistentFlags())
	config.AddBackupFlags(backupCmd.PersistentFlags())
	backupCmd.AddCommand(newFullBackupCmd())
	return backupCmd
}

func newFullBackupCmd() *cobra.Command {
	fullBackupCmd := &cobra.Command{
		Use:   "full",
		Short: "Full backup data files",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := log.SetLog(cmd.Flags())
			if err != nil {
				return fmt.Errorf("init logger failed: %w", err)
			}

			cfg := &config.BackupConfig{}
			err = cfg.ParseFullFlags(cmd.Flags())
			if err != nil {
				return fmt.Errorf("parse flags failed: %w", err)
			}

			b, err := backup.NewBackup(context.TODO(), cfg)
			if err != nil {
				return err
			}

			startTime := time.Now()
			log2.Info("Start to full backup cluster...")
			backupName, err := b.FullBackup()
			if err != nil {
				log2.Errorf("Full backup failed: %v, will try to clean the remaining garbage...\n", err)
				if backupName != "" {
					if cleanErr := clean(&config.CleanupConfig{
						BackupName: backupName,
						Backend:    cfg.Backend,
						TLSConfig:  cfg.TLSConfig,
					}); cleanErr != nil {
						return cleanErr
					}
					log2.Errorf("Cleanup full backup %s successfully after backup failed.", backupName)
				}
				return err
			}

			log2.Infof("Full backup succeed, time spent: %v", time.Since(startTime))
			return nil
		},
	}

	return fullBackupCmd
}

func clean(cfg *config.CleanupConfig) error {
	c, err := cleanup.NewCleanup(context.TODO(), cfg)
	if err != nil {
		return fmt.Errorf("create cleanup for %s failed: %w", cfg.BackupName, err)
	}

	if err = c.Clean(); err != nil {
		return fmt.Errorf("cleanup %s failed when backup failed: %w", cfg.BackupName, err)
	}
	return nil
}
