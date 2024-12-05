package meta

import (
	"context"
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/golang/internal/generated_code/v5.0.0/proto/admin"
)

type (
	ListBackupServiceGroupsReq struct{}

	ListBackupServiceGroupsResp struct {
		*HeaderResponse
		MetaServiceGroup []*ServiceInfo
		ServiceGroups    []*ServiceGroupServiceInfo
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

	ServiceGroupBackupInfo struct {
		ServiceGroupId int64
		PartitionNum   int32
		ReplicaFactor  int32
		MetaBackups    []string
		StorageInfos   []*StorageCheckpointInfo
	}

	CreateBackupReq struct {
		BackupName      string
		ServiceGroupIds []int64
	}

	CreateBackupResp struct {
		*HeaderResponse
		BackupName              string
		CreateTime              int64
		ServiceGroupBackupInfos []*ServiceGroupBackupInfo
	}

	DropBackupReq struct {
		BackupNames []string
	}

	DropBackupResp struct {
		*HeaderResponse
	}

	ServiceGroupRestoreInfo struct {
		NewServiceGroupId int64
		MetaBackups       []string
		ServiceMap        map[int64]int64
		CatalogOwner      string
	}

	RestoreReq struct {
		ServiceGroupMap          map[int64]int64
		ServiceGroupRestoreInfos []*ServiceGroupRestoreInfo
		Force                    bool
	}

	RestoreResp struct {
		*HeaderResponse
		PartServiceMap map[int64][]int64
	}
)

func (c *metaClient) CreateBackup(req *CreateBackupReq) (*CreateBackupResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.connectTimeout)
	defer cancel()
	in := &admin.CreateBackupRequest{
		Header:          &admin.RequestHeader{Token: c.token},
		BackupName:      []byte(req.BackupName),
		ServiceGroupIds: req.ServiceGroupIds,
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

	ServiceGroupBackupInfos := make([]*ServiceGroupBackupInfo, 0, len(response.ServiceGroupInfos))
	for _, cbInfo := range response.ServiceGroupInfos {
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

		ServiceGroupBackupInfos = append(ServiceGroupBackupInfos, &ServiceGroupBackupInfo{
			ServiceGroupId: cbInfo.ServiceGroupId,
			PartitionNum:   cbInfo.PartitionNum,
			ReplicaFactor:  cbInfo.ReplicaFactor,
			MetaBackups:    bytesToStrings(cbInfo.MetaBackups),
			StorageInfos:   storageInfos,
		})
	}

	return &CreateBackupResp{
		HeaderResponse:          header,
		BackupName:              string(response.BackupName),
		CreateTime:              response.CreateTime,
		ServiceGroupBackupInfos: ServiceGroupBackupInfos,
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

	ServiceGroupInfos := make([]*admin.ServiceGroupRestoreInfo, 0)
	for _, info := range req.ServiceGroupRestoreInfos {
		ServiceGroupInfos = append(ServiceGroupInfos, &admin.ServiceGroupRestoreInfo{
			NewServiceGroupId: info.NewServiceGroupId,
			MetaBackups:       stringsToBytes(info.MetaBackups),
			ServiceMap:        info.ServiceMap,
			CatalogOwner:      []byte(info.CatalogOwner),
		})

	}

	in := &admin.RestoreRequest{
		Header:            &admin.RequestHeader{Token: c.token},
		ServiceGroupIdMap: req.ServiceGroupMap,
		ServiceGroupInfos: ServiceGroupInfos,
		Force:             req.Force,
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

	partServiceMap := make(map[int64][]int64)
	for _, partService := range response.Parts {
		partServiceMap[int64(partService.PartId)] = partService.ServiceIds
	}

	return &RestoreResp{
		HeaderResponse: header,
		PartServiceMap: partServiceMap,
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

func NewCreateBackupReq(backupName string, ServiceGroupIds []int64) *CreateBackupReq {
	return &CreateBackupReq{
		BackupName:      backupName,
		ServiceGroupIds: ServiceGroupIds,
	}
}

func NewDropBackupReq(backupNames []string) *DropBackupReq {
	return &DropBackupReq{
		BackupNames: backupNames,
	}
}
