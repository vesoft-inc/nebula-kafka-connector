package service_admin

import (
	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
)

type serviceFlagsType struct {
	serviceType string
	host        string
	port        int32
	svcgrpName  string
	configFile  string
	output      string
}

var serviceFlags serviceFlagsType

func validateAddorDropFlags() error {
	var flags = serviceFlags
	if flags.svcgrpName == "" {
		return common.NgctlError("svcgrp name is empty", "")
	}
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
	if flags.svcgrpName == "" {
		return common.NgctlError("svcgrp name is empty", "")
	}
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
	default:
		return meta.ServiceTypeGraphd, common.NgctlError("Invalid service type, "+typ, "")
	}
}

var ServiceAdminCmd = &cobra.Command{
	Use:   "service",
	Short: "Run commands managing services in the group",
	Long:  `Run commands managing services in the group`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ServiceAdminCmd.AddCommand(addServiceCmd)
	ServiceAdminCmd.AddCommand(showServiceCmd)
	ServiceAdminCmd.AddCommand(dropServiceCmd)
	ServiceAdminCmd.AddCommand(startServiceCmd)
	ServiceAdminCmd.AddCommand(stopServiceCmd)
	ServiceAdminCmd.AddCommand(configServiceCmd)
}
