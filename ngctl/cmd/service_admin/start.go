package service_admin

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/cmd/common"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/config"
	"github.com/vesoft-inc/nebula-ng-tools/ngctl/pkg/job"
)

var startServiceCmd = &cobra.Command{
	Use:   "start",
	Short: "Start service",
	Long:  "Start service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateOperationFlags(); err != nil {
			return err
		}
		c, err := config.NewConfigFromFile(serviceFlags.configFile)
		if err != nil {
			return err
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
		if err := job.NgadmJob.ServiceOperation(c, hosts, serviceFlags.serviceType, "start"); err != nil {
			return common.NgctlError("Failed to start service", err.Error())
		}
		fmt.Fprintf(common.MetaOutput, "Start service successfully.\n")
		return nil
	},
}

func init() {
	startServiceCmd.Flags().StringVarP(&serviceFlags.svcgrpName, "svcgrp", "s", "", "svcgrp name")
	startServiceCmd.Flags().StringVarP(&serviceFlags.serviceType, "type", "t", "", "service type")
	startServiceCmd.Flags().StringVarP(&serviceFlags.host, "host", "H", "", "host")
	startServiceCmd.Flags().StringVarP(&serviceFlags.configFile, "config", "f", "", "config file for ngctl")
}
