package yamlparser

import (
	"fmt"
	"os"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"gopkg.in/yaml.v3"
)

type SpecMap struct {
	Spec map[string]types.Process `yaml:"spec,omitempty"`
}

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
	anyMap := &SpecMap{}
	err = yaml.Unmarshal(data, &anyMap)
	if err != nil {
		return nil, err
	}
	ParseUtilsProcess(jobSpec, anyMap)
	// append needed config fields
	CheckMetaSpec(jobSpec.Spec.Metad)
	return jobSpec, err
}

func ParseUtilsProcess(jobSpec *types.JobSpec, specMap *SpecMap) {
	if jobSpec.UtilsProcesses == nil {
		jobSpec.UtilsProcesses = make(map[string]*types.Process)
	}
	spec := specMap.Spec
	if spec == nil {
		return
	}
	for key, process := range spec {
		if key == "metad" {
			continue
		}
		// apply default process
		defaultProcess, ok := UtilsProcesses[key]
		if ok {
			if process.ExecShellPath == "" {
				process.ExecShellPath = defaultProcess.ExecShellPath
			}
			if process.ExecStartPath == "" {
				process.ExecStartPath = defaultProcess.ExecStartPath
			}
			if process.WorkingDir == "" {
				process.WorkingDir = defaultProcess.WorkingDir
			}
			if process.ConfigPath == "" {
				process.ConfigPath = defaultProcess.ConfigPath
			}
			if process.Name == "" {
				process.Name = defaultProcess.Name
			}
		}
		if process.Name == "" {
			process.Name = key
		}
		jobSpec.UtilsProcesses[key] = &process
	}
}

func CheckMetaSpec(metaSpec *types.MetadSpec) error {
	if metaSpec == nil {
		return fmt.Errorf("metad spec is nil")
	}
	if metaSpec.Hosts == nil {
		return fmt.Errorf("metad hosts is nil")
	}
	if metaSpec.ServiceGroups == nil {
		return fmt.Errorf("metad clusters is nil")
	}
	if metaSpec.Config == nil {
		metaSpec.Config = make(map[string]any)
	}
	for _, cluster := range metaSpec.ServiceGroups {
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
	err = os.WriteFile("ngadm.yaml", yamlString, 0644)
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
