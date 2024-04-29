package restore

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/storage"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/async"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/clients"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/br-ent/pkg/utils"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

func GetBackupSuffix() string {
	return fmt.Sprintf("_old_%d", time.Now().Unix())
}

type Restore struct {
	ctx  context.Context
	cfg  *config.RestoreConfig
	meta *clients.NebulaMeta
	amg  *clients.AgentManager

	sto storage.ExternalStorage

	metaCluster []*clients.ServiceInfo
	clusters    *utils.NebulaClusters

	// backupMetas store backup meta list for restore
	backupMetas []*meta.CreateBackupResp

	clusterIdMapping map[int64]int64

	rootUri    string
	backupName string
	backSuffix string
}

func NewRestore(ctx context.Context, cfg *config.RestoreConfig) (*Restore, error) {
	r := &Restore{
		ctx:              ctx,
		cfg:              cfg,
		rootUri:          cfg.Backend.Uri(),
		backupName:       cfg.BackupName,
		clusterIdMapping: cfg.ClusterIdMapping,
	}

	var err error
	r.amg, err = clients.NewAgentManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("create agent manager failed: %w", err)

	}
	r.meta, err = clients.NewMeta(cfg.MetaAddr, cfg.Username, cfg.Password, nil)
	if err != nil {
		return nil, fmt.Errorf("create meta client failed: %s", err)
	}

	r.sto, err = storage.New(cfg.Backend)
	if err != nil {
		return nil, fmt.Errorf("create storage failed: %w", err)
	}

	clusters, err := r.meta.ListClusters()
	if err != nil {
		return nil, fmt.Errorf("list cluster failed: %w", err)
	}
	r.clusters = utils.NewNebulaClusters(clusters, r.amg)

	return r, nil
}

// Restore restores all nebula cluster which storage is one engine for one part
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
func (r *Restore) Restore() error {
	logger := log.WithField("full restore", r.cfg.BackupName)

	// create localTmp dir to save tmp backup.meta file
	if err := utils.EnsureDir(utils.LocalTmpDir); err != nil {
		return err
	}
	defer func() {
		if err := utils.RemoveDir(utils.LocalTmpDir); err != nil {
			log.WithError(err).Errorf("Remove tmp dir %s failed.", utils.LocalTmpDir)
		}
	}()

	// load cluster's backup meta list
	if err := r.loadBakMetas(r.cfg.BackupName); err != nil {
		return err
	}
	logger.WithField("download backup meta", r.backupMetas).Info("download backup meta success")

	// check this cluster's topology with info kept in backup meta
	if err := r.checkPhysicalTopology(r.backupMetas[0].ClusterBackupInfos); err != nil {
		return fmt.Errorf("physical topology not consistent: %w", err)
	}

	// stop all cluster
	err := r.stopAllClustersWithLM()
	if err != nil {
		return fmt.Errorf("stop cluster failed: %w", err)
	}
	time.Sleep(time.Second * 10)
	logger.Info("Stop cluster successfully.")

	// backup original data
	err = r.backupOriginalData()
	if err != nil {
		return fmt.Errorf("backup origin data path failed: %w", err)
	}
	logger.Info("Backup origin cluster data successfully.")

	// download backup meta data from external storage to cluster
	err = r.downloadMeta()
	if err != nil {
		return fmt.Errorf("download meta data to cluster failed: %w", err)
	}
	logger.Info("Download meta data to cluster successfully.")

	// start meta service first
	err = r.startMetaCluster()
	if err != nil {
		return fmt.Errorf("start meta service failed: %w", err)
	}
	time.Sleep(time.Second * 10)
	logger.Info("Start meta service successfully.")

	// restore meta service by map
	hostPartMap, err := r.restoreMeta(r.backupMetas[0])
	if err != nil {
		return fmt.Errorf("restore cluster meta failed: %w", err)
	}
	logger.Info("Restore meta service successfully.")

	//download backup storage data from external storage to cluster
	if err = r.downloadStorage(hostPartMap, r.backupMetas[0]); err != nil {
		return fmt.Errorf("download storage data to cluster failed: %w", err)
	}
	logger.Info("Download storage data to cluster successfully.")

	// start storage and graph service
	err = r.startAllClusterStorage()
	if err != nil {
		return fmt.Errorf("start storage service failed: %w", err)
	}
	err = r.startAllClusterGraph()
	if err != nil {
		return fmt.Errorf("start graph service failed: %w", err)
	}
	logger.Info("Start storage and graph services successfully.")

	// after success restore, cleanup the backup data if needed
	err = r.cleanupOriginalData()
	if err != nil {
		return fmt.Errorf("clean up origin data failed: %w", err)
	}
	logger.Info("Cleanup origin data successfully.")

	return nil
}

func (r *Restore) checkPhysicalTopology(bakClusters []*meta.ClusterBackupInfo) error {
	if len(r.clusters.GetClusterIds()) != len(bakClusters) {
		return fmt.Errorf("cluster count not consistent: %d vs %d", len(r.clusters.GetClusterIds()), len(bakClusters))
	}
	for _, bakCluster := range bakClusters {
		currentClusterId := r.clusterIdMapping[bakCluster.ClusterId]
		if len(r.clusters.GetStorages(currentClusterId)) != len(bakCluster.StorageInfos) {
			return fmt.Errorf("cluster %d and %d storage count are not consistent: %d vs %d",
				currentClusterId, bakCluster.ClusterId,
				len(r.clusters.GetStorages(currentClusterId)), len(bakCluster.StorageInfos))
		}
	}
	return nil
}

func (r *Restore) backupOriginalData() error {
	r.backSuffix = GetBackupSuffix()

	// backup meta data
	for _, service := range r.metaCluster {
		agent, err := r.amg.GetAgent(service.Host)
		if err != nil {
			return fmt.Errorf("get agent %s failed: %w", service.Host, err)
		}

		if len(service.DataPaths) != 1 {
			return fmt.Errorf("meta service: %s should only have one data dir, but %d",
				service.Host, len(service.DataPaths))
		}

		srcPath := filepath.Join(service.DataPaths[0], "nebula")
		dstPath := fmt.Sprintf("%s%s", srcPath, r.backSuffix)
		if err = agent.MoveDir(srcPath, dstPath); err != nil && !utils.IsNotExist(err) {
			return fmt.Errorf("move dir from %s to %s failed: %w", srcPath, dstPath, err)
		}

		log.WithField("origin path", srcPath).
			WithField("backup path", dstPath).
			WithField("origin not exist", utils.IsNotExist(err)).
			Info("Backup origin storage data path successfully.")
	}

	// backup storage data
	for _, cluster := range r.clusters.GetClusters() {
		for _, service := range cluster {
			if service.ServiceType == meta.ServiceTypeStoraged {
				agent, err := r.amg.GetAgent(service.Host)
				if err != nil {
					return fmt.Errorf("get agent %s failed: %w", service.Host, err)
				}

				for _, d := range service.DataPaths {
					srcPath := filepath.Join(d, "nebula")
					dstPath := fmt.Sprintf("%s%s", srcPath, r.backSuffix)
					if err = agent.MoveDir(srcPath, dstPath); err != nil && !utils.IsNotExist(err) {
						return fmt.Errorf("move dir from %s to %s failed: %w", srcPath, dstPath, err)
					}

					log.WithField("origin path", srcPath).
						WithField("backup path", dstPath).
						WithField("origin not exist", utils.IsNotExist(err)).
						Info("Backup origin storage data path successfully.")
				}
			}
		}
	}

	return nil
}

func (r *Restore) cleanupOriginalData() error {
	// cleanup backup meta data
	for _, service := range r.metaCluster {
		agent, err := r.amg.GetAgent(service.Host)
		if err != nil {
			return fmt.Errorf("get agent %s failed: %w", service.Host, err)
		}

		if len(service.DataPaths) != 1 {
			return fmt.Errorf("meta service: %s should only have one data dir, but %d",
				service.Host, len(service.DataPaths))
		}

		srcPath := filepath.Join(service.DataPaths[0], "nebula")
		dstPath := fmt.Sprintf("%s%s", srcPath, r.backSuffix)
		if err = agent.RemoveDir(dstPath); err != nil {
			return fmt.Errorf("remove dir %s failed: %w", dstPath, err)
		}

		log.WithField("backup path", dstPath).
			Info("Cleanup backup origin storage data path successfully.")
	}

	// cleanup backup storage data
	for _, cluster := range r.clusters.GetClusters() {
		for _, service := range cluster {
			if service.ServiceType == meta.ServiceTypeStoraged {
				agent, err := r.amg.GetAgent(service.Host)
				if err != nil {
					return fmt.Errorf("get agent %s failed: %w", service.Host, err)
				}

				for _, d := range service.DataPaths {
					srcPath := filepath.Join(d, "nebula")
					dstPath := fmt.Sprintf("%s%s", srcPath, r.backSuffix)
					if err = agent.RemoveDir(dstPath); err != nil {
						return fmt.Errorf("remove dir %s failed: %w", dstPath, err)
					}

					log.WithField("backup path", dstPath).
						Info("Cleanup backup origin storage data path successfully.")
				}
			}
		}
	}

	return nil
}

// loadBakMetas load the backup chain through base reference until base == "", only support
func (r *Restore) loadBakMetas(backupName string) error {
	// check backup dir existence
	rootUri, err := utils.UriJoin(r.cfg.Backend.Uri(), backupName)
	if err != nil {
		return err
	}
	exist := r.sto.ExistDir(r.ctx, rootUri)
	if !exist {
		return fmt.Errorf("backup dir %s does not exist", rootUri)
	}

	// download and parse backup meta file
	backupMetaName := fmt.Sprintf("%s.meta", backupName)
	metaUri, _ := utils.UriJoin(rootUri, backupMetaName)
	tmpLocalPath := filepath.Join(utils.LocalTmpDir, backupMetaName)
	err = r.sto.Download(r.ctx, tmpLocalPath, metaUri, false)
	if err != nil {
		return fmt.Errorf("download %s to %s failed: %w", metaUri, tmpLocalPath, err)
	}
	bakMeta, err := utils.ParseMetaFromFile(tmpLocalPath)
	if err != nil {
		return fmt.Errorf("parse backup meta file %s failed: %w", tmpLocalPath, err)
	}

	r.backupMetas = append(r.backupMetas, bakMeta)

	return nil
}

func (r *Restore) downloadMeta() error {
	// {backupRoot}/{backupName}/meta/*.sst
	externalUri, _ := utils.UriJoin(r.rootUri, r.backupName, "meta")
	backend, err := r.sto.GetDir(r.ctx, externalUri)
	if err != nil {
		return fmt.Errorf("get storage backend for %s failed: %w", externalUri, err)
	}

	// download meta backup files to every meta service
	for _, service := range r.metaCluster {
		agent, err := r.amg.GetAgent(service.Host)
		if err != nil {
			return fmt.Errorf("get agent for meta %s failed: %w", service.Host, err)
		}

		// meta kv data path: {nebulaData}/meta/{backup_name}/{cluster_id}/
		localDir := filepath.Join(service.DataPaths[0], "checkpoint")
		if err = agent.DownloadFile(backend, localDir, true); err != nil {
			return fmt.Errorf("download meta files from %s to %s failed: %w", externalUri, localDir, err)
		}
	}

	return nil
}

func (r *Restore) restoreMeta(backupRes *meta.CreateBackupResp) (map[string][]string, error) {
	storageIdMap := make(map[int64]int64)
	hostPartMap := make(map[string][]string)

	bakClusterMap := utils.FlattenClusterMap(backupRes)

	for cur, old := range r.clusterIdMapping {
		curStorages := r.clusters.GetStorages(cur)
		oldStorages := bakClusterMap[old]
		for i, s := range curStorages {
			storageIdMap[s.ServiceId] = oldStorages[i].ServiceId
			for _, info := range oldStorages[i].CkptInfos {
				key := utils.GenPartKey(cur, info.PartId)
				hostPartMap[key] = append(hostPartMap[key], s.Host)
			}
		}
	}

	req := &meta.RestoreReq{
		//MetaBackups: backupRes.MetaBackups,
		ClusterMap: r.clusterIdMapping,
		ServiceMap: storageIdMap,
	}

	if _, err := r.meta.Restore(req); err != nil {
		return nil, fmt.Errorf("restore meta failed: %w", err)
	}

	return hostPartMap, nil
}

func (r *Restore) downloadStorage(hostPartMap map[string][]string, backupRes *meta.CreateBackupResp) error {
	group := async.NewGroup(context.TODO(), r.cfg.Concurrency, "download storaged partition")
	for _, bakCluster := range backupRes.ClusterBackupInfos {
		curClusterId := r.clusterIdMapping[bakCluster.ClusterId]

		dataPathSelector := utils.NewPathSelectorMap(r.clusters.GetStorages(curClusterId))
		dataPathMap := make(map[string]string)

		parts := utils.FlattenClusterBackupInfo(bakCluster)
		for _, part := range parts {
			storageUri, _ := utils.UriJoin(r.rootUri, r.backupName, "data")
			key := utils.GenPartKey(part.ClusterId, part.PartId)

			strClusterId := strconv.Itoa(int(curClusterId))
			strPartId := strconv.Itoa(int(part.PartId))
			for _, host := range hostPartMap[key] {
				agent, err := r.amg.GetAgent(host)
				if err != nil {
					return fmt.Errorf("get agent for storage %s failed: %w", host, err)
				}

				externalUri, _ := utils.UriJoin(storageUri, strClusterId, strPartId)
				// avoid agent.DownloadFile prefix bugs
				externalUri += "/"

				source, err := r.sto.GetDir(r.ctx, externalUri)
				if err != nil {
					return fmt.Errorf("get storage backend for %s failed: %w", externalUri, err)
				}

				// ensure every part's all checkpoint place in same dataPath
				dataPathKey := utils.GenDataPathKey(host, key)
				dataPath, hasPlace := dataPathMap[dataPathKey]
				if !hasPlace {
					dataPath = dataPathSelector[host].EvenlyGet()
					dataPathMap[dataPathKey] = dataPath
				}

				target := filepath.Join(dataPath, "checkpoint", backupRes.BackupName, strPartId)

				// source: {backupRoot}/{backupName}/data/{clusterId}/{partId}/
				// target: {nebulaDataPath}/checkpoint/{backupName}/{partId}
				worker := func() error {
					if err = agent.DownloadFile(source, target, true); err != nil {
						return fmt.Errorf("download %s to %s failed:%w", externalUri, target, err)
					}
					return nil
				}

				group.Add(func(stopCh chan interface{}) {
					stopCh <- worker()
				})
			}
		}
	}

	return group.Wait()
}
