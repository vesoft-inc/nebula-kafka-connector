package storage

import (
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/internal/types"
	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/storage"
)

func (s *storageService) S3Download(req *types.S3DownloadReq) (*types.S3DownloadResp, error) {
	backend := &storage.Backend{
		S3: &storage.S3Config{
			Endpoint:  req.Endpoint,
			Region:    req.Region,
			AccessKey: req.AccessKey,
			SecretKey: req.SecretKey,
			Bucket:    req.Bucket,
			Path:      req.Path,
		},
	}

	s3client, err := storage.NewS3(backend)
	if err != nil {
		return nil, err
	}

	if err = s3client.Download(s.ctx, req.LocalPath, backend.Uri(), true); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *storageService) S3Upload(req *types.S3UploadReq) (*types.S3UploadResp, error) {
	backend := &storage.Backend{
		S3: &storage.S3Config{
			Endpoint:  req.Endpoint,
			Region:    req.Region,
			AccessKey: req.AccessKey,
			SecretKey: req.SecretKey,
			Bucket:    req.Bucket,
			Path:      req.Path,
		},
	}

	s3client, err := storage.NewS3(backend)
	if err != nil {
		return nil, err
	}

	if err = s3client.Upload(s.ctx, backend.Uri(), req.LocalPath, true); err != nil {
		return nil, err
	}

	return nil, nil
}
