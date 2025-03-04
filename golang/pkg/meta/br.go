package meta

import (
	"context"
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/golang/internal/generated_code/v5.0.0/proto/admin"
)

type (
	CreateBackupReq struct {
		BackupName       string
		ServiceGroupName string
		IncludeMeta      bool
	}

	CreateBackupResp struct {
		BackupName             string                  `json:"BackupName"`
		CreateTime             int64                   `json:"CreateTime"`
		MetaInfo               *CheckpointInfo         `json:"MetaInfo"`
		ServiceGroupBackupInfo *ServiceGroupBackupInfo `json:"ServiceGroupBackupInfo"`
	}

	ServiceGroupBackupInfo struct {
		ServiceGroupId int64             `json:"ServiceGroupId"`
		PartitionNum   int32             `json:"PartitionNum"`
		Checkpoints    []*CheckpointInfo `json:"Checkpoints"`
	}

	CkptPartInfo struct {
		PartId   uint32 `json:"PartId"`
		DataPath string `json:"DataPath"`
		WalPath  string `json:"WalPath"`
	}

	CheckpointInfo struct {
		ServiceId int64           `json:"ServiceId"`
		Host      string          `json:"Host"`
		Port      uint32          `json:"Port"`
		CkptParts []*CkptPartInfo `json:"CkptParts"`
	}

	DropBackupReq struct {
		BackupNames []string
	}

	ServiceGroupRestoreInfo struct {
		NewServiceGroupId int64
		MetaBackups       []string
		ServiceMap        map[int64]int64
		CatalogOwner      string
	}

	RestoreReq struct {
		ServiceGroup string
	}
)

func (c *metaClient) CreateBackup(req *CreateBackupReq) (*CreateBackupResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	in := &admin.CreateBackupRequest{
		Header:       &admin.RequestHeader{Token: c.token},
		BackupName:   []byte(req.BackupName),
		ServiceGroup: []byte(req.ServiceGroupName),
		IncludeMeta:  req.IncludeMeta,
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
	return convertCreateResponse(response), nil
}

func convertCreateResponse(resp *admin.CreateBackupResponse) *CreateBackupResp {
	if resp == nil {
		return nil
	}
	cr := &CreateBackupResp{
		BackupName: string(resp.BackupName),
		CreateTime: resp.CreateTime,
	}

	cr.MetaInfo = convertCheckpointInfo(resp.MetaInfo)
	cr.ServiceGroupBackupInfo = convertServiceGroupBackupInfo(resp.ServiceGroupInfo)
	return cr
}

func convertCheckpointInfo(resp *admin.CheckpointInfo) *CheckpointInfo {
	if resp == nil {
		return nil
	}
	ci := &CheckpointInfo{
		ServiceId: resp.ServiceId,
		Host:      string(resp.Host.GetHost()),
		Port:      resp.Host.GetPort(),
		CkptParts: make([]*CkptPartInfo, 0, len(resp.CkptParts)),
	}
	for _, part := range resp.CkptParts {
		ci.CkptParts = append(ci.CkptParts, &CkptPartInfo{
			PartId:   part.PartId,
			DataPath: string(part.DataPath),
			WalPath:  string(part.WalPath),
		})
	}
	return ci
}

func convertServiceGroupBackupInfo(resp *admin.ServiceGroupBackupInfo) *ServiceGroupBackupInfo {
	if resp == nil {
		return nil
	}
	sgi := &ServiceGroupBackupInfo{
		ServiceGroupId: resp.ServiceGroupId,
		PartitionNum:   resp.PartitionNum,
		Checkpoints:    make([]*CheckpointInfo, 0, len(resp.StorageInfos)),
	}
	for _, info := range resp.StorageInfos {
		sgi.Checkpoints = append(sgi.Checkpoints, convertCheckpointInfo(info))
	}
	return sgi
}

func (c *metaClient) DropBackups(req *DropBackupReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	in := &admin.DropBackupRequest{
		Header:      &admin.RequestHeader{Token: c.token},
		BackupNames: stringsToBytes(req.BackupNames),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.DropBackup(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)

}

func (c *metaClient) Restore(req *RestoreReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	in := &admin.RestoreRequest{
		Header:       &admin.RequestHeader{Token: c.token},
		ServiceGroup: []byte(req.ServiceGroup),
	}
	resp, err := c.execute(func() (responseHeader, error) {
		return c.client.Restore(ctx, in)
	})
	if err != nil {
		return err
	}
	return responseIsErr(resp)

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

func NewCreateBackupReq(backupName string, serviceGroupName string, includeMeta bool) *CreateBackupReq {
	return &CreateBackupReq{
		BackupName:       backupName,
		ServiceGroupName: serviceGroupName,
		IncludeMeta:      includeMeta,
	}
}

func NewDropBackupReq(backupNames []string) *DropBackupReq {
	return &DropBackupReq{
		BackupNames: backupNames,
	}
}

func NewRestoreReq(serviceGroup string) *RestoreReq {
	return &RestoreReq{
		ServiceGroup: serviceGroup,
	}
}
