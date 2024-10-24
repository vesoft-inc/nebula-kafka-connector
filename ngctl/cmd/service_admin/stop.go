package service_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/job"
)

var stopServiceCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop service",
	Long:  "Stop service",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return common.MetaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateOperationFlags(); err != nil {
			return err
		}
		c, err := config.NewConfigFromFile(serviceFlags.configFile)
		if err != nil {
			return common.NgctlError("Failed to parse config file for the install path", err.Error())
		}
		instances, err := c.GetServiceInstances(serviceFlags.svcgrpName,
			serviceFlags.serviceType, serviceFlags.host)
		if err != nil {
			return err
		}
		hosts := make([]string, 0)
		for _, inst := range instances {
			hosts = append(hosts, inst.Host)
		}
		if err := job.NgadmJob.ServiceOperation(c, hosts, serviceFlags.serviceType, "stop"); err != nil {
			return err
		}

		fmt.Fprintf(common.MetaOutput, "Stop service successfully.\n")
		return nil
	},
}

func init() {
	stopServiceCmd.Flags().StringVarP(&serviceFlags.svcgrpName, "svcgrp", "s", "", "svcgrp name")
	stopServiceCmd.Flags().StringVarP(&serviceFlags.serviceType, "type", "t", "", "service type")
	stopServiceCmd.Flags().StringVarP(&serviceFlags.host, "host", "H", "", "host")
	stopServiceCmd.Flags().StringVarP(&serviceFlags.configFile, "config", "f", "", "config file for ngctl")
}
