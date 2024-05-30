package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/log"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/restore"

	log2 "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRestoreCmd() *cobra.Command {
	restoreCmd := &cobra.Command{
		Use:          "restore",
		Short:        "Restore data files, notice that it will restart the cluster",
		SilenceUsage: true,
	}
	config.AddCommonFlags(restoreCmd.PersistentFlags())
	config.AddRestoreFlags(restoreCmd.PersistentFlags())
	restoreCmd.AddCommand(newFullRestoreCmd())
	return restoreCmd
}

func newFullRestoreCmd() *cobra.Command {
	fullRestoreCmd := &cobra.Command{
		Use:   "full",
		Short: "full restore all clusters data files",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := log.SetLog(cmd.Flags())
			if err != nil {
				return fmt.Errorf("init logger failed: %w", err)
			}

			cfg := &config.RestoreConfig{}
			err = cfg.ParseFlags(cmd.Flags())
			if err != nil {
				return err
			}

			r, err := restore.NewRestore(context.TODO(), cfg)
			if err != nil {
				return err
			}

			startTime := time.Now()
			log2.Info("Start to restore clusters...")
			err = r.Restore()
			if err != nil {
				return err
				log2.Errorf("Restore failed: %v, will try to recovery original data...\n", err)
				f, ferr := restore.NewFixFrom(r)
				if ferr != nil {
					return err
				}

				ferr = f.Fix()
				if ferr != nil {
					log2.Errorf("Fix failed when restore failed %v", ferr)
				}

				return err
			}
			log2.Infof("Restore cluster succeed, time spent: %v", time.Since(startTime))
			return nil
		},
	}

	return fullRestoreCmd
}
