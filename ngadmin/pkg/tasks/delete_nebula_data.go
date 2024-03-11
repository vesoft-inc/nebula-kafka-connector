package tasks

import (
	"fmt"
	"path"
	"strings"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/nebula"
	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
)

type DeleteNebulaDataParams struct {
	Host  string
	Path  string // nebula install path
	Drain bool
}

func NewDeleteNebulaDataTask(taskSpec *types.TaskSpec, taskContext *JobContext) (Task, error) {
	params, ok := taskSpec.Params.(*DeleteNebulaDataParams)
	if !ok {
		return nil, fmt.Errorf("unexpected type for params: %T", taskSpec.Params)
	}
	return &DeleteNebulaData{
		JobContext: taskContext,
		taskSpec:   taskSpec,
		host:       params.Host,
		path:       params.Path,
		drain:      params.Drain,
	}, nil
}

type DeleteNebulaData struct {
	JobContext *JobContext
	taskSpec   *types.TaskSpec
	host       string
	path       string
	drain      bool
}

var NebulaFiles = []string{
	"3rd", "bin", "etc", "include", "lib", "logs", "pids", "plugins", "scripts", "sids",
}

func (d *DeleteNebulaData) Execute() error {
	executor := d.JobContext.GetExecuter(d.host)
	if executor == nil {
		return fmt.Errorf("executor not found for host: %s", d.host)
	}

	// get config data from remote
	cmd := fmt.Sprintf("cat %s.default", path.Join(d.path, "etc/nebula-storaged.conf"))
	stdout, stderr, err := executor.Shell(cmd, false)
	if err != nil || len(stderr) > 0 {
		return fmt.Errorf("stderr: %s, err: %s", string(stderr), err)
	}
	storageConfig := string(stdout)

	cmd = fmt.Sprintf("cat %s.default", path.Join(d.path, "etc/nebula-metad.conf"))
	stdout, stderr, err = executor.Shell(cmd, false)
	if err != nil || len(stderr) > 0 {
		return fmt.Errorf("stderr: %s, err: %s", string(stderr), err)
	}
	metaConfig := string(stdout)
	// get process config
	cfg, err := nebula.NewTemplateWithData([]byte(storageConfig))
	if err != nil {
		return err
	}
	storage_data_path, err := cfg.GetValue("", "--data_path")
	if err != nil || len(storage_data_path) == 0 {
		storage_data_path = "data/storage"
	}
	storage_data_path = path.Join(d.path, storage_data_path)

	cfg, err = nebula.NewTemplateWithData([]byte(metaConfig))
	if err != nil {
		return err
	}
	meta_data_path, err := cfg.GetValue("", "--data_path")
	if err != nil || len(meta_data_path) == 0 {
		meta_data_path = "data/meta"
	}
	meta_data_path = path.Join(d.path, meta_data_path)

	// delete install path except data_path dir
	if d.drain {
		cmd = fmt.Sprintf("rm -rf %s && rm -rf %s && rm -rf %s", storage_data_path, meta_data_path, d.path)
	} else {
		if strings.Contains(storage_data_path, d.path) || strings.Contains(meta_data_path, d.path) {
			files := []string{}
			for _, file := range NebulaFiles {
				filePath := path.Join(d.path, file)
				files = append(files, filePath)
			}
			cmd = fmt.Sprintf("rm -rf %s", strings.Join(files, " ")) // delete nebulagraph files
		} else { // delete hole nebula install path if not data path
			cmd = fmt.Sprintf("rm -rf %s", d.path)
		}
	}
	_, stderr, err = executor.Shell(cmd, false)
	if err != nil || len(stderr) > 0 {
		return fmt.Errorf("stderr: %s, err: %s", string(stderr), err)
	}
	return nil
}

func (d *DeleteNebulaData) Rollback() error {
	return nil
}

func (d *DeleteNebulaData) String() string {
	return fmt.Sprintf("delete data %s drain:%v", d.host, d.drain)
}
