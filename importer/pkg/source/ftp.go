package source

import (
	"fmt"
	"io"
	"time"

	"github.com/jlaffaye/ftp"
)

var _ Source = (*ftpSource)(nil)

type (
	ftpSource struct {
		c    *Config
		obj  io.ReadCloser
		size int64
	}
)

func openFTPFile(c *Config) (*ftpSource, error) {
	// open connection to FTP server
	client, err := ftp.Dial(fmt.Sprintf("%s:%d", c.FTP.Host, c.FTP.Port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	// login to ftp server
	err = client.Login(c.FTP.Username, c.FTP.Password)
	if err != nil {
		return nil, fmt.Errorf("login to ftp server failed: %w", err)
	}

	// get the file size
	size, err := client.FileSize(c.FTP.Path)
	if err != nil {
		return nil, fmt.Errorf("getting file size failed: %w", err)
	}

	// open the file
	obj, err := client.Retr(c.FTP.Path)
	if err != nil {
		return nil, fmt.Errorf("opening file failed: %w", err)
	}

	return &ftpSource{
		c:    c,
		obj:  obj,
		size: size,
	}, nil
}

func (f *ftpSource) Config() *Config {
	return f.c
}

func (f *ftpSource) Size() (int64, error) {
	return f.size, nil
}

func (f *ftpSource) Read(p []byte) (int, error) {
	return f.obj.Read(p)
}

func (f *ftpSource) Close() error {
	return f.obj.Close()
}
