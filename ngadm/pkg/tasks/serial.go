package tasks

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

type Serial struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	tasks      []Task
	name       string
}
type SerialParams struct {
}

func NewSerial(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	tasks := []Task{}
	for _, taskSpec := range taskSpec.SubTasks {
		task, err := GetTask(taskSpec, taskContext)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return &Serial{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		tasks:      tasks,
		name:       taskSpec.Description,
	}, nil
}

func (d *Serial) Execute() error {
	for _, task := range d.tasks {
		select {
		case <-d.JobContext.Sigs:
			return fmt.Errorf("signal received")
		default:
		}
		taskName := task.String()
		if taskName != "" {
			d.JobContext.Logger.Info(fmt.Sprintf("%s start", task))
		}
		d.JobContext.PushTask(task)
		err := task.Execute()
		if err != nil {
			return err
		}
		if taskName != "" {
			d.JobContext.Logger.Info(fmt.Sprintf("%s done", task))
		}
	}
	return nil
}

func (d *Serial) Rollback() error {
	for i := len(d.tasks) - 1; i >= 0; i-- {
		task := d.tasks[i]
		err := task.Rollback()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Serial) String() string {
	return d.name
}
