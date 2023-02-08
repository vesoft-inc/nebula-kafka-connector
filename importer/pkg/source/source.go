//go:generate mockgen -source=source.go -destination source_mock.go -package source Source,Sizer
package source

import (
	"io"
)

type (
	Source interface {
		Config() *Config
		Sizer
		io.Reader
		io.Closer
	}

	Sizer interface {
		Size() (int64, error)
	}
)

func Open(c *Config) (Source, error) {
	// TODO: support blob and so on
	switch {
	case c.S3 != nil:
		return openS3File(c)
	case c.OSS != nil:
		return openOSSFile(c)
	case c.FTP != nil:
		return openFTPFile(c)
	case c.SFTP != nil:
		return openSFTPFile(c)
	case c.HDFS != nil:
		return openHDFSFile(c)
	default:
		return openLocalFile(c)
	}
}
