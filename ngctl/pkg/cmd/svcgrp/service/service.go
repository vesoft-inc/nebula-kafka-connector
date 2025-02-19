package service

import (
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/cmd/common"
)

type serviceFlagsType struct {
	serviceType string
	host        string
	port        int32
	configFile  string
	output      string
}

var serviceFlags serviceFlagsType

func validateAddOrDropFlags() error {
	var flags = serviceFlags
	if flags.host == "" {
		return common.NgctlError("must provide host info", "")
	}
	if flags.port == -1 {
		return common.NgctlError("must provide port info", "")
	}
	if flags.serviceType == "" {
		return common.NgctlError("must provide service type", "")
	}
	return nil
}

func validateOperationFlags() error {
	var flags = serviceFlags
	if flags.configFile == "" {
		return common.NgctlError("config file is empty", "")
	}
	if flags.serviceType == "" {
		return common.NgctlError("must provide service type", "")
	}
	return nil
}

// get service type, not include metad
func getServiceType(typ string) (meta.ServiceType, error) {
	switch typ {
	case "graphd":
		return meta.ServiceTypeGraphd, nil
	case "storaged":
		return meta.ServiceTypeStoraged, nil
	case "analyticd":
		return meta.ServiceTypeAnalyticd, nil
	default:
		return meta.ServiceTypeGraphd, common.NgctlError("Invalid service type, "+typ, "")
	}
}
