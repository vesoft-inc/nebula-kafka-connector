package source

import "os"

var _ Source = (*localSource)(nil)

type (
	localSource struct {
		c *Config
		f *os.File
	}
)

func openLocalFile(c *Config) (*localSource, error) {
	f, err := os.Open(c.Path)
	if err != nil {
		return nil, err
	}
	return &localSource{
		c: c,
		f: f,
	}, nil
}

func (r *localSource) Size() (int64, error) {
	fi, err := r.f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func (r *localSource) Read(p []byte) (int, error) {
	return r.f.Read(p)
}

func (r *localSource) Close() error {
	return r.f.Close()
}
