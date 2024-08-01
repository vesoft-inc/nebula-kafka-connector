package restore

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"path/filepath"
	"strconv"
	"strings"
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

	// only support restore one cluster now
	clusterId       int64
	backupClusterId int64
	catalogOwner    string

	force bool

	rootUri    string
	backupName string
	backSuffix string
}

func NewRestore(ctx context.Context, cfg *config.RestoreConfig) (*Restore, error) {
	r := &Restore{
		ctx:             ctx,
		cfg:             cfg,
		rootUri:         cfg.Backend.Uri(),
		backupName:      cfg.BackupName,
		clusterId:       cfg.ClusterId,
		backupClusterId: cfg.BackupClusterId,
		catalogOwner:    cfg.CatalogOwner,
		force:           cfg.Force,
	}

	var err error
	r.amg, err = clients.NewAgentManager(ctx, cfg.AgentsAddr)
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

	// get cluster
	clusters, err := r.meta.ListClusters(r.amg, r.cfg.ClusterId)
	if err != nil {
		return nil, fmt.Errorf("list cluster failed: %w", err)
	}

	restoreCluster := make([]*clients.ClusterServiceInfo, 0)
	for _, cluster := range clusters {
		if cluster.ClusterId == r.clusterId {
			restoreCluster = append(restoreCluster, cluster)
			break
		}
	}
	if len(restoreCluster) == 0 {
		return nil, fmt.Errorf("restore cluster %d not found", r.clusterId)
	}

	r.clusters, err = utils.NewNebulaClusters(restoreCluster, r.amg)
	if err != nil {
		return nil, fmt.Errorf("new nebula clusters failed: %w", err)
	}

	// get meta cluster
	metaResp, err := r.meta.ShowMeta()
	if err != nil {
		return nil, fmt.Errorf("show meta failed: %w", err)
	}
	metaCluster := make([]*clients.ServiceInfo, 0)
	for _, service := range metaResp.Services {
		agent, err := r.amg.GetAgent(service.Host)
		if err != nil {
			return nil, fmt.Errorf("get agent %s failed: %w", service.Host, err)
		}
		installPath, err := agent.GetInstallPath(service.Type)
		if err != nil {
			return nil, fmt.Errorf("get metad %s install path failed: %w", service.Host, err)
		}
		dataPaths, err := agent.GetDataPaths(service.Type, installPath)
		if err != nil {
			return nil, fmt.Errorf("get metad %s data path failed: %w", service.Host, err)
		}

		metaCluster = append(metaCluster, &clients.ServiceInfo{
			ServiceId:   service.Id,
			ServiceType: service.Type,
			Host:        service.Host,
			Port:        service.Port,
			InstallPath: installPath,
			DataPaths:   dataPaths,
		})
	}

	r.metaCluster = metaCluster

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

	// stop specify cluster's storageds
	err := r.stopStorageds()
	if err != nil {
		return fmt.Errorf("stop cluster storageds failed: %w", err)
	}

	// stop specify cluster's graphds
	err = r.stopGraphds()
	if err != nil {
		return fmt.Errorf("stop cluster graphd failed: %w", err)
	}

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

	// restore meta service by map
	hostPartMap, serviceIdMap, err := r.restoreMeta(r.backupMetas[0])
	if err != nil {
		return fmt.Errorf("restore cluster meta failed: %w", err)
	}
	logger.Info("Restore meta service successfully.")

	//download backup storage data from external storage to cluster
	if err = r.downloadStorage(hostPartMap, r.backupMetas[0]); err != nil {
		return fmt.Errorf("download storage data to cluster failed: %w", err)
	}
	logger.Info("Download storage data to cluster successfully.")

	// play back storage data
	if err = r.playBackStorageData(mapToString(serviceIdMap)); err != nil {
		return fmt.Errorf("playback storage data failed: %w", err)
	}
	logger.Info("Play back storage data successfully.")

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

	//after success restore, cleanup the download checkpoints if needed
	err = r.cleanupDownloadCheckpoints()
	if err != nil {
		return fmt.Errorf("clean up download checkpoints failed: %w", err)
	}
	logger.Info("Cleanup download checkpoints successfully.")

	// after success restore, cleanup the backup data if needed
	err = r.cleanupOriginalData()
	if err != nil {
		return fmt.Errorf("clean up origin data failed: %w", err)
	}
	logger.Info("Cleanup origin data successfully.")

	return nil
}

func (r *Restore) checkPhysicalTopology(bakClusters []*meta.ClusterBackupInfo) error {
	if len(bakClusters) != 1 {
		return fmt.Errorf("backup cluster count should be 1, but %d", len(bakClusters))
	}

	if len(r.clusters.GetStorages(r.clusterId)) != len(bakClusters[0].StorageInfos) {
		return fmt.Errorf("cluster %d and %d storage count are not consistent: %d vs %d",
			r.clusterId, bakClusters[0].ClusterId,
			len(r.clusters.GetStorages(r.clusterId)), len(bakClusters[0].StorageInfos))
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

		srcPath := filepath.Join(service.DataPaths[0], "data")
		dstPath := fmt.Sprintf("%s%s", srcPath, r.backSuffix)
		if err = agent.MoveDir(srcPath, dstPath); err != nil && !utils.IsNotExist(err) {
			return fmt.Errorf("move dir from %s to %s failed: %w", srcPath, dstPath, err)
		}

		log.WithField("origin path", srcPath).
			WithField("backup path", dstPath).
			WithField("origin not exist", utils.IsNotExist(err)).
			Info("Backup origin meta data path successfully.")
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
					srcPath := filepath.Join(d, "data")
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

func (r *Restore) cleanupDownloadCheckpoints() error {
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

		dstPath := filepath.Join(service.DataPaths[0], "checkpoint")
		if err = agent.RemoveDir(dstPath); err != nil {
			return fmt.Errorf("remove dir %s failed: %w", dstPath, err)
		}

		log.WithField("download metad checkpoint path", dstPath).
			Info("Cleanup download metad checkpoint data successfully.")
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
					dstPath := filepath.Join(d, "checkpoint")
					if err = agent.RemoveDir(dstPath); err != nil {
						return fmt.Errorf("remove dir %s failed: %w", dstPath, err)
					}

					log.WithField("download storaged checkpoint path", dstPath).
						Info("Cleanup download storaged checkpoint data successfully.")
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

		srcPath := filepath.Join(service.DataPaths[0], "data")
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
					srcPath := filepath.Join(d, "data")
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

		// meta kv data path: {nebulaData}/checkpoint/{backup_name}
		localDir := filepath.Join(service.DataPaths[0], "checkpoint", r.backupName)

		if err = agent.DownloadFile(backend, localDir, true); err != nil {
			return fmt.Errorf("download meta files from %s to %s failed: %w", externalUri, localDir, err)
		}
	}

	return nil
}

func (r *Restore) restoreMeta(backupRes *meta.CreateBackupResp) (map[string][]string, map[int64]int64, error) {
	storageIdMap := make(map[int64]int64)
	hostPartMap := make(map[string][]string)
	storageIdHostMap := make(map[int64]string)

	bakClusterMap := utils.FlattenClusterMap(backupRes)
	oldStorages := bakClusterMap[r.backupClusterId]
	curStorages := r.clusters.GetStorages(r.clusterId)

	for _, s := range curStorages {
		storageIdHostMap[s.ServiceId] = s.Host
	}

	for i, s := range curStorages {
		storageIdMap[oldStorages[i].ServiceId] = s.ServiceId
	}

	clusterRestoreInfos := []*meta.ClusterRestoreInfo{
		{
			NewClusterId: r.clusterId,
			MetaBackups:  backupRes.ClusterBackupInfos[0].MetaBackups,
			ServiceMap:   storageIdMap,
			CatalogOwner: r.catalogOwner,
		},
	}

	req := &meta.RestoreReq{
		ClusterMap:          map[int64]int64{r.backupClusterId: r.clusterId},
		ClusterRestoreInfos: clusterRestoreInfos,
		Force:               r.force,
	}

	log.Infof("restore req clustermap: %+v, ", req.ClusterMap)
	for _, info := range req.ClusterRestoreInfos {
		log.Infof("restore req cluster restore info: %+v, ", info)
	}

	resp, err := r.meta.Restore(req)
	if err != nil {
		return nil, nil, fmt.Errorf("restore meta failed: %w", err)
	}

	for partId, serviceIds := range resp.PartServiceMap {
		key := utils.GenPartKey(r.backupClusterId, partId)
		hostPartMap[key] = make([]string, 0)
		for _, serviceId := range serviceIds {
			hostPartMap[key] = append(hostPartMap[key], storageIdHostMap[serviceId])
		}
	}

	return hostPartMap, storageIdMap, nil
}

func (r *Restore) downloadStorage(hostPartMap map[string][]string, backupRes *meta.CreateBackupResp) error {
	logx.Infof("download storage data, hostpartmap: %v", hostPartMap)

	group := async.NewGroup(context.TODO(), r.cfg.Concurrency, "download storaged partition")
	curClusterId := r.clusterId

	dataPathSelector := utils.NewPathSelectorMap(r.clusters.GetStorages(curClusterId))
	dataPathMap := make(map[string]string)

	parts := utils.FlattenClusterBackupInfo(backupRes.ClusterBackupInfos[0])
	for _, part := range parts {
		storageUri, _ := utils.UriJoin(r.rootUri, r.backupName, "data")
		key := utils.GenPartKey(part.ClusterId, part.PartId)

		strClusterId := strconv.Itoa(int(r.backupClusterId))
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

	return group.Wait()
}

func (r *Restore) playBackStorageData(serviceIdMap string) error {
	group := async.NewGroup(context.TODO(), r.cfg.Concurrency, "playback storaged data")

	storages := r.clusters.GetStorages(r.clusterId)
	for _, s := range storages {
		agent, err := r.amg.GetAgent(s.Host)
		if err != nil {
			return fmt.Errorf("get agent for storage %s failed: %w", s.Host, err)
		}

		worker := func() error {
			if err = agent.DBPlayBack(r.backupName, s.InstallPath, strings.Join(s.DataPaths, ","), serviceIdMap); err != nil {
				return fmt.Errorf("playback storaged data failed: %w", err)
			}
			return nil
		}

		group.Add(func(stopCh chan interface{}) {
			stopCh <- worker()
		})
	}
	return group.Wait()
}

func mapToString(m map[int64]int64) string {
	pairs := make([]string, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, fmt.Sprintf("%d:%d", k, v))
	}
	return strings.Join(pairs, ",")
}
