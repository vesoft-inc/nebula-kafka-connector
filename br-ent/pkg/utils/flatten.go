package utils

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

// PartInfo save a unique part's information for FlattenBackupMeta
type PartInfo struct {
	ClusterId      int64
	ServiceId      int64
	PartId         int64
	Host           string
	Port           uint32
	CheckpointPath string
}

type FlattenedParts map[string]*PartInfo

func FlattenClusterMap(resp *meta.CreateBackupResp) map[int64][]*meta.StorageCheckpointInfo {
	clusterMap := make(map[int64][]*meta.StorageCheckpointInfo)
	for _, cluster := range resp.ClusterBackupInfos {
		clusterMap[cluster.ClusterId] = cluster.StorageInfos
	}
	return clusterMap
}

// FlattenClusterBackupInfo flatten backup resp to a map for convenience
// because of (clusterId + partId) can specify a unique part
func FlattenClusterBackupInfo(backupInfo *meta.ClusterBackupInfo) FlattenedParts {
	backupMap := make(map[string]*PartInfo)
	for _, storage := range backupInfo.StorageInfos {
		for _, part := range storage.CkptInfos {
			key := GenPartKey(backupInfo.ClusterId, part.PartId)
			backupMap[key] = &PartInfo{
				ClusterId:      backupInfo.ClusterId,
				ServiceId:      storage.ServiceId,
				Host:           storage.Host,
				Port:           storage.Port,
				PartId:         part.PartId,
				CheckpointPath: part.CkptPath,
			}
		}
	}

	return backupMap
}

// GenPartKey generate a unique key for part
func GenPartKey(clusterId, partId int64) string {
	return fmt.Sprintf("%d-%d", clusterId, partId)
}

// GenDataPathKey generate a unique key for data_path
func GenDataPathKey(host string, partKey string) string {
	return fmt.Sprintf("%s-%s", host, partKey)
}
