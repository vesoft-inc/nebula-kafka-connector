package tasks

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
	"golang.org/x/sync/errgroup"
)

type Parallel struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	tasks      []Task
	name       string
}

type ParallelParams struct {
}

func NewParallel(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	tasks := []Task{}
	for _, taskSpec := range taskSpec.SubTasks {
		task, err := GetTask(taskSpec, taskContext)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return &Parallel{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		tasks:      tasks,
		name:       taskSpec.Description,
	}, nil
}

func (d *Parallel) Execute() error {
	g := &errgroup.Group{}
	workerPool := make(chan struct{}, d.JobContext.MaxParallel)
	for _, t := range d.tasks {
		workerPool <- struct{}{}
		t := t // https://golang.org/doc/faq#closures_and_goroutines
		taskName := t.String()
		g.Go(func() error {
			defer func() {
				if taskName != "" {
					d.JobContext.Logger.Info(fmt.Sprintf("parallel %s done", t))
				}
				<-workerPool
			}()
			select {
			case <-d.JobContext.Sigs:
				return fmt.Errorf("signal received")
			default:
			}
			if taskName != "" {
				d.JobContext.Logger.Info(fmt.Sprintf("parallel %s start", t))
			}
			d.JobContext.PushTask(t)
			return t.Execute()
		})
	}
	err := g.Wait()
	return err
}

func (d *Parallel) Rollback() error {
	g := &errgroup.Group{}
	workerPool := make(chan struct{}, d.JobContext.MaxParallel)
	for i := len(d.tasks) - 1; i >= 0; i-- {
		workerPool <- struct{}{}
		task := d.tasks[i]
		g.Go(func() error {
			defer func() {
				<-workerPool
			}()
			return task.Rollback()
		})
	}
	err := g.Wait()
	return err
}

func (d *Parallel) String() string {
	return d.name
}
