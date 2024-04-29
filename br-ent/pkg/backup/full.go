package backup

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/async"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/utils"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"

	log "github.com/sirupsen/logrus"
)

// FullBackup backs up full data in given external storage, and return the backup name
/*
backup_root/backup_name
  - meta
    - xxx.sst
    - ...
  - data
    - clusterId
      - partId
      - ...
    - ...
  - backup_name.meta
*/
func (b *Backup) FullBackup() (string, error) {
	// call the meta service, create backup files in each local
	backupRes, err := b.meta.CreateFullBackup(b.cfg.BackupName, b.cfg.ClusterId)
	if err != nil {
		if backupRes != nil {
			return backupRes.BackupName, nil
		}
		return "", err
	}

	if len(backupRes.ClusterBackupInfos) == 0 {
		return "", fmt.Errorf("no backup info returned")
	}

	backupName := backupRes.BackupName
	logger := log.WithField("name", backupName)
	logger.WithField("backup meta", utils.StringifyBackup(backupRes)).Debugf("Create backup checkpoints in machine's local.")

	// ensure root dir
	rootUri, err := utils.UriJoin(b.cfg.Backend.Uri(), backupName)
	if err != nil {
		return backupName, err
	}
	err = b.sto.EnsureDir(b.ctx, rootUri, false)
	if err != nil {
		return backupName, fmt.Errorf("ensure dir %s failed: %w", rootUri, err)
	}
	logger.WithField("root", rootUri).Info("Ensure backup root dir.")

	// upload meta files
	metaUri, err := utils.UriJoin(rootUri, "meta")
	if err != nil {
		return backupName, err
	}
	if err = b.uploadMeta(backupRes, metaUri); err != nil {
		return backupName, err
	}
	logger.WithField("meta", metaUri).Info("Upload meta successfully.")

	// upload storage files
	storageUri, err := utils.UriJoin(rootUri, "data")
	if err != nil {
		return backupName, err
	}
	err = b.uploadFullStorage(backupRes, storageUri)
	if err != nil {
		return backupName, fmt.Errorf("upload storage failed %w", err)
	}
	logger.WithField("data", storageUri).Info("Upload data backup successfully.")

	// generate backup meta files and upload
	if err = utils.EnsureDir(utils.LocalTmpDir); err != nil {
		return backupName, err
	}
	defer func() {
		if err = utils.RemoveDir(utils.LocalTmpDir); err != nil {
			log.WithError(err).Errorf("Remove tmp dir %s failed.", utils.LocalTmpDir)
		}
	}()

	tmpMetaPath, err := b.generateMetaFile(backupRes)
	if err != nil {
		return backupName, fmt.Errorf("write meta to tmp path failed: %w", err)
	}
	logger.WithField("tmp path", tmpMetaPath).Info("Write meta data to local tmp file successfully.")
	backupMetaPath, err := utils.UriJoin(rootUri, filepath.Base(tmpMetaPath))
	if err != nil {
		return backupName, err
	}
	err = b.sto.Upload(b.ctx, backupMetaPath, tmpMetaPath, false)
	if err != nil {
		return backupName, fmt.Errorf("upload local tmp file to remote storage %s failed: %w", backupMetaPath, err)
	}
	logger.WithField("remote path", backupMetaPath).Info("Upload tmp backup meta file to remote.")

	// drop backup files in cluster machine local and local tmp files
	_, err = b.meta.DropBackup(meta.NewDropBackupReq([]string{backupName}))
	if err != nil {
		return backupName, fmt.Errorf("drop backup %s in cluster local failed: %w", backupName, err)
	}
	logger.Info("Drop backup in cluster and local tmp folder successfully.")

	return backupName, nil
}

func (b *Backup) uploadFullStorage(backupRes *meta.CreateBackupResp, targetUri string) error {
	group := async.NewGroup(context.TODO(), b.cfg.Concurrency, "full upload storaged partition")

	for _, bakCluster := range backupRes.ClusterBackupInfos {
		parts := utils.FlattenClusterBackupInfo(bakCluster)
		for _, part := range parts {
			agent, err := b.amg.GetAgent(part.Host)
			if err != nil {
				return err
			}

			// source: {nebulaDataPath}/checkpoints/{backupName}/{partId}
			// target: {backupRoot}/{backupName}/data/{clusterId}/{partId}
			target, _ := utils.UriJoin(targetUri, strconv.Itoa(int(part.ClusterId)), strconv.Itoa(int(part.PartId)))
			source := part.CheckpointPath
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
	}

	return group.Wait()
}
