package yamlparser

import "github.com/vesoft-inc/nebula-ng-tools/ngadm/pkg/types"

var UtilsProcesses = map[string]types.Process{
	"license-manager": {
		Name:          "license-manager",                                          // the name of the process for install directory & product name
		ExecShellPath: "./nebula-license-manager/scripts/license-manager.service", // for shell start
		ExecStartPath: "./nebula-license-manager/nebula-license-manager",          // for systemd start
		WorkingDir:    "./nebula-license-manager/",                                // for systemd start
		ConfigPath:    "./nebula-license-manager/etc/nebula-license-manager.yaml", // for merge user config to your product config
	},
	"agent": {
		Name:          "nebula-agent",            // the name of the process for install directory & product name
		ExecShellPath: "./scripts/agent.service", // for shell start
		ExecStartPath: "./agent",                 // for systemd start
		WorkingDir:    "./",                      // for systemd start
		ConfigPath:    "./etc/agent-api.yaml",    // for merge user config to your product config
	},
}
