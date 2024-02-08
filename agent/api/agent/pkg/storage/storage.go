package storage

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

// Downloader download data from externalUri in ExternalStorage to the localPath
// If localPath not exist, create it;
// If it exists, overwrite it when it is file, copy to it when it is folder
type Downloader interface {
	Download(ctx context.Context, localPath, externalUri string, recursively bool) error
}

// Uploader upload data from localPath to the externalUri in ExternalStorage
// When recursively is false, upload file from localPath to externalUri
// Else upload files in localPath folder to externalUri folder
type Uploader interface {
	Upload(ctx context.Context, externalUri, localPath string, recursively bool) error
}

// ExternalStorage will keep the authentication information and other configuration for the storage
// Then the functions only need uri as the parameter
type ExternalStorage interface {
	Downloader
	Uploader
}

func New(b *Backend) (ExternalStorage, error) {
	logx.Infof("Create type: %s storage, uri :%s", b.Type(), b.Uri())

	switch b.Type() {
	case S3Type:
		return NewS3(b)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", b.Type())
	}
}
