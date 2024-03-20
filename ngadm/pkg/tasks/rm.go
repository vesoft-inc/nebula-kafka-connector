package tasks

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

type RMParams struct {
	Host string
	Path string
}

func NewRM(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*RMParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	return &RM{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		host:       params.Host,
		path:       params.Path,
	}, nil
}

type RM struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	host       string
	path       string
}

func (d *RM) Execute() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	cmd := fmt.Sprintf("rm -rf %s", d.path)
	stdout, stderr, err := executor.Shell(cmd, false)
	if err != nil || len(stderr) > 0 {
		return fmt.Errorf("execute command failed: %s, %s, %s", cmd, stdout, stderr)
	}
	return nil
}

func (d *RM) Rollback() error {
	// delete can't rollback
	return nil
}

func (d *RM) String() string {
	return fmt.Sprintf("delete path %s:%s", d.host, d.path)
}
