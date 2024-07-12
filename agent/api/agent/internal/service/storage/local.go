package storage

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/audit"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/storage"
)

func (s *storageService) LocalUpload(req *types.LocalUploadReq) (*types.LocalUploadResp, error) {
	if err := audit.RecordOperation(s.ctx, audit.OpLocalUpload, fmt.Sprintf("upload file to %s from %s", req.Path, req.LocalPath)); err != nil {
		return nil, err
	}

	backend := &storage.Backend{
		Local: &storage.LocalConfig{
			Path: req.Path,
		},
	}

	localclient, err := storage.NewLocal(backend)
	if err != nil {
		return nil, err
	}

	if err = localclient.Upload(s.ctx, backend.Uri(), req.LocalPath, req.Recursively); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *storageService) LocalDownload(req *types.LocalDownloadReq) (*types.LocalDownloadResp, error) {
	if err := audit.RecordOperation(s.ctx, audit.OpLocalDownload, fmt.Sprintf("download file from %s to %s", req.Path, req.LocalPath)); err != nil {
		return nil, err
	}

	backend := &storage.Backend{
		Local: &storage.LocalConfig{
			Path: req.Path,
		},
	}

	localclient, err := storage.NewLocal(backend)
	if err != nil {
		return nil, err
	}

	if err = localclient.Download(s.ctx, req.LocalPath, backend.Uri(), req.Recursively); err != nil {
		return nil, err
	}

	return nil, nil
}
