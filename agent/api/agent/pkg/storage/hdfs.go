package storage

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/limiter"

	"github.com/colinmarc/hdfs/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

type HDFSConfig struct {
	Address  string
	Username string
	Path     string
	Kerberos KerberosConfig
}

func (h *HDFSConfig) DeepCopy() *HDFSConfig {
	cp := &HDFSConfig{
		Address:  h.Address,
		Username: h.Username,
		Path:     h.Path,
		Kerberos: h.Kerberos,
	}
	return cp
}

type hdfsClient struct {
	backend *Backend
	client  *hdfs.Client
}

func NewHDFS(b *Backend) (ExternalStorage, error) {
	if b.Type() != HDFSType {
		return nil, fmt.Errorf("bad format hdfs uri: %s", b.Uri())
	}

	clientOptions := hdfs.ClientOptions{
		Addresses: []string{b.HDFS.Address},
		User:      b.HDFS.Username,
	}

	if b.HDFS.Kerberos.Enable {
		krbClient, err := getKerberosClient(b.HDFS.Kerberos)
		if err != nil {
			logx.Errorf("failed to create hdfs kerberos client: %v", err)
			return nil, err
		}
		clientOptions.KerberosClient = krbClient
		clientOptions.KerberosServicePrincipleName = b.HDFS.Kerberos.KerberosServicePrincipleName
	}

	client, err := hdfs.NewClient(clientOptions)
	if err != nil {
		logx.Errorf("failed to create hdfs client: %v", err)
		return nil, err
	}

	return &hdfsClient{
		backend: b,
		client:  client,
	}, nil
}

func (h *hdfsClient) downloadPrefix(localDir, prefix string) error {
	files, err := h.client.ReadDir(prefix)
	if err != nil {
		return fmt.Errorf("failed to read hdfs directory: %w", err)
	}

	for _, file := range files {
		hdfsPath := filepath.Join(prefix, file.Name())
		localPath := filepath.Join(localDir, file.Name())

		if file.IsDir() {
			err = os.MkdirAll(localPath, 0o775)
			if err != nil {
				return fmt.Errorf("failed to create local directory: %w", err)
			}

			err = h.downloadPrefix(localPath, hdfsPath)
			if err != nil {
				return err
			}
		} else {
			err = h.downloadToFile(localPath, hdfsPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *hdfsClient) downloadToFile(file, key string) error {
	hdfsFile, err := h.client.Open(key)
	if err != nil {
		return fmt.Errorf("failed to open hdfs file: %w", err)
	}
	defer hdfsFile.Close()

	localFile, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, hdfsFile)
	if err != nil {
		return fmt.Errorf("failed to download file from hdfs: %w", err)
	}

	return nil
}

func (h *hdfsClient) Download(ctx context.Context, localPath, externalUri string, recursively bool) error {
	b := h.backend.DeepCopy()
	err := b.SetUri(externalUri)
	if err != nil {
		return fmt.Errorf("download, check and set hdfs uri %s failed: %w", externalUri, err)
	}

	srcInfo, err := h.client.Stat(b.HDFS.Path)
	if os.IsNotExist(err) {
		return fmt.Errorf("hdfs path not exist: %s", b.HDFS.Path)
	}
	if srcInfo.IsDir() && !recursively {
		return fmt.Errorf("%s is directory, must download recursively", externalUri)
	}

	if srcInfo.IsDir() {
		return h.downloadPrefix(localPath, b.HDFS.Path)
	} else {
		return h.downloadToFile(localPath, b.HDFS.Path)
	}
}

func (h *hdfsClient) uploadPrefix(prefix, localDir string) error {
	walker := make(fileWalk)
	go func() {
		if err := filepath.Walk(localDir, walker.Walk); err != nil {
			logx.Errorf("Walk failed: %v", err)
		}
		close(walker)
	}()

	for path := range walker {
		rel, err := filepath.Rel(localDir, path)
		if err != nil {
			logx.Errorf("Unable to get relative path: %s.", path)
		}
		key := filepath.Join(prefix, rel)

		err = h.uploadToStorage(key, path)
		if err != nil {
			return fmt.Errorf("upload from %s to %s failed: %w", path, key, err)
		}
	}

	logx.Infof("Upload from %s to %s recursively.", localDir, h.backend.Uri())
	return nil
}

func (h *hdfsClient) uploadToStorage(key, file string) error {
	if limiter.Rate.IsSet() {
		srcInfo, err := os.Stat(file)
		if err != nil {
			return err
		}
		limiter.Rate.Wait(srcInfo.Size())
	}

	localFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer localFile.Close()

	fileWriter, err := h.client.Create(key)
	if err != nil {
		return fmt.Errorf("failed to create hdfs file: %w", err)
	}
	defer fileWriter.Close()

	_, err = io.Copy(fileWriter, localFile)
	if err != nil {
		return fmt.Errorf("failed to upload file to hdfs: %w", err)
	}

	return nil
}

func (h *hdfsClient) Upload(ctx context.Context, externalUri, localPath string, recursively bool) error {
	b := h.backend.DeepCopy()
	err := b.SetUri(externalUri)
	if err != nil {
		return fmt.Errorf("upload, check and set hdfs uri %s failed: %w", externalUri, err)
	}

	srcInfo, err := h.client.Stat(b.HDFS.Path)
	if os.IsNotExist(err) {
		return fmt.Errorf("hdfs path not exist: %s", b.HDFS.Path)
	}
	if srcInfo.IsDir() && !recursively {
		return fmt.Errorf("%s is directory, must upload recursively", externalUri)
	}

	if srcInfo.IsDir() {
		return h.uploadPrefix(b.HDFS.Path, localPath)
	} else {
		return h.uploadToStorage(b.HDFS.Path, localPath)
	}
}

func (h *hdfsClient) ExistDir(ctx context.Context, uri string) bool {
	b := h.backend.DeepCopy()
	err := b.SetUri(uri)
	if err != nil {
		log.WithError(err).WithField("uri", uri).Error("check and set uri failed when test existDir")
		return false
	}

	_, err = h.client.Stat(b.HDFS.Path)
	if err != nil {
		log.WithError(err).WithField("uri", uri).Error("check hdfs path not exist")
		return false
	}

	return true
}

func (h *hdfsClient) EnsureDir(ctx context.Context, uri string, recursively bool) error {
	b := h.backend.DeepCopy()
	err := b.SetUri(uri)
	if err != nil {
		return fmt.Errorf("ensure dir, check and set hdfs uri %s failed: %w", uri, err)
	}
	return nil
}

func (h *hdfsClient) GetDir(ctx context.Context, uri string) (*Backend, error) {
	b := h.backend.DeepCopy()
	err := b.SetUri(uri)
	if err != nil {
		return nil, fmt.Errorf("get dir, check and set hdfs uri %s failed: %w", uri, err)
	}

	return b, nil
}

func (h *hdfsClient) ListDir(ctx context.Context, uri string) ([]string, error) {
	b := h.backend.DeepCopy()
	err := b.SetUri(uri)
	if err != nil {
		return nil, fmt.Errorf("list dir, check and set hdfs uri %s failed: %w", uri, err)
	}

	files, err := h.client.ReadDir(b.HDFS.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read hdfs directory: %w", err)
	}

	dirs := make([]string, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}

	return dirs, nil
}

func (h *hdfsClient) RemoveDir(ctx context.Context, uri string) error {
	b := h.backend.DeepCopy()
	err := b.SetUri(uri)
	if err != nil {
		return fmt.Errorf("remove dir, check and set hdfs uri %s failed: %w", uri, err)
	}

	if err = h.client.Remove(b.HDFS.Path); err != nil {
		return fmt.Errorf("failed to remove hdfs path: %w", err)
	}

	return nil
}
