package tasks

import (
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

func NewDebug(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*DebugParams)
	if !ok {
		params = &DebugParams{
			Message: "debug task",
		}
	}
	return &Debug{
		JobContext: taskContext,
		params:     params,
	}, nil // Change Debug{} to &Debug{}
}

type DebugParams struct {
	Message string
}
type Debug struct {
	JobContext *JobContext
	params     *DebugParams
}

func (d *Debug) Execute() error { // Change (d *Debug) to (d *Debug)
	d.JobContext.Logger.Info("debug task executed!")
	//for test parallel log time
	time.Sleep(1 * time.Second)
	d.JobContext.Logger.Info(d.params.Message + " executed! time:" + time.Now().String())
	return nil
}

func (d *Debug) Rollback() error {
	d.JobContext.Logger.Info("debug task rollback!")
	time.Sleep(1 * time.Second)
	d.JobContext.Logger.Info(d.params.Message + " rollback! time:" + time.Now().String())
	return nil
}

func (d *Debug) String() string {
	return "debug" // debug task don't print anything
}
