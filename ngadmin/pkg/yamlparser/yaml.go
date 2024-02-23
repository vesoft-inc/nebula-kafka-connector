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
		metaSpec.Config = make(map[string]any)
	}
	for _, cluster := range metaSpec.Clusters {
		if cluster.Graphd.Hosts == nil {
			return fmt.Errorf("graphd hosts is nil")
		}
		if cluster.Graphd.Config == nil {
			cluster.Graphd.Config = make(map[string]any)
		}
		if cluster.Storaged.Hosts == nil {
			return fmt.Errorf("storaged hosts is nil")
		}
		if cluster.Storaged.Config == nil {
			cluster.Storaged.Config = make(map[string]any)
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

type anyMap map[string]interface{}

func merge(dst, src anyMap) anyMap {
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			srcMap, srcOk := srcVal.(anyMap)
			dstMap, dstOk := dstVal.(anyMap)
			if srcOk && dstOk {
				dst[key] = merge(dstMap, srcMap)
			} else {
				dst[key] = srcVal
			}
		} else {
			dst[key] = srcVal
		}
	}
	return dst
}

func ApplyConfigByYamlString(yamlString string, config anyMap) (string, error) {
	var values anyMap = make(map[string]any)
	err := yaml.Unmarshal([]byte(yamlString), &values)
	if err != nil {
		return "", err
	}

	values = merge(values, config)

	newYamlString, err := yaml.Marshal(values)
	if err != nil {
		return "", err
	}

	return string(newYamlString), nil
}
