package tasks

// this task is for checking if the dir exists & the dir is empty, if not, return error for avoiding overwrite the dir
import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

type CheckDir struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	path       string
	host       string
	force      bool
}

type CheckDirParams struct {
	Path  string `json:"path"`
	Host  string `json:"host"`
	Force bool   `json:"force"`
}

func NewCheckDir(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*CheckDirParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	return &CheckDir{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		path:       params.Path,
		host:       params.Host,
		force:      params.Force,
	}, nil
}

func (d *CheckDir) Execute() error {
	if d.path == "" {
		return nil
	}

	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	if d.force {
		return d.RmDirs()
	}
	// check if the dir exists
	cmd := fmt.Sprintf("mkdir -p %s && test -d %s && ls %s", d.path, d.path, d.path)
	stdout, stderr, err := executor.Shell(cmd, false)
	if err != nil {
		return fmt.Errorf("check dir failed %s, %s, %v", stdout, stderr, err)
	}
	if len(stdout) > 0 {
		return fmt.Errorf("dir %s:%s is not empty", d.host, d.path)
	}
	return nil
}

func (d *CheckDir) RmDirs() error {
	if d.force {
		executor := d.JobContext.GetExecuter(d.host)
		if executor == nil {
			return fmt.Errorf("executor not found for host: %s", d.host)
		}
		cmd := fmt.Sprintf("rm -rf %s", d.path)
		stdout, stderr, err := executor.Shell(cmd, false)
		if err != nil {
			return fmt.Errorf("rm dir failed %s, %s, %w", stdout, stderr, err)
		}
	}
	return nil
}

func (d *CheckDir) Rollback() error {
	return nil
}

func (d *CheckDir) String() string {
	return fmt.Sprintf("CheckDir %s", d.path)
}
