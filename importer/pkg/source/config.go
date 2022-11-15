package source

type Config struct {
	Path string `yaml:"path"`
}

func (c *Config) String() string {
	return c.Path
}
