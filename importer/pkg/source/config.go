package source

type (
	Config struct {
		Path string     `yaml:"path"`
		CSV  *CSVConfig `yaml:"csv,omitempty"`
	}

	CSVConfig struct {
		Delimiter string `yaml:"delimiter,omitempty"`
	}
)

func (c *Config) String() string {
	return c.Path
}
