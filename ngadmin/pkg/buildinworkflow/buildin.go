package buildinworkflow

import (
	"fmt"

	"github.com/vesoft-inc/nebula-ng-tools/ngadmin/pkg/types"
)

type WorkflowConverter func(args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error)

// buildin workflow converter map, don't need lock,just read
var BuildInWorkflowConverterMap = map[string]WorkflowConverter{
	"install":   Install,
	"verify":    Verify,
	"operation": Operation,
	"uninstall": Uninstall,
	"status":    Status,
	"config":    Config,
}

func GetBuildinWorkflow(cmd string, args map[string]any, spec *types.JobSpec) (*types.WorkflowSpec, error) {
	if converter, ok := BuildInWorkflowConverterMap[cmd]; ok {
		return converter(args, spec)
	}
	return nil, fmt.Errorf("workflow %s not found", cmd)
}
