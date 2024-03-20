package tasks

import (
	"os"
	"sync"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/executor"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/logger"
	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

type TaskGernerator func(taskSpec *types.TaskSpec, TaskContext *JobContext) (Task, error)

type JobContext struct {
	Logger      logger.Logger
	Progress    *types.Progress
	TasksTree   []Task
	mu          *sync.Mutex
	ExecuterMap map[string]executor.Executor
	ValueMap    map[string]any
	MaxParallel int
	Sigs        chan os.Signal
}

func NewJobContext() *JobContext {
	progress := &types.Progress{}
	return &JobContext{
		Progress:    progress,
		TasksTree:   []Task{},
		ValueMap:    map[string]any{},
		Logger:      logger.NewStdLogger(progress),
		ExecuterMap: map[string]executor.Executor{},
		mu:          &sync.Mutex{},
		MaxParallel: 100,
		Sigs:        make(chan os.Signal, 1),
	}
}

func (j *JobContext) SetExecuter(key string, executer executor.Executor) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.ExecuterMap[key] = executer
}

func (j *JobContext) GetExecuter(key string) executor.Executor {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.ExecuterMap[key]
}

func (j *JobContext) SetValue(key string, value any) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.ValueMap[key] = value
}

func (j *JobContext) GetValue(key string) any {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.ValueMap[key]
}

func (j *JobContext) GetLogger() logger.Logger {
	return j.Logger
}

func (j *JobContext) SetLogger(logger logger.Logger) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.Logger = logger
}

func (j *JobContext) PushTask(task Task) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.TasksTree = append(j.TasksTree, task)
	j.Progress.Current++
	if j.Progress.Current > j.Progress.Total {
		j.Progress.Current = j.Progress.Total
	}
}

func (j *JobContext) GetTasksTree() []Task {
	return j.TasksTree
}
