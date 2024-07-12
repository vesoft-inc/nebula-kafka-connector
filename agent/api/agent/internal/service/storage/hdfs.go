package storage

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/audit"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/storage"
)

func (s *storageService) HDFSDownload(req *types.HDFSDownloadReq) (*types.HDFSDownloadResp, error) {
	if err := audit.RecordOperation(s.ctx, audit.OpHDFSDownload, fmt.Sprintf("download file from %s to %s", req.Path, req.LocalPath)); err != nil {
		return nil, err
	}

	backend := &storage.Backend{
		HDFS: &storage.HDFSConfig{
			Address:  req.Address,
			Username: req.Username,
			Path:     req.Path,
			Kerberos: storage.KerberosConfig{
				Enable:                       req.Kerberos.Enable,
				Principle:                    req.Kerberos.Principle,
				KeytabFilePath:               req.Kerberos.KeytabFilePath,
				ConfigFilePath:               req.Kerberos.ConfigFilePath,
				KerberosServicePrincipleName: req.Kerberos.KerberosServicePrincipleName,
			},
		},
	}

	hdfsclient, err := storage.NewHDFS(backend)
	if err != nil {
		return nil, err
	}

	if err = hdfsclient.Download(s.ctx, req.LocalPath, backend.Uri(), true); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *storageService) HDFSUpload(req *types.HDFSUploadReq) (*types.HDFSUploadResp, error) {
	if err := audit.RecordOperation(s.ctx, audit.OpHDFSUpload, fmt.Sprintf("upload file to %s from %s", req.Path, req.LocalPath)); err != nil {
		return nil, err
	}

	backend := &storage.Backend{
		HDFS: &storage.HDFSConfig{
			Address:  req.Address,
			Username: req.Username,
			Path:     req.Path,
			Kerberos: storage.KerberosConfig{
				Enable:                       req.Kerberos.Enable,
				Principle:                    req.Kerberos.Principle,
				KeytabFilePath:               req.Kerberos.KeytabFilePath,
				ConfigFilePath:               req.Kerberos.ConfigFilePath,
				KerberosServicePrincipleName: req.Kerberos.KerberosServicePrincipleName,
			},
		},
	}

	hdfsclient, err := storage.NewHDFS(backend)
	if err != nil {
		return nil, err
	}

	if err = hdfsclient.Upload(s.ctx, backend.Uri(), req.LocalPath, true); err != nil {
		return nil, err
	}

	return nil, nil
}
