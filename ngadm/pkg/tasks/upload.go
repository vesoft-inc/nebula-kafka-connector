package tasks

import (
	"fmt"
	"path"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

type UploadParams struct {
	SrcPath string `json:"srcPath"`
	DstPath string `json:"dstPath"`
	Host    string `json:"host"`
}

func NewUpload(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*UploadParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}

	return &Upload{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		srcPath:    params.SrcPath,
		dstPath:    params.DstPath,
		host:       params.Host,
	}, nil
}

type Upload struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	srcPath    string
	dstPath    string
	host       string
}

func (d *Upload) Execute() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	// make dir
	stdout, stderr, err := executor.Shell(fmt.Sprintf("mkdir -p %s", path.Dir(d.dstPath)), false)
	if err != nil || stderr != "" {
		return fmt.Errorf("mkdir failed:%s, %s,%w", stdout, stderr, err)
	}

	stdout, stderr, err = executor.Upload(d.srcPath, d.dstPath)
	if err != nil || stderr != "" {
		return fmt.Errorf("upload failed:%s, %s,%w", stdout, stderr, err)
	}
	return nil
}

func (d *Upload) Rollback() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}
	stdout, stderr, err := executor.Shell(fmt.Sprintf("rm -rf %s", path.Dir(d.dstPath)), false)
	if err != nil || stderr != "" {
		return fmt.Errorf("rollback failed:%s, %s,%w", stdout, stderr, err)
	}
	return nil
}

func (d *Upload) String() string {
	return "Upload " + d.srcPath + " to " + d.host + ":" + d.dstPath
}
