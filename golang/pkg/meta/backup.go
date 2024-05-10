package meta

import (
	"context"
	"fmt"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/admin"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/common"
)

type (
	ListBackupClustersReq struct{}

	ListBackupClustersResp struct {
		*HeaderResponse
		MetaCluster []*ServiceInfo
		Clusters    []*ClusterServiceInfo
	}

	CkptInfo struct {
		PartId   int64
		CkptPath string
	}

	StorageCheckpointInfo struct {
		ServiceId int64
		Host      string
		Port      uint32
		CkptInfos []CkptInfo
	}

	ClusterBackupInfo struct {
		ClusterId     int64
		PartitionNum  int32
		ReplicaFactor int32
		MetaBackups   []string
		StorageInfos  []*StorageCheckpointInfo
	}

	CreateBackupReq struct {
		BackupName string
		ClusterIds []int64
	}

	CreateBackupResp struct {
		*HeaderResponse
		BackupName         string
		CreateTime         int64
		ClusterBackupInfos []*ClusterBackupInfo
	}

	DropBackupReq struct {
		BackupNames []string
	}

	DropBackupResp struct {
		*HeaderResponse
	}

	RestoreReq struct {
		MetaBackups []string
		ClusterMap  map[int64]int64
		ServiceMap  map[int64]int64
	}

	RestoreResp struct {
		*HeaderResponse
	}

	ShowMetaResp struct {
		*HeaderResponse
		Services []*ServiceInfo
	}
)

func (c *metaClient) CreateBackup(req *CreateBackupReq) (*CreateBackupResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.connectTimeout)
	defer cancel()
	in := &admin.CreateBackupRequest{
		Header:     &admin.RequestHeader{Token: c.token},
		BackupName: []byte(req.BackupName),
		ClusterIds: req.ClusterIds,
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.CreateBackup(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	if err = responseIsErr(resp); err != nil {
		return nil, err
	}
	response, ok := resp.(*admin.CreateBackupResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	header, err := getResponseHeader(response)
	if err != nil {
		return nil, fmt.Errorf("get response header failed: %w", err)
	}

	clusterBackupInfos := make([]*ClusterBackupInfo, 0, len(response.ClusterInfos))
	for _, cbInfo := range response.ClusterInfos {
		storageInfos := make([]*StorageCheckpointInfo, 0, len(cbInfo.StorageInfos))
		for _, storage := range cbInfo.StorageInfos {
			ckptInfos := make([]CkptInfo, 0, len(storage.CkptParts))
			for _, ckptPart := range storage.CkptParts {
				ckptInfos = append(ckptInfos, CkptInfo{
					PartId:   int64(ckptPart.PartId),
					CkptPath: ckptPart.CkptPath,
				})
			}

			storageInfos = append(storageInfos, &StorageCheckpointInfo{
				ServiceId: storage.ServiceId,
				Host:      string(storage.Host.GetHost()),
				Port:      storage.Host.GetPort(),
				CkptInfos: ckptInfos,
			})
		}

		clusterBackupInfos = append(clusterBackupInfos, &ClusterBackupInfo{
			ClusterId:     cbInfo.ClusterId,
			PartitionNum:  cbInfo.PartitionNum,
			ReplicaFactor: cbInfo.ReplicaFactor,
			MetaBackups:   bytesToStrings(cbInfo.MetaBackups),
			StorageInfos:  storageInfos,
		})
	}

	return &CreateBackupResp{
		HeaderResponse:     header,
		BackupName:         string(response.BackupName),
		CreateTime:         response.CreateTime,
		ClusterBackupInfos: clusterBackupInfos,
	}, nil
}

func (c *metaClient) DropBackup(req *DropBackupReq) (*DropBackupResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.connectTimeout)
	defer cancel()
	in := &admin.DropBackupRequest{
		Header:      &admin.RequestHeader{Token: c.token},
		BackupNames: stringsToBytes(req.BackupNames),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.DropBackup(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	if err = responseIsErr(resp); err != nil {
		return nil, err
	}
	response, ok := resp.(*admin.DropBackupResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	header, err := getResponseHeader(response)
	if err != nil {
		return nil, fmt.Errorf("get response header failed: %w", err)
	}
	return &DropBackupResp{
		HeaderResponse: header,
	}, err
}

func (c *metaClient) Restore(req *RestoreReq) (*RestoreResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.connectTimeout)
	defer cancel()
	in := &admin.RestoreRequest{
		Header: &admin.RequestHeader{Token: c.token},
		//MetaBackups: stringsToBytes(req.MetaBackups),
		ClusterMap: req.ClusterMap,
		//ServiceMap:  req.ServiceMap,
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.Restore(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	if err = responseIsErr(resp); err != nil {
		return nil, err
	}
	response, ok := resp.(*admin.RestoreResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	header, err := getResponseHeader(response)
	if err != nil {
		return nil, fmt.Errorf("get response header failed: %w", err)
	}

	return &RestoreResp{
		HeaderResponse: header,
	}, nil
}

func (c *metaClient) ShowMeta() (*ShowMetaResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.connectTimeout)
	defer cancel()
	in := &admin.ShowMetaRequest{
		Header: &admin.RequestHeader{Token: c.token},
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.ShowMetaInfo(ctx, in)
	})
	if err != nil {
		return nil, err
	}
	if err = responseIsErr(resp); err != nil {
		return nil, err
	}
	response, ok := resp.(*admin.ShowMetaResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}
	header, err := getResponseHeader(response)
	if err != nil {
		return nil, fmt.Errorf("get response header failed: %w", err)
	}

	services := make([]*ServiceInfo, 0, len(response.Info.Peers))
	for _, host := range response.Info.Peers {
		services = append(services, &ServiceInfo{
			Host: string(host.GetHost()),
			Port: host.GetPort(),
			Type: ServiceType(common.ServiceType_META),
		})
	}

	return &ShowMetaResp{
		HeaderResponse: header,
		Services:       services,
	}, nil
}

func bytesToStrings(b [][]byte) []string {
	s := make([]string, 0, len(b))
	for _, v := range b {
		s = append(s, string(v))
	}
	return s
}

func stringsToBytes(s []string) [][]byte {
	b := make([][]byte, 0, len(s))
	for _, v := range s {
		b = append(b, []byte(v))
	}
	return b
}

func NewCreateBackupReq(backupName string, clusterIds []int64) *CreateBackupReq {
	return &CreateBackupReq{
		BackupName: backupName,
		ClusterIds: clusterIds,
	}
}

func NewDropBackupReq(backupNames []string) *DropBackupReq {
	return &DropBackupReq{
		BackupNames: backupNames,
	}
}
