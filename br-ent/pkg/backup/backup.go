package backup

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/async"

	agentstorage "github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/storage"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/clients"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/utils"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

type Backup struct {
	ctx  context.Context
	cfg  *config.BackupConfig
	meta *clients.NebulaMeta
	amg  *clients.AgentManager

	clusters []*clients.ServiceGroupServiceInfo
	sto      agentstorage.ExternalStorage
}

func NewBackup(ctx context.Context, cfg *config.BackupConfig) (*Backup, error) {
	b := &Backup{
		ctx: ctx,
		cfg: cfg,
	}

	var err error
	b.amg, err = clients.NewAgentManager(ctx, cfg.AgentsAddr)
	if err != nil {
		return nil, fmt.Errorf("create agent manager failed: %w", err)
	}

	b.meta, err = clients.NewMeta(cfg.MetaAddr, cfg.Username, cfg.Password, nil)
	if err != nil {
		return nil, fmt.Errorf("create meta client failed: %s", err)
	}

	b.sto, err = agentstorage.New(cfg.Backend)
	if err != nil {
		return nil, fmt.Errorf("create storage failed: %w", err)
	}

	clusters, err := b.meta.ListServiceGroups(b.amg, b.cfg.ServiceGroupId, b.cfg.Spec)
	if err != nil {
		return nil, fmt.Errorf("list cluster failed: %w", err)
	}

	b.clusters = clusters

	return b, nil
}

// upload the meta backup files in host to external uri
// localDir are absolute meta checkpoint folder in host filesystem
// targetUri is external storage's uri, which is meta's root dir,
// has pattern like local://xxx, s3://xxx
func (b *Backup) uploadMeta(backupRes *meta.CreateBackupResp, targetUri string) error {
	group := async.NewGroup(context.TODO(), b.cfg.Concurrency, "upload meta backup files")

	meatHost := strings.Split(b.meta.LeaderAddr(), ":")[0]
	agent, err := b.amg.GetAgent(meatHost)
	if err != nil {
		return err
	}

	for _, bakServiceGroup := range backupRes.ServiceGroupBackupInfos {
		if len(bakServiceGroup.MetaBackups) == 0 {
			return fmt.Errorf("meta backup files are empty")
		}

		// source: {metaDataPath}/checkpoint/{backupName}/{cluster_id}
		// target: {backupRoot}/{backupName}/meta/{clusterId}
		target, _ := utils.UriJoin(targetUri, strconv.Itoa(int(bakServiceGroup.ServiceGroupId)))
		source := path.Dir(bakServiceGroup.MetaBackups[0])
		backend, err := b.sto.GetDir(b.ctx, target)
		if err != nil {
			return fmt.Errorf("get storage backend for %s failed: %w", targetUri, err)
		}

		worker := func() error {
			if err = agent.UploadFile(backend, source, true); err != nil {
				return fmt.Errorf("upload file by agent failed: %w", err)
			}
			return nil
		}

		group.Add(func(stopCh chan interface{}) {
			stopCh <- worker()
		})
	}

	return group.Wait()
}

func (b *Backup) generateMetaFile(backupRes *meta.CreateBackupResp) (string, error) {
	tmpMetaPath := filepath.Join(utils.LocalTmpDir, fmt.Sprintf("%s.meta", backupRes.BackupName))
	return tmpMetaPath, utils.DumpMetaToFile(backupRes, tmpMetaPath)
}
