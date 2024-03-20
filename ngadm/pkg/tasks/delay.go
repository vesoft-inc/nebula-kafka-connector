package tasks

import (
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

type DelayParams struct {
	Duration time.Duration
}

type Delay struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	params     *DelayParams
}

func NewDelay(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*DelayParams)
	if !ok {
		params = &DelayParams{
			Duration: 1 * time.Second,
		}
	}
	return &Delay{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		params:     params,
	}, nil
}

func (d *Delay) Execute() error {
	time.Sleep(d.params.Duration)
	return nil
}

func (d *Delay) Rollback() error {
	time.Sleep(d.params.Duration)
	return nil
}

func (d *Delay) String() string {
	return d.taskSpec.Description
}
