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
	// TODO: support hdfs, s3, blob and so on
	return openLocalFile(c)
}
