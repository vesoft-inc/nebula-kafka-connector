package tasks

import (
	"fmt"
	"regexp"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/nebula"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/yamlparser"
)

type ConfigParams struct {
	Host   string
	Config map[string]any
	Dst    string
	Sudo   bool
	Type   string //yaml,nebulagraph
}

func NewConfig(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*ConfigParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	return &Config{
		JobContext:   taskContext,
		taskSpec:     taskSpec,
		host:         params.Host,
		config:       params.Config,
		dst:          params.Dst,
		originConfig: "",
		sudo:         params.Sudo,
		typo:         params.Type,
	}, nil
}

type Config struct {
	JobContext   *JobContext
	taskSpec     *types.TaskSpec
	host         string
	config       map[string]any
	dst          string
	originConfig string
	sudo         bool
	typo         string
}

func (d *Config) Execute() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}

	// get config data from remote
	cmd := fmt.Sprintf("cat %s", d.dst)
	stdout, stderr, err := executor.Shell(cmd, d.sudo)
	if err != nil {
		return fmt.Errorf("stderr: %s, err: %s", string(stderr), err)
	}
	d.originConfig = string(stdout)
	configString := ""
	if d.typo == "nebulagraph" {
		// update process config
		cfg, err := nebula.NewTemplateWithData([]byte(d.originConfig))
		if err != nil {
			return err
		}
		//map[string]any to map[string]string
		config := make(map[string]string)
		for k, v := range d.config {
			vs, ok := v.(string)
			if ok {
				config["--"+k] = vs
			}
		}
		cfg.UpdateConfig(config)
		configString, err = cfg.String()
		if err != nil {
			return err
		}
	} else if d.typo == "yaml" {
		// update yaml config
		configString, err = yamlparser.ApplyConfigByYamlString(d.originConfig, d.config)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported config type: %s", d.typo)
	}

	// use regex to replace config ' "
	re := regexp.MustCompile("`")
	configString = re.ReplaceAllString(configString, "'")

	// write config to remote
	cmd = fmt.Sprintf(`cat>%s<<EOF
%s
EOF`, d.dst, configString)
	_, stderr, err = executor.Shell(cmd, d.sudo)
	if err != nil {
		return fmt.Errorf("stderr: %s, err: %s", string(stderr), err)
	}

	return nil
}

func (d *Config) Rollback() error {
	if d.originConfig == "" {
		return nil
	}
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	cmd := fmt.Sprintf(`cat>%s<<EOF
%s
EOF`, d.dst, d.originConfig)
	_, stderr, err := executor.Shell(cmd, d.sudo)
	if err != nil {
		return fmt.Errorf("stderr: %s, err: %s", string(stderr), err)
	}
	return nil
}

func (d *Config) String() string {
	return "Config " + d.host + " " + d.dst
}
