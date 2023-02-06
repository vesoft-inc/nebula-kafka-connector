package config

import (
	"io"
	"os"

	configbase "github.com/vesoft-inc/nebula-ng-tools/importer/pkg/config/base"
	configv3 "github.com/vesoft-inc/nebula-ng-tools/importer/pkg/config/v3"
	configv5 "github.com/vesoft-inc/nebula-ng-tools/importer/pkg/config/v5"
	"github.com/vesoft-inc/nebula-ng-tools/importer/pkg/errors"

	"gopkg.in/yaml.v3"
)

type (
	Client       = configbase.Client
	Log          = configbase.Log
	Configurator = configbase.Configurator
)

func FromBytes(content []byte) (Configurator, error) {
	type tmpConfig struct {
		Client struct {
			Version string `yaml:"version"`
		} `yaml:"client"`
	}
	var tc tmpConfig
	if err := yaml.Unmarshal(content, &tc); err != nil {
		return nil, err
	}
	var c Configurator
	switch tc.Client.Version {
	case configbase.ClientVersion3:
		c = &configv3.Config{}
	case "", configbase.ClientVersion5:
		c = &configv5.Config{}
	default:
		return nil, errors.ErrUnsupportedClientVersion
	}

	if err := yaml.Unmarshal(content, c); err != nil {
		return nil, err
	}
	return c, nil
}

func FromReader(r io.Reader) (Configurator, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return FromBytes(content)
}

func FromFile(name string) (Configurator, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return FromReader(f)
}
