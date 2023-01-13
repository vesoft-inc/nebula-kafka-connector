package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Client  `yaml:"client"`
		Manager `yaml:"manager"`
		Sources `yaml:"sources"`
		*Log    `yaml:"log,omitempty"`
	}
)

func (c *Config) FromBytes(content []byte) error {
	return yaml.Unmarshal(content, &c)
}

func (c *Config) FromReader(r io.Reader) error {
	return yaml.NewDecoder(r).Decode(c)
}

func (c *Config) FromFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return c.FromReader(f)
}

func (c *Config) Yaml() (data []byte, err error) {
	return yaml.Marshal(c)
}

func (c *Config) Optimize(configPath string) error {
	if err := c.Log.OptimizeFilePath(configPath); err != nil {
		return err
	}

	if err := c.Sources.OptimizeConfigPath(configPath); err != nil {
		return err
	}
	sources, err := c.OptimizeWildCard()
	if err != nil {
		return err
	}
	c.Sources = sources

	return nil
}
