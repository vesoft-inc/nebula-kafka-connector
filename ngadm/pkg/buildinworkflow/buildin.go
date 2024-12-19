package buildinworkflow

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"
)

type WorkflowConverter func(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error)

const (
	INSTALL_PATH_FILE = "install_path"
	INSTALL_PATH_DIR  = "~/.nebulagraph"
)

// buildin workflow converter map, don't need lock,just read
var BuildInWorkflowConverterMap = map[string]WorkflowConverter{
	"install":                 Install,
	"verify":                  Verify,
	"operation":               Operation,
	"uninstall":               Uninstall,
	"status":                  Status,
	"config":                  Config,
	"install-agent":           InstallAgent,
	"install-host":            InstallHost,
	"uninstall-host":          UninstallHost,
	"upload-file":             UploadFile,
	"install-service-group":   InstallServiceGroup,
	"uninstall-service-group": UninstallServiceGroup,
}

func GetBuildinWorkflow(cmd string, args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	if converter, ok := BuildInWorkflowConverterMap[cmd]; ok {
		return converter(args, spec)
	}
	return nil, fmt.Errorf("workflow %s not found", cmd)
}
