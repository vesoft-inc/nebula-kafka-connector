package tasks

import (
	"fmt"
	"sync"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
)

var TasksMap = map[string]TaskGernerator{}
var mu = &sync.Mutex{}

func Init() {
	TasksMap = map[string]TaskGernerator{
		"debug":            NewDebug,
		"shell":            NewShell,
		"connect":          NewConnect,
		"upload":           NewUpload,
		"serial":           NewSerial,
		"parallel":         NewParallel,
		"config":           NewConfig,
		"init_config":      NewInitConfig,
		"create_cluster":   NewCreateCluster,
		"nebula_operation": NewNebulaOperation,
		"extract":          NewExtract,
		"check_dir":        NewCheckDir,
	}

}

func RegisterTask(name string, f TaskGernerator) {
	mu.Lock()
	defer mu.Unlock()
	TasksMap[name] = f
}

func UnregisterTask(name string) {
	mu.Lock()
	defer mu.Unlock()
	delete(TasksMap, name)
}

func GetTask(task *types.TaskSpec, jobContext *JobContext) (Task, error) {
	mu.Lock()
	f, ok := TasksMap[task.Type]
	mu.Unlock() //unlock before call f, avoid deadlock
	if !ok {
		return nil, fmt.Errorf("task %s not found", task.Type)
	}
	taskInstance, err := f(task, jobContext)
	if err != nil {
		return nil, err
	}
	return taskInstance, nil
}
