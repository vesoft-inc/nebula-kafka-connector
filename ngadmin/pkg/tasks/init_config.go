package tasks

import (
	"fmt"
	"regexp"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/nebula"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
)

type InitConfigParams struct {
	Host         string
	Dst          string
	ChangeMap    map[string]string
	OriginConfig string
	Sudo         bool
}

var NebulaInitConfigKey = "nebula_init_config_"

func NewInitConfig(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*InitConfigParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	changeMap := map[string]string{}
	for k, v := range params.ChangeMap {
		changeMap["--"+k] = v
	}
	return &InitConfig{
		JobContext:   taskContext,
		taskSpec:     taskSpec,
		host:         params.Host,
		dst:          params.Dst,
		changeMap:    changeMap,
		originConfig: params.OriginConfig,
		sudo:         params.Sudo,
	}, nil
}

type InitConfig struct {
	JobContext   *JobContext
	taskSpec     *types.TaskSpec
	host         string
	dst          string
	changeMap    map[string]string
	originConfig string
	sudo         bool
}

func (d *InitConfig) Execute() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}

	if d.originConfig == "" {
		// get config data from remote
		cmd := fmt.Sprintf("cat %s.default", d.dst)
		stdout, stderr, err := executor.Shell(cmd, d.sudo)
		if err != nil {
			return fmt.Errorf("stderr: %s, err: %s", string(stderr), err.Error())
		}
		d.originConfig = string(stdout)
	}

	// update process config
	cfg, err := nebula.NewTemplateWithData([]byte(d.originConfig))
	if err != nil {
		return err
	}
	cfg.UpdateConfig(d.changeMap)
	config, err := cfg.String()
	if err != nil {
		return err
	}

	// use regex to replace config ' "
	re := regexp.MustCompile("`")
	config = re.ReplaceAllString(config, "'")

	// write config to remote
	cmd := fmt.Sprintf(`cat>%s<<EOF
%s
EOF`, d.dst, config)
	_, stderr, err := executor.Shell(cmd, d.sudo)
	if err != nil {
		return fmt.Errorf("stderr: %s, err: %s", string(stderr), err.Error())
	}

	// for save config to database
	key := NebulaInitConfigKey + d.host + d.dst
	d.JobContext.SetValue(key, config)
	return nil
}

func (d *InitConfig) Rollback() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	// executor.Shell(fmt.Sprintf("rm -rf %s", d.dst), d.sudo) //dont need rm config file,rollback for install will rm all
	return nil
}

func (d *InitConfig) String() string {
	if d.taskSpec.Description != "" {
		return d.taskSpec.Description
	}
	return "InitConfig"
}
