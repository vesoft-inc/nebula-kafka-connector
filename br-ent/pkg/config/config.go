package config

import (
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/yamlparser"
)

func ParseConfig(configFilePath string) (configSpec *types.JobSpec, err error) {
	spec, err := yamlparser.ParseYamlByPath(configFilePath)
	if err != nil {
		return nil, err
	}

	return spec, nil
}
