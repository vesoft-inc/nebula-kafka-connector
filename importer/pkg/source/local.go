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
func (s *localSource) Config() *Config {
	return s.c
}

func (s *localSource) Size() (int64, error) {
	fi, err := s.f.Stat()
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func (s *localSource) Read(p []byte) (int, error) {
	return s.f.Read(p)
}

func (s *localSource) Close() error {
	return s.f.Close()
}
