package source

import (
	"fmt"
	"github.com/colinmarc/hdfs/v2"
	"io"
)

var _ Source = (*hdfsSource)(nil)

type (
	hdfsSource struct {
		c    *Config
		obj  io.ReadCloser
		size int64
	}
)

func openHDFSFile(c *Config) (*hdfsSource, error) {
	// open connection to HDFS server
	client, err := hdfs.New(fmt.Sprintf("%s:%d", c.HDFS.Host, c.HDFS.Port))
	if err != nil {
		return nil, err
	}

	// open the file
	obj, err := client.Open(c.HDFS.Path)
	if err != nil {
		return nil, err
	}

	return &hdfsSource{
		c:    c,
		obj:  obj,
		size: obj.Stat().Size(),
	}, nil
}

func (o *hdfsSource) Config() *Config {
	return o.c
}

func (o *hdfsSource) Size() (int64, error) {
	return o.size, nil
}

func (o *hdfsSource) Read(p []byte) (int, error) {
	return o.obj.Read(p)
}

func (o *hdfsSource) Close() error {
	return o.obj.Close()
}
