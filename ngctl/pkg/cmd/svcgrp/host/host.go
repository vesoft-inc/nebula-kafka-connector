package host

import (
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

type hostFlagsType struct {
	host       string
	agentPort  uint32
	configFile string
	output     string
	drain      bool
}

var hostFlags hostFlagsType

func validateOperationFlags() error {
	var flags = hostFlags

	if flags.host == "" {
		return common.NgctlError("must provide host info", "")
	}
	if flags.configFile == "" {
		return common.NgctlError("config file is empty", "")
	}
	return nil
}
