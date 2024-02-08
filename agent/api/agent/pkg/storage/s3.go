package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vesoft-inc/nebula-ng-tools/agent/api/agent/pkg/limiter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/zeromicro/go-zero/core/logx"
)

type S3Config struct {
	Endpoint  string
	Region    string
	Bucket    string
	Path      string
	AccessKey string
	SecretKey string
}

func (s *S3Config) DeepCopy() *S3Config {
	cp := &S3Config{
		Endpoint:  s.Endpoint,
		Region:    s.Region,
		Bucket:    s.Bucket,
		Path:      s.Path,
		AccessKey: s.AccessKey,
		SecretKey: s.SecretKey,
	}
	return cp
}

type s3Client struct {
	backend *Backend
	sess    *session.Session
	client  *s3.S3
}

func NewS3(b *Backend) (ExternalStorage, error) {
	if b.Type() != S3Type {
		return nil, fmt.Errorf("bad format s3 uri: %s", b.Uri())
	}

	creds := credentials.NewStaticCredentials(b.S3.AccessKey, b.S3.SecretKey, "")
	forcePath := CheckEndpoint(b.S3.Endpoint)
	region := "default"
	if b.S3.Region != "" {
		region = b.S3.Region
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(b.S3.Endpoint),
		S3ForcePathStyle: aws.Bool(forcePath), // ip:port
		Credentials:      creds,
	}))

	logx.Infof("Create s3 backend with endpoint: %s, region: %s", b.S3.Endpoint, region)

	return &s3Client{
		backend: b,
		sess:    sess,
		client:  s3.New(sess),
	}, nil
}

func (s *s3Client) downloadToFile(file, key string) error {
	// Take rate limiter count by object size
	if limiter.Rate.IsSet() {
		size, err := s.GetObjectSize(s.backend.S3.Bucket, key)
		if err != nil {
			return err
		}
		limiter.Rate.Wait(size)
	}

	// Create the directories in the path
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0775); err != nil {
		return fmt.Errorf("ensure dir %s failed: %w", dir, err)
	}

	// Set up the local file
	fd, err := os.Create(file)
	if err != nil {
		if err != nil {
			return fmt.Errorf("create file %s failed: %w", file, err)
		}
	}
	defer fd.Close()

	// Download the file using the AWS SDK for Go
	req := &s3.GetObjectInput{
		Bucket: aws.String(s.backend.S3.Bucket),
		Key:    aws.String(key),
	}

	downloader := s3manager.NewDownloader(s.sess)
	numBytes, err := downloader.Download(fd, req)
	if err != nil {
		return fmt.Errorf("download from %s to %s failed: %w", key, file, err)
	}

	logx.Infof("Download from %s to %s successfully, bytes=%d.", key, file, numBytes)
	return nil
}

func (s *s3Client) downloadPrefix(localDir, prefix string) error {
	req := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.backend.S3.Bucket),
		Prefix: aws.String(prefix),
	}

	err := s.client.ListObjectsV2Pages(req, func(p *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range p.Contents {
			relPath, err := filepath.Rel(prefix, *obj.Key)
			if err != nil {
				logx.Errorf("Get relative path failed, key=%s, prefix=%s", *obj.Key, prefix)
				return false
			}

			localFile := filepath.Join(localDir, relPath)
			err = s.downloadToFile(localFile, *obj.Key)
			if err != nil {
				logx.Errorf("Download file from %s to %s failed: %v", *obj.Key, localFile, err)
				return false
			}
		}
		return true
	})
	if err != nil {
		return fmt.Errorf("download %s recursively failed: %w", s.backend.Uri(), err)
	}

	logx.Infof("Download from %s to %s successfully.", s.backend.Uri(), localDir)
	return nil
}

func (s *s3Client) Download(ctx context.Context, localPath, externalUri string, recursively bool) error {
	b := s.backend.DeepCopy()
	err := b.SetUri(externalUri)
	if err != nil {
		return fmt.Errorf("download, check and set s3 uri %s failed: %w", externalUri, err)
	}

	if recursively {
		return s.downloadPrefix(localPath, b.S3.Path)
	} else {
		return s.downloadToFile(localPath, b.S3.Path)
	}
}

func (s *s3Client) uploadToStorage(key, file string) error {
	// Take rate limiter count by file size
	if limiter.Rate.IsSet() {
		srcInfo, err := os.Stat(file)
		if err != nil {
			return err
		}
		limiter.Rate.Wait(srcInfo.Size())
	}

	fd, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("open file %s failed: %w when upload", file, err)
	}
	defer fd.Close()

	uploader := s3manager.NewUploader(s.sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.backend.S3.Bucket),
		Key:    aws.String(key),
		Body:   fd,
	})
	if err != nil {
		return fmt.Errorf("upload from %s to %s failed: %w", file, key, err)
	}

	logx.Infof("Upload from %s to %s successfully.", file, key)
	return nil
}

func (s *s3Client) uploadPrefix(prefix, localDir string) error {
	walker := make(fileWalk)
	go func() {
		// Gather the files to upload by walking the path recursively
		if err := filepath.Walk(localDir, walker.Walk); err != nil {
			logx.Errorf("Walk failed: %v", err)
		}
		close(walker)
	}()

	// For each file found walking, upload it to s3Client
	for path := range walker {
		rel, err := filepath.Rel(localDir, path)
		if err != nil {
			logx.Errorf("Unable to get relative path: %s.", path)
		}
		key := filepath.Join(prefix, rel)

		err = s.uploadToStorage(key, path)
		if err != nil {
			return fmt.Errorf("upload from %s to %s failed: %w", path, key, err)
		}
	}

	logx.Infof("Upload from %s to %s recursively.", localDir, s.backend.Uri())
	return nil
}

func (s *s3Client) Upload(ctx context.Context, externalUri, localPath string, recursively bool) error {
	b := s.backend.DeepCopy()
	err := b.SetUri(externalUri)
	if err != nil {
		return fmt.Errorf("upload, check and set s3 uri %s failed: %w", externalUri, err)
	}

	if recursively {
		return s.uploadPrefix(b.S3.Path, localPath)
	} else {
		return s.uploadToStorage(b.S3.Path, localPath)
	}
}

func (s *s3Client) GetObjectSize(bucket, key string) (int64, error) {
	req := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	resp, err := s.client.HeadObject(req)
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}
