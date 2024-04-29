package storage

import (
	"fmt"
	"net/url"
	"strings"
)

type Backend struct {
	S3    *S3Config
	HDFS  *HDFSConfig
	Local *LocalConfig
}

const (
	S3Prefix    = "s3://"
	HDFSPrefix  = "hdfs://"
	LocalPrefix = "local://"
)

type BackendType int

const (
	S3Type BackendType = iota
	HDFSType
	LocalType
)

func (t BackendType) String() string {
	switch t {
	case S3Type:
		return "s3"
	case HDFSType:
		return "hdfs"
	case LocalType:
		return "local"
	default:
		return "unknown"
	}
}

func ParseType(uri string) BackendType {
	if strings.HasPrefix(uri, S3Prefix) {
		return S3Type
	}
	if strings.HasPrefix(uri, HDFSPrefix) {
		return HDFSType
	}
	if strings.HasPrefix(uri, LocalPrefix) {
		return LocalType
	}

	return -1
}

func (b *Backend) Type() BackendType {
	t := ParseType(b.Uri())
	return t
}

func (b *Backend) Uri() string {
	if b.S3 != nil {
		return S3Prefix + b.S3.Bucket + "/" + b.S3.Path
	}
	if b.HDFS != nil {
		return HDFSPrefix + b.HDFS.Path
	}
	if b.Local != nil {
		return LocalPrefix + b.Local.Path
	}

	return "nil path"
}

func (b *Backend) SetUri(uri string) error {
	t := ParseType(uri)
	switch t {
	case S3Type:
		u, err := parseUri(uri)
		if err != nil {
			return err
		}

		if b.S3 == nil {
			b.S3 = &S3Config{}
		}

		b.S3.Bucket = u.Host
		b.S3.Path = u.Path
	case HDFSType:
		if b.HDFS == nil {
			b.HDFS = &HDFSConfig{}
		}
		if strings.HasPrefix(uri, HDFSPrefix) {
			b.HDFS.Path = strings.TrimPrefix(uri, HDFSPrefix)
		}
	case LocalType:
		if b.Local == nil {
			b.Local = &LocalConfig{}
		}
		b.Local.Path = strings.TrimPrefix(uri, LocalPrefix)
	default:
		return fmt.Errorf("unknow storage backend type")
	}
	return nil
}

func (b *Backend) DeepCopy() *Backend {
	cp := &Backend{}
	switch b.Type() {
	case S3Type:
		if b.S3 != nil {
			cp.S3 = b.S3.DeepCopy()
		}
	case HDFSType:
		if b.HDFS != nil {
			cp.HDFS = b.HDFS.DeepCopy()
		}
	case LocalType:
		if b.Local != nil {
			cp.Local = b.Local.DeepCopy()

		}
	}
	return cp
}

func parseUri(uri string) (*url.URL, error) {
	t := ParseType(uri)
	u, err := url.Parse(uri)
	if err != nil || strings.ToLower(u.Scheme) != t.String() {
		return nil, fmt.Errorf("invalid uri: %s", uri)
	}
	u.Host = strings.TrimRight(u.Host, "/ ")
	u.Path = strings.TrimLeft(u.Path, "/ ")
	return u, nil
}
