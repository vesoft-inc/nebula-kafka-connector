package source

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var _ Source = (*ossSource)(nil)

type (
	ossSource struct {
		c    *Config
		obj  io.ReadCloser
		size int64
	}
)

func openOSSFile(c *Config) (*ossSource, error) {
	client, err := oss.New(c.OSS.Endpoint, c.OSS.AccessKey, c.OSS.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("create oss client faild: %w", err)
	}

	bucket, err := client.Bucket(c.OSS.Bucket)
	if err != nil {
		return nil, fmt.Errorf("get oss bucket faild: %w", err)
	}

	// get object size
	key := strings.TrimLeft(c.OSS.Key, "/")
	meta, err := bucket.GetObjectMeta(key)
	if err != nil {
		return nil, err
	}
	contentLength := meta.Get("Content-Length")
	size, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return nil, err
	}

	// get object
	obj, err := bucket.GetObject(key)
	if err != nil {
		return nil, err
	}

	return &ossSource{
		c:    c,
		obj:  obj,
		size: size,
	}, nil
}

func (o *ossSource) Config() *Config {
	return o.c
}

func (o *ossSource) Size() (int64, error) {
	return o.size, nil
}

func (o *ossSource) Read(p []byte) (int, error) {
	return o.obj.Read(p)
}

func (o *ossSource) Close() error {
	return o.obj.Close()
}
