package cleanup

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	agentstorage "github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/storage"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/clients"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/utils"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

type Cleanup struct {
	ctx    context.Context
	cfg    *config.CleanupConfig
	client *clients.NebulaMeta
	amg    *clients.AgentManager
	sto    agentstorage.ExternalStorage
}

func NewCleanup(ctx context.Context, cfg *config.CleanupConfig) (*Cleanup, error) {
	sto, err := agentstorage.New(cfg.Backend)
	if err != nil {
		return nil, err
	}

	client, err := clients.NewMeta(cfg.MetaAddr, cfg.Username, cfg.Password, nil)

	amg, err := clients.NewAgentManager(ctx, cfg.AgentsAddr)
	if err != nil {
		return nil, err
	}

	return &Cleanup{
		ctx:    ctx,
		cfg:    cfg,
		client: client,
		amg:    amg,
		sto:    sto,
	}, nil
}

func (c *Cleanup) cleanNebula() error {
	if _, err := c.client.DropBackup(meta.NewDropBackupReq([]string{c.cfg.BackupName})); err != nil {
		return fmt.Errorf("drop backup failed: %w", err)
	}

	log.Debugf("Drop backup %s successfully.", c.cfg.BackupName)

	return nil
}

func (c *Cleanup) cleanExternal() error {
	backupUri, err := utils.UriJoin(c.cfg.Backend.Uri(), c.cfg.BackupName)
	if err != nil {
		return err
	}

	err = c.sto.RemoveDir(c.ctx, backupUri)
	if err != nil {
		return fmt.Errorf("remove %s in external storage failed: %w", backupUri, err)
	}
	log.Debugf("Remove %s successfully.", backupUri)

	// Local backend's data lay in different cluster machines,
	// which should be handled separately
	if c.cfg.Backend.Local != nil {
		for addr, agent := range c.amg.GetAgents() {
			backupPath := strings.TrimPrefix(backupUri, agentstorage.LocalPrefix)
			if err = agent.RemoveDir(backupPath); err != nil {
				return fmt.Errorf("remove %s in host: %s failed: %w", backupPath, addr, err)
			}

			log.Debugf("Remove local data %s in %s successfully.", backupPath, addr)
		}
	}

	return nil
}

func (c *Cleanup) Clean() error {
	logger := log.WithField("backup name", c.cfg.BackupName)

	logger.Info("Start to cleanup data in cluster self.")
	err := c.cleanNebula()
	if err != nil {
		log.Errorf("clean nebula local data failed: %v", err)
	}

	logger.Info("Start cleanup data in external storage.")
	err = c.cleanExternal()
	if err != nil {
		return fmt.Errorf("clean external storage data failed: %w", err)
	}

	logger.Info("Clean up backup data successfully.")
	return nil
}
