package tasks

import (
	"fmt"
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
)

type ExtractParams struct {
	PkgPath     string
	ExtractPath string
	Host        string
	Sudo        bool
}

func NewExtract(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*ExtractParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	return &Extract{
		JobContext:  taskContext,
		taskSpec:    taskSpec,
		pkgPath:     params.PkgPath,
		extractPath: params.ExtractPath,
		host:        params.Host,
		sudo:        params.Sudo,
	}, nil
}

type Extract struct {
	JobContext  *JobContext
	taskSpec    *types.TaskSpec
	pkgPath     string
	extractPath string
	host        string
	sudo        bool
}

func (d *Extract) Execute() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	//mkdir first
	cmd := fmt.Sprintf("mkdir -p %s", d.extractPath)
	_, stderr, err := executor.Shell(cmd, d.sudo)
	if err != nil || stderr != "" {
		return fmt.Errorf("mkdir failed %s stderr: %s, err: %s", d.extractPath, string(stderr), err.Error())
	}
	pkgType := path.Ext(d.pkgPath)
	if pkgType == ".sh" {
		//chmod +x
		cmd = fmt.Sprintf("chmod +x %s", d.pkgPath)
		_, stderr, err = executor.Shell(cmd, d.sudo)
		if err != nil || stderr != "" {
			return fmt.Errorf("chmod failed %s stderr: %s, err: %s", d.pkgPath, string(stderr), err.Error())
		}
		cmd = fmt.Sprintf("%s --prefix=%s", d.pkgPath, d.extractPath)
	}
	_, stderr, err = executor.Shell(cmd, d.sudo)
	if err != nil {
		return fmt.Errorf("extract failed %s stderr: %s, err: %s", d.pkgPath, string(stderr), err.Error())
	}

	return nil
}

func (d *Extract) Rollback() error {
	cmd := fmt.Sprintf("rm -rf %s", d.extractPath)
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	_, stderr, err := executor.Shell(cmd, d.sudo)
	if err != nil {
		return fmt.Errorf("stderr: %s, err: %s", string(stderr), err.Error())
	}
	return nil
}

func (d *Extract) String() string {
	return "Extract"
}
