package yamlparser

import (
	"fmt"
	"os"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
	"gopkg.in/yaml.v2"
)

// todo: parser
func ParseYamlByPath(path string) (jobSpec *types.JobSpec, err error) {
	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Unmarshal the YAML data into a JobSpec object
	jobSpec = &types.JobSpec{}
	err = yaml.Unmarshal(data, jobSpec)
	if err != nil {
		return nil, err
	}
	// append needed config fields
	CheckMetaSpec(jobSpec.Spec.Metad)
	return jobSpec, err
}

func CheckMetaSpec(metaSpec *types.MetadSpec) error {
	if metaSpec == nil {
		return fmt.Errorf("metad spec is nil")
	}
	if metaSpec.Hosts == nil {
		return fmt.Errorf("metad hosts is nil")
	}
	if metaSpec.Clusters == nil {
		return fmt.Errorf("metad clusters is nil")
	}
	if metaSpec.Config == nil {
		metaSpec.Config = make(map[string]string)
	}
	if _, ok := metaSpec.Config["port"]; !ok {
		metaSpec.Config["port"] = "9559"
	}
	for _, cluster := range metaSpec.Clusters {
		if cluster.Graphd.Hosts == nil {
			return fmt.Errorf("graphd hosts is nil")
		}
		if cluster.Graphd.Config == nil {
			cluster.Graphd.Config = make(map[string]string)
		}
		if cluster.Storaged.Hosts == nil {
			return fmt.Errorf("storaged hosts is nil")
		}
		if cluster.Storaged.Config == nil {
			cluster.Storaged.Config = make(map[string]string)
		}
	}
	return nil
}

func CacheConfig(jobSpec *types.JobSpec) error {
	yamlString, err := yaml.Marshal(jobSpec)
	if err != nil {
		return err
	}
	err = os.WriteFile("ngadmin.yaml", yamlString, 0644)
	if err != nil {
		return err
	}
	return nil
}
