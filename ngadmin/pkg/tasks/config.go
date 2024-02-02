package tasks

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
)

type ConfigParams struct {
	Host   string
	Config map[string]string
}

func NewConfig(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*ConfigParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	return &Config{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		host:       params.Host,
		config:     params.Config,
	}, nil
}

type Config struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	host       string
	config     map[string]string
}

func (d *Config) Execute() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}

	return nil
}

func (d *Config) Rollback() error {
	return nil
}

func (d *Config) String() string {
	return "Config"
}
