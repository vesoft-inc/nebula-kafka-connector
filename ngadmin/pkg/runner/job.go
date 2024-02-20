package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/buildinworkflow"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/logger"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/tasks"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
)

type Job struct {
	Name         string
	WorkflowSpec *types.WorkflowSpec //save input info for lookup the job later
	JobSpec      *types.JobSpec
	CMD          string
	Args         map[string]any
	StartTime    time.Time
	EndTime      time.Time
	Context      *tasks.JobContext
}

// only for one instance
func NewJob(name string) *Job {
	return &Job{
		Name:    name,
		Context: tasks.NewJobContext(),
	}
}

func (j *Job) Run(cmd string, args map[string]any, spec *types.JobSpec) error {
	workflow, err := buildinworkflow.GetBuildinWorkflow(cmd, args, spec)
	if workflow == nil || err != nil {
		return fmt.Errorf("workflow %s build err:%s", cmd, err.Error())
	}
	j.CMD = cmd
	j.Args = args
	j.JobSpec = spec
	return j.RunWorkflow(workflow)
}

func (j *Job) RunWorkflow(workflow *types.WorkflowSpec) error {
	j.WorkflowSpec = workflow
	j.Context.TasksTree = []tasks.Task{}
	j.Context.Progress.Total = j.GetTotalTaskNum(workflow)
	j.Context.Logger.Info(fmt.Sprintf("job \"%s\" start", j.Name))

	tasks := workflow.Tasks
	if workflow.Type == "" {
		workflow.Type = "serial"
	}
	taskSpec := &types.TaskSpec{
		Type:     workflow.Type,
		Params:   workflow.Params,
		SubTasks: tasks,
	}
	j.ListenSafeExit()
	defer j.UnListenSafeExit()
	_, err := j.RunTask(taskSpec)
	if err != nil {
		if workflow.Rollback {
			rollbackErr := j.Rollback()
			if rollbackErr != nil {
				j.Context.Logger.Error(err.Error())
				return fmt.Errorf("run task faild: %v \n rollback failed: %v", err, rollbackErr)
			}
		}
		return err
	}
	j.Context.Logger.Info(fmt.Sprintf("job \"%s\" done", j.Name))
	return nil
}

func (j *Job) ListenSafeExit() {
	signal.Notify(j.Context.Sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-j.Context.Sigs // all parallel & serial task will stop on next subtask
		j.Context.Logger.Info(fmt.Sprintf("job %s exit", j.Name))
		if j.WorkflowSpec.Rollback {
			rollbackErr := j.Rollback()
			if rollbackErr != nil {
				j.Context.Logger.Error(fmt.Sprintf("rollback failed: %v", rollbackErr))
			}
		}
		os.Exit(0)
	}()
}

func (j *Job) UnListenSafeExit() {
	signal.Stop(j.Context.Sigs)
}

func (j *Job) Rollback() error {
	for i := len(j.Context.TasksTree) - 1; i >= 0; i-- {
		task := j.Context.TasksTree[i]
		err := task.Rollback()
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *Job) RunTask(task *types.TaskSpec) (tasks.Task, error) {
	taskInstance, err := tasks.GetTask(task, j.Context)
	if err != nil {
		j.Context.Logger.Error(err.Error())
		return nil, err
	}
	taskName := taskInstance.String()
	if taskName != "" {
		j.Context.Logger.Info(fmt.Sprintf("%s start", taskInstance))
	}
	err = taskInstance.Execute()
	j.Context.PushTask(taskInstance)
	if err != nil {
		j.Context.Logger.Error(err.Error())
		return nil, err
	}
	if taskName != "" {
		j.Context.Logger.Info(fmt.Sprintf("%s done", taskInstance))
	}
	return taskInstance, nil
}

func (j *Job) SetLogger(logger logger.Logger) {
	j.Context.Logger = logger
}

func (j *Job) JSON() (string, error) {
	bytes, err := json.Marshal(map[string]any{
		"name":         j.Name,
		"progress":     j.Context.Progress,
		"valueMap":     j.Context.ValueMap,
		"workflowSpec": j.WorkflowSpec,
		"jobSpec":      j.JobSpec,
	})
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (j *Job) String() string {
	return fmt.Sprintf("job %s", j.Name)
}

func (j *Job) GetTotalTaskNum(workflow *types.WorkflowSpec) int {
	total := 0
	for _, task := range workflow.Tasks {
		total += dfs(task, 0)
	}
	return total
}

func dfs(task *types.TaskSpec, num int) int {
	num++
	for _, t := range task.SubTasks {
		num = dfs(t, num)
	}
	return num
}
